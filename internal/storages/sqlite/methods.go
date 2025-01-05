package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EvansTrein/exchanger_gRPC/internal/storages"
)

func (s *SQLiteDB) Rate(ctx context.Context, currency string) (*storages.Rate, error) {
	var rate storages.Rate

	s.log.Info("Rate method used")

	err := s.db.QueryRowContext(ctx, "SELECT name FROM Currencies WHERE code = ?", currency).Scan(&rate.Currency)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("rate not found for currency: %s", currency)
		}
		return nil, fmt.Errorf("failed to query rate: %w", err)
	}

	rate.Value = 0.71

	return &rate, nil
}

func (s *SQLiteDB) Rates(ctx context.Context) ([]storages.Rate, error) {

	return []storages.Rate{}, nil
}
