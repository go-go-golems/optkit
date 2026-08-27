package sqlite

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/go-go-golems/optkit/budget"
	sqlitedb "github.com/go-go-golems/optkit/internal/sqlite"
	"github.com/go-go-golems/optkit/record"
)

func (s *Store) Define(ctx context.Context, campaignID record.CampaignID, limits []budget.Limit) error {
	if err := record.ValidateID("campaign", string(campaignID)); err != nil {
		return err
	}
	normalized, err := budget.NormalizeLimits(limits)
	if err != nil {
		return err
	}
	return s.db.InTx(ctx, func(conn *sqlitedb.Conn) error {
		head, err := conn.Query(ctx, `SELECT version FROM campaign_heads WHERE campaign_id = ?`, string(campaignID))
		if err != nil {
			return err
		}
		if len(head) != 1 {
			return fmt.Errorf("%w: campaign %s has no journal head", budget.ErrNotDefined, campaignID)
		}
		existing, err := loadLimitsConn(ctx, conn, campaignID)
		if err != nil {
			return err
		}
		if len(existing) > 0 {
			if !limitsEqual(existing, normalized) {
				return budget.ErrLimitConflict
			}
			return nil
		}
		for _, limit := range normalized {
			if _, err := conn.Exec(ctx, `INSERT INTO budget_limits(campaign_id, resource, limit_units) VALUES(?, ?, ?)`, string(campaignID), string(limit.Resource), limit.Units); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) Reserve(ctx context.Context, request budget.ReserveRequest) (budget.Reservation, error) {
	claims, err := budget.NormalizeQuantities(request.Claims, false)
	if err != nil {
		return budget.Reservation{}, err
	}
	request.Claims = claims
	reservationID, err := budget.ReservationID(request)
	if err != nil {
		return budget.Reservation{}, err
	}
	var result budget.Reservation
	err = s.db.InTx(ctx, func(conn *sqlitedb.Conn) error {
		existing, found, err := reservationByWorkConn(ctx, conn, request.Work)
		if err != nil {
			return err
		}
		if found {
			if existing.ID == reservationID && existing.Campaign == request.Campaign && quantitiesEqual(existing.Requested, claims) {
				result = existing
				return nil
			}
			return budget.ErrReservationConflict
		}
		limits, err := loadLimitsConn(ctx, conn, request.Campaign)
		if err != nil {
			return err
		}
		if len(limits) == 0 {
			return budget.ErrNotDefined
		}
		usage, _, err := usageConn(ctx, conn, request.Campaign, "")
		if err != nil {
			return err
		}
		limitMap := limitMap(limits)
		for _, claim := range claims {
			limit, ok := limitMap[claim.Resource]
			if !ok {
				return fmt.Errorf("%w: resource %s has no limit", budget.ErrNotDefined, claim.Resource)
			}
			if usage[claim.Resource]+claim.Units > limit {
				return fmt.Errorf("%w: resource %s needs %d, available %d", budget.ErrInsufficient, claim.Resource, claim.Units, limit-usage[claim.Resource])
			}
		}
		requestedJSON, err := record.CanonicalJSON(claims)
		if err != nil {
			return err
		}
		emptyActual, err := record.CanonicalJSON([]budget.Quantity{})
		if err != nil {
			return err
		}
		now := s.now().UTC()
		if _, err := conn.Exec(ctx, `INSERT INTO budget_reservations(id, campaign_id, work_id, requested_json, actual_json, status, overage, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?, 0, ?, ?)`,
			string(reservationID), string(request.Campaign), string(request.Work), requestedJSON, emptyActual, string(budget.StatusReserved), formatTime(now), formatTime(now)); err != nil {
			return err
		}
		result = budget.Reservation{
			ID: reservationID, Campaign: request.Campaign, Work: request.Work,
			Requested: claims, Status: budget.StatusReserved, CreatedAt: now, UpdatedAt: now,
		}
		return nil
	})
	return result, err
}

func (s *Store) Commit(ctx context.Context, reservationID record.ReservationID, actual []budget.Quantity) (budget.Reservation, error) {
	if err := record.ValidateID("reservation", string(reservationID)); err != nil {
		return budget.Reservation{}, err
	}
	normalized, err := budget.NormalizeQuantities(actual, true)
	if err != nil {
		return budget.Reservation{}, err
	}
	var result budget.Reservation
	err = s.db.InTx(ctx, func(conn *sqlitedb.Conn) error {
		reservation, found, err := reservationByIDConn(ctx, conn, reservationID)
		if err != nil {
			return err
		}
		if !found {
			return budget.ErrReservationNotFound
		}
		switch reservation.Status {
		case budget.StatusCommitted:
			if quantitiesEqual(reservation.Actual, normalized) {
				result = reservation
				return nil
			}
			return budget.ErrReservationConflict
		case budget.StatusReleased:
			return budget.ErrReservationTerminal
		case budget.StatusReserved:
		default:
			return fmt.Errorf("unknown reservation status %q", reservation.Status)
		}
		requested := quantityMap(reservation.Requested)
		for _, value := range normalized {
			if _, ok := requested[value.Resource]; !ok {
				return fmt.Errorf("actual resource %s was not reserved", value.Resource)
			}
		}
		limits, err := loadLimitsConn(ctx, conn, reservation.Campaign)
		if err != nil {
			return err
		}
		usage, _, err := usageConn(ctx, conn, reservation.Campaign, reservation.ID)
		if err != nil {
			return err
		}
		before := make(map[budget.Resource]int64, len(usage))
		for resource, units := range usage {
			before[resource] = units
		}
		overage := false
		for _, value := range normalized {
			if value.Units > requested[value.Resource] {
				overage = true
			}
			usage[value.Resource] += value.Units
		}
		for _, limit := range limits {
			if before[limit.Resource] <= limit.Units && usage[limit.Resource] > limit.Units {
				overage = true
			}
		}
		actualJSON, err := record.CanonicalJSON(normalized)
		if err != nil {
			return err
		}
		now := s.now().UTC()
		changes, err := conn.Exec(ctx, `UPDATE budget_reservations SET actual_json = ?, status = ?, overage = ?, updated_at = ? WHERE id = ? AND status = ?`,
			actualJSON, string(budget.StatusCommitted), boolInt(overage), formatTime(now), string(reservationID), string(budget.StatusReserved))
		if err != nil {
			return err
		}
		if changes != 1 {
			return budget.ErrReservationConflict
		}
		reservation.Actual = normalized
		reservation.Status = budget.StatusCommitted
		reservation.Overage = overage
		reservation.UpdatedAt = now
		result = reservation
		return nil
	})
	return result, err
}

func (s *Store) Release(ctx context.Context, reservationID record.ReservationID) (budget.Reservation, error) {
	if err := record.ValidateID("reservation", string(reservationID)); err != nil {
		return budget.Reservation{}, err
	}
	var result budget.Reservation
	err := s.db.InTx(ctx, func(conn *sqlitedb.Conn) error {
		reservation, found, err := reservationByIDConn(ctx, conn, reservationID)
		if err != nil {
			return err
		}
		if !found {
			return budget.ErrReservationNotFound
		}
		if reservation.Status == budget.StatusReleased {
			result = reservation
			return nil
		}
		if reservation.Status == budget.StatusCommitted {
			return budget.ErrReservationTerminal
		}
		now := s.now().UTC()
		changes, err := conn.Exec(ctx, `UPDATE budget_reservations SET status = ?, updated_at = ? WHERE id = ? AND status = ?`,
			string(budget.StatusReleased), formatTime(now), string(reservationID), string(budget.StatusReserved))
		if err != nil {
			return err
		}
		if changes != 1 {
			return budget.ErrReservationConflict
		}
		reservation.Status = budget.StatusReleased
		reservation.UpdatedAt = now
		result = reservation
		return nil
	})
	return result, err
}

func (s *Store) GetReservation(ctx context.Context, reservationID record.ReservationID) (budget.Reservation, error) {
	rows, err := s.db.Query(ctx, budgetReservationSelect+` WHERE id = ?`, string(reservationID))
	if err != nil {
		return budget.Reservation{}, err
	}
	if len(rows) == 0 {
		return budget.Reservation{}, budget.ErrReservationNotFound
	}
	return decodeReservation(rows[0])
}

func (s *Store) ReservationForWork(ctx context.Context, workID record.WorkID) (budget.Reservation, error) {
	rows, err := s.db.Query(ctx, budgetReservationSelect+` WHERE work_id = ?`, string(workID))
	if err != nil {
		return budget.Reservation{}, err
	}
	if len(rows) == 0 {
		return budget.Reservation{}, budget.ErrReservationNotFound
	}
	return decodeReservation(rows[0])
}

func (s *Store) Snapshot(ctx context.Context, campaignID record.CampaignID) (budget.Snapshot, error) {
	if err := record.ValidateID("campaign", string(campaignID)); err != nil {
		return budget.Snapshot{}, err
	}
	var result budget.Snapshot
	err := s.db.InTx(ctx, func(conn *sqlitedb.Conn) error {
		value, err := snapshotConn(ctx, conn, campaignID)
		if err != nil {
			return err
		}
		result = value
		return nil
	})
	return result, err
}

func snapshotConn(ctx context.Context, conn *sqlitedb.Conn, campaignID record.CampaignID) (budget.Snapshot, error) {
	limitsRows, err := conn.Query(ctx, `SELECT resource, limit_units FROM budget_limits WHERE campaign_id = ? ORDER BY resource`, string(campaignID))
	if err != nil {
		return budget.Snapshot{}, err
	}
	if len(limitsRows) == 0 {
		return budget.Snapshot{}, budget.ErrNotDefined
	}
	resources := make([]budget.ResourceSnapshot, 0, len(limitsRows))
	index := make(map[budget.Resource]int)
	for _, row := range limitsRows {
		resourceRaw, err := text(row, "resource")
		if err != nil {
			return budget.Snapshot{}, err
		}
		units, err := integer(row, "limit_units")
		if err != nil {
			return budget.Snapshot{}, err
		}
		index[budget.Resource(resourceRaw)] = len(resources)
		resources = append(resources, budget.ResourceSnapshot{Resource: budget.Resource(resourceRaw), Limit: units, Available: units})
	}
	rows, err := conn.Query(ctx, budgetReservationSelect+` WHERE campaign_id = ?`, string(campaignID))
	if err != nil {
		return budget.Snapshot{}, err
	}
	violated := false
	for _, row := range rows {
		reservation, err := decodeReservation(row)
		if err != nil {
			return budget.Snapshot{}, err
		}
		violated = violated || reservation.Overage
		values := reservation.Requested
		committed := false
		switch reservation.Status {
		case budget.StatusReserved:
		case budget.StatusCommitted:
			values = reservation.Actual
			committed = true
		case budget.StatusReleased:
			continue
		default:
			return budget.Snapshot{}, fmt.Errorf("unknown reservation status %q", reservation.Status)
		}
		for _, value := range values {
			position, ok := index[value.Resource]
			if !ok {
				return budget.Snapshot{}, fmt.Errorf("reservation references undefined resource %s", value.Resource)
			}
			if committed {
				resources[position].Committed += value.Units
			} else {
				resources[position].Reserved += value.Units
			}
		}
	}
	for index := range resources {
		resources[index].Available = resources[index].Limit - resources[index].Reserved - resources[index].Committed
		if resources[index].Available < 0 {
			violated = true
		}
	}
	return budget.Snapshot{Campaign: campaignID, Resources: resources, Violated: violated}, nil
}

func reservationByIDConn(ctx context.Context, conn *sqlitedb.Conn, id record.ReservationID) (budget.Reservation, bool, error) {
	rows, err := conn.Query(ctx, budgetReservationSelect+` WHERE id = ?`, string(id))
	if err != nil {
		return budget.Reservation{}, false, err
	}
	if len(rows) == 0 {
		return budget.Reservation{}, false, nil
	}
	value, err := decodeReservation(rows[0])
	return value, true, err
}

func reservationByWorkConn(ctx context.Context, conn *sqlitedb.Conn, id record.WorkID) (budget.Reservation, bool, error) {
	rows, err := conn.Query(ctx, budgetReservationSelect+` WHERE work_id = ?`, string(id))
	if err != nil {
		return budget.Reservation{}, false, err
	}
	if len(rows) == 0 {
		return budget.Reservation{}, false, nil
	}
	value, err := decodeReservation(rows[0])
	return value, true, err
}

func loadLimitsConn(ctx context.Context, conn *sqlitedb.Conn, campaignID record.CampaignID) ([]budget.Limit, error) {
	rows, err := conn.Query(ctx, `SELECT resource, limit_units FROM budget_limits WHERE campaign_id = ? ORDER BY resource`, string(campaignID))
	if err != nil {
		return nil, err
	}
	limits := make([]budget.Limit, 0, len(rows))
	for _, row := range rows {
		resource, err := text(row, "resource")
		if err != nil {
			return nil, err
		}
		units, err := integer(row, "limit_units")
		if err != nil {
			return nil, err
		}
		limits = append(limits, budget.Limit{Resource: budget.Resource(resource), Units: units})
	}
	return limits, nil
}

func usageConn(ctx context.Context, conn *sqlitedb.Conn, campaignID record.CampaignID, exclude record.ReservationID) (map[budget.Resource]int64, bool, error) {
	rows, err := conn.Query(ctx, budgetReservationSelect+` WHERE campaign_id = ?`, string(campaignID))
	if err != nil {
		return nil, false, err
	}
	usage := make(map[budget.Resource]int64)
	overage := false
	for _, row := range rows {
		reservation, err := decodeReservation(row)
		if err != nil {
			return nil, false, err
		}
		if exclude != "" && reservation.ID == exclude {
			continue
		}
		overage = overage || reservation.Overage
		switch reservation.Status {
		case budget.StatusReserved:
			for _, value := range reservation.Requested {
				usage[value.Resource] += value.Units
			}
		case budget.StatusCommitted:
			for _, value := range reservation.Actual {
				usage[value.Resource] += value.Units
			}
		case budget.StatusReleased:
		default:
			return nil, false, fmt.Errorf("unknown reservation status %q", reservation.Status)
		}
	}
	return usage, overage, nil
}

func decodeReservation(row sqlitedb.Row) (budget.Reservation, error) {
	id, err := text(row, "id")
	if err != nil {
		return budget.Reservation{}, err
	}
	campaignID, err := text(row, "campaign_id")
	if err != nil {
		return budget.Reservation{}, err
	}
	workID, err := text(row, "work_id")
	if err != nil {
		return budget.Reservation{}, err
	}
	requestedJSON, err := bytesValue(row, "requested_json")
	if err != nil {
		return budget.Reservation{}, err
	}
	actualJSON, err := bytesValue(row, "actual_json")
	if err != nil {
		return budget.Reservation{}, err
	}
	status, err := text(row, "status")
	if err != nil {
		return budget.Reservation{}, err
	}
	overageValue, err := integer(row, "overage")
	if err != nil {
		return budget.Reservation{}, err
	}
	createdRaw, err := text(row, "created_at")
	if err != nil {
		return budget.Reservation{}, err
	}
	updatedRaw, err := text(row, "updated_at")
	if err != nil {
		return budget.Reservation{}, err
	}
	createdAt, err := parseTime(createdRaw)
	if err != nil {
		return budget.Reservation{}, err
	}
	updatedAt, err := parseTime(updatedRaw)
	if err != nil {
		return budget.Reservation{}, err
	}
	var requested []budget.Quantity
	if err := json.Unmarshal(requestedJSON, &requested); err != nil {
		return budget.Reservation{}, err
	}
	var actual []budget.Quantity
	if err := json.Unmarshal(actualJSON, &actual); err != nil {
		return budget.Reservation{}, err
	}
	return budget.Reservation{
		ID: record.ReservationID(id), Campaign: record.CampaignID(campaignID), Work: record.WorkID(workID),
		Requested: requested, Actual: actual, Status: budget.ReservationStatus(status),
		Overage: overageValue != 0, CreatedAt: createdAt, UpdatedAt: updatedAt,
	}, nil
}

const budgetReservationSelect = `SELECT id, campaign_id, work_id, requested_json, actual_json, status, overage, created_at, updated_at FROM budget_reservations`

func quantityMap(values []budget.Quantity) map[budget.Resource]int64 {
	out := make(map[budget.Resource]int64, len(values))
	for _, value := range values {
		out[value.Resource] = value.Units
	}
	return out
}

func limitMap(values []budget.Limit) map[budget.Resource]int64 {
	out := make(map[budget.Resource]int64, len(values))
	for _, value := range values {
		out[value.Resource] = value.Units
	}
	return out
}

func quantitiesEqual(left, right []budget.Quantity) bool {
	leftJSON, leftErr := record.CanonicalJSON(left)
	rightJSON, rightErr := record.CanonicalJSON(right)
	return leftErr == nil && rightErr == nil && bytes.Equal(leftJSON, rightJSON)
}

func limitsEqual(left, right []budget.Limit) bool {
	left = append([]budget.Limit(nil), left...)
	right = append([]budget.Limit(nil), right...)
	sort.Slice(left, func(i, j int) bool { return left[i].Resource < left[j].Resource })
	sort.Slice(right, func(i, j int) bool { return right[i].Resource < right[j].Resource })
	leftJSON, leftErr := record.CanonicalJSON(left)
	rightJSON, rightErr := record.CanonicalJSON(right)
	return leftErr == nil && rightErr == nil && bytes.Equal(leftJSON, rightJSON)
}

func boolInt(value bool) int64 {
	if value {
		return 1
	}
	return 0
}

var _ budget.Ledger = (*Store)(nil)
