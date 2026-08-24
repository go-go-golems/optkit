package budget

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"time"

	"github.com/go-go-golems/optkit/record"
)

var (
	ErrNotDefined          = errors.New("campaign budget is not defined")
	ErrLimitConflict       = errors.New("campaign budget limits conflict")
	ErrInsufficient        = errors.New("insufficient campaign budget")
	ErrReservationConflict = errors.New("budget reservation conflict")
	ErrReservationNotFound = errors.New("budget reservation not found")
	ErrReservationTerminal = errors.New("budget reservation is terminal")
)

type Resource string

var resourcePattern = regexp.MustCompile(`^[a-z][a-z0-9._-]*$`)

func (r Resource) Validate() error {
	if !resourcePattern.MatchString(string(r)) {
		return fmt.Errorf("invalid resource %q", r)
	}
	return nil
}

type Quantity struct {
	Resource Resource `json:"resource"`
	Units    int64    `json:"units"`
}

type Limit struct {
	Resource Resource `json:"resource"`
	Units    int64    `json:"units"`
}

type ReservationStatus string

const (
	StatusReserved  ReservationStatus = "reserved"
	StatusCommitted ReservationStatus = "committed"
	StatusReleased  ReservationStatus = "released"
)

type Reservation struct {
	ID        record.ReservationID `json:"id"`
	Campaign  record.CampaignID    `json:"campaign"`
	Work      record.WorkID        `json:"work"`
	Requested []Quantity           `json:"requested"`
	Actual    []Quantity           `json:"actual,omitempty"`
	Status    ReservationStatus    `json:"status"`
	Overage   bool                 `json:"overage"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}

type ReserveRequest struct {
	Campaign record.CampaignID `json:"campaign"`
	Work     record.WorkID     `json:"work"`
	Claims   []Quantity        `json:"claims"`
}

type ResourceSnapshot struct {
	Resource  Resource `json:"resource"`
	Limit     int64    `json:"limit"`
	Reserved  int64    `json:"reserved"`
	Committed int64    `json:"committed"`
	Available int64    `json:"available"`
}

type Snapshot struct {
	Campaign  record.CampaignID  `json:"campaign"`
	Resources []ResourceSnapshot `json:"resources"`
	Violated  bool               `json:"violated"`
}

type Ledger interface {
	Define(context.Context, record.CampaignID, []Limit) error
	Reserve(context.Context, ReserveRequest) (Reservation, error)
	Commit(context.Context, record.ReservationID, []Quantity) (Reservation, error)
	Release(context.Context, record.ReservationID) (Reservation, error)
	GetReservation(context.Context, record.ReservationID) (Reservation, error)
	ReservationForWork(context.Context, record.WorkID) (Reservation, error)
	Snapshot(context.Context, record.CampaignID) (Snapshot, error)
}

func NormalizeLimits(limits []Limit) ([]Limit, error) {
	if len(limits) == 0 {
		return nil, fmt.Errorf("at least one budget limit is required")
	}
	out := append([]Limit(nil), limits...)
	sort.Slice(out, func(i, j int) bool { return out[i].Resource < out[j].Resource })
	for index, limit := range out {
		if err := limit.Resource.Validate(); err != nil {
			return nil, err
		}
		if limit.Units < 0 {
			return nil, fmt.Errorf("budget limit %s cannot be negative", limit.Resource)
		}
		if index > 0 && out[index-1].Resource == limit.Resource {
			return nil, fmt.Errorf("duplicate budget limit %s", limit.Resource)
		}
	}
	return out, nil
}

func NormalizeQuantities(values []Quantity, allowEmpty bool) ([]Quantity, error) {
	if len(values) == 0 && !allowEmpty {
		return nil, fmt.Errorf("at least one resource quantity is required")
	}
	out := append([]Quantity(nil), values...)
	sort.Slice(out, func(i, j int) bool { return out[i].Resource < out[j].Resource })
	for index, value := range out {
		if err := value.Resource.Validate(); err != nil {
			return nil, err
		}
		if value.Units < 0 {
			return nil, fmt.Errorf("resource quantity %s cannot be negative", value.Resource)
		}
		if index > 0 && out[index-1].Resource == value.Resource {
			return nil, fmt.Errorf("duplicate resource quantity %s", value.Resource)
		}
	}
	return out, nil
}

func ReservationID(request ReserveRequest) (record.ReservationID, error) {
	if err := record.ValidateID("campaign", string(request.Campaign)); err != nil {
		return "", err
	}
	if err := record.ValidateID("work", string(request.Work)); err != nil {
		return "", err
	}
	claims, err := NormalizeQuantities(request.Claims, false)
	if err != nil {
		return "", err
	}
	digest, _, err := record.SemanticDigest("schema:optkit.budget-reservation-identity/v1", struct {
		Campaign record.CampaignID `json:"campaign"`
		Work     record.WorkID     `json:"work"`
		Claims   []Quantity        `json:"claims"`
	}{request.Campaign, request.Work, claims})
	if err != nil {
		return "", err
	}
	raw, err := record.ContentID("reservation", digest)
	if err != nil {
		return "", err
	}
	return record.ReservationID(raw), nil
}
