package sqlite

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/EvansTrein/exchanger_gRPC/internal/storages"
)

func (s *SQLiteDB) Rate(ctx context.Context, currency string) (*storages.Rate, error) {
	var rate storages.Rate

	return &rate, nil
}

func (s *SQLiteDB) AllRates(ctx context.Context) (map[string]float32, error) {
	answer := map[string]float32{}

	return answer, nil
}

func (s *SQLiteDB) RatesDownloadFromExternalAPI(TableNameForCurrencyRates string) error {

	return fmt.Errorf("error no func RatesDownloadFromExternalAPI")
}

func (s *SQLiteDB) LoadDefaultRates(TableNameForCurrencyRates string) error {
	queries := make([]string, 0, 4)
	today := time.Now().Format(time.DateTime)

	rubQuery := fmt.Sprintf(
        `INSERT INTO Rates (baseCurrencyID, toCurrencyID, rate, date)
        VALUES
        ((SELECT id FROM Currencies WHERE code = 'RUB'), (SELECT id FROM Currencies WHERE code = 'EUR'), 0.01, '%s'),
        ((SELECT id FROM Currencies WHERE code = 'RUB'), (SELECT id FROM Currencies WHERE code = 'USD'), 0.012, '%s'),
        ((SELECT id FROM Currencies WHERE code = 'RUB'), (SELECT id FROM Currencies WHERE code = 'CNY'), 0.08, '%s');`,
        today, today, today,
    )

	eurQuery := fmt.Sprintf(
        `INSERT INTO Rates (baseCurrencyID, toCurrencyID, rate, date)
        VALUES
        ((SELECT id FROM Currencies WHERE code = 'EUR'), (SELECT id FROM Currencies WHERE code = 'RUB'), 100.0, '%s'),
        ((SELECT id FROM Currencies WHERE code = 'EUR'), (SELECT id FROM Currencies WHERE code = 'USD'), 1.05, '%s'),
        ((SELECT id FROM Currencies WHERE code = 'EUR'), (SELECT id FROM Currencies WHERE code = 'CNY'), 7.8, '%s');`,
        today, today, today,
    )

	usdQuery := fmt.Sprintf(
        `INSERT INTO Rates (baseCurrencyID, toCurrencyID, rate, date)
        VALUES
        ((SELECT id FROM Currencies WHERE code = 'USD'), (SELECT id FROM Currencies WHERE code = 'RUB'), 85.00, '%s'),
        ((SELECT id FROM Currencies WHERE code = 'USD'), (SELECT id FROM Currencies WHERE code = 'EUR'), 0.95, '%s'),
        ((SELECT id FROM Currencies WHERE code = 'USD'), (SELECT id FROM Currencies WHERE code = 'CNY'), 7.2, '%s');`,
        today, today, today,
    )

	cnyQuery := fmt.Sprintf(
        `INSERT INTO Rates (baseCurrencyID, toCurrencyID, rate, date)
        VALUES
        ((SELECT id FROM Currencies WHERE code = 'CNY'), (SELECT id FROM Currencies WHERE code = 'RUB'), 12.5, '%s'),
        ((SELECT id FROM Currencies WHERE code = 'CNY'), (SELECT id FROM Currencies WHERE code = 'EUR'), 0.13, '%s'),
        ((SELECT id FROM Currencies WHERE code = 'CNY'), (SELECT id FROM Currencies WHERE code = 'USD'), 0.14, '%s');`,
        today, today, today,
    )

	queries = append(queries, rubQuery, eurQuery, usdQuery, cnyQuery)

	for _, query := range queries {
        _, err := s.db.Exec(query)
        if err != nil {
			s.log.Error("failed to execute sql query", slog.String("error", err.Error()))
            return err
        }
    }

	return nil
}

func (s *SQLiteDB) IsTableEmpty(tableName string) (bool, error) {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)

	err := s.db.QueryRow(query).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if the table is empty: %v", err)
	}

	return count == 0, nil
}
