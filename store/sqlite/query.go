package sqlite

import (
	"context"
	"sort"

	"github.com/go-go-golems/optkit/record"
)

// ListCampaigns returns all locally known campaign IDs in stable order. It is
// a read-only query-plane capability; campaign creation remains journal-owned.
func (s *Store) ListCampaigns(ctx context.Context) ([]record.CampaignID, error) {
	rows, err := s.db.Query(ctx, `SELECT campaign_id FROM campaign_heads ORDER BY campaign_id`)
	if err != nil {
		return nil, err
	}
	ids := make([]record.CampaignID, 0, len(rows))
	for _, row := range rows {
		raw, err := text(row, "campaign_id")
		if err != nil {
			return nil, err
		}
		ids = append(ids, record.CampaignID(raw))
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	return ids, nil
}
