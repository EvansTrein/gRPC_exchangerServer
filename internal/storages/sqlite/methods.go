package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/EvansTrein/gRPC_exchangerServer/internal/storages"
	currencyrateapi "github.com/EvansTrein/gRPC_exchangerServer/pkg/currencyRateAPI"
	"github.com/EvansTrein/gRPC_exchangerServer/pkg/utils"
)

func (s *SQLiteDB) Rate(ctx context.Context, fromCurrency, toCurrency string) (*storages.Rate, error) {
	const op = "func Rate"
	log := s.log.With(
		slog.String("operation", op),
		slog.Any("calling context", ctx),
		slog.String("function argument fromCurrency", fromCurrency),
		slog.String("function argument toCurrency", toCurrency),
	)

	log.Debug("call of the Rate SQL method")

	var rate storages.Rate

	query := `
		SELECT 
			BaseCurrency.code AS baseCurrencyCode,
			ToCurrency.code AS toCurrencyCode,
			Rates.rate
		FROM 
			Rates
		JOIN 
			Currencies AS BaseCurrency ON Rates.baseCurrencyID = BaseCurrency.id
		JOIN 
			Currencies AS ToCurrency ON Rates.toCurrencyID = ToCurrency.id
		WHERE 
			BaseCurrency.code = ? AND ToCurrency.code = ?
		`

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("failed to prepare SQL query", "error", err)
		return nil, fmt.Errorf("failed to prepare SQL query: %v", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, fromCurrency, toCurrency)

	if err = row.Scan(&rate.BaseCurrency, &rate.ToCurrency, &rate.Rate); err != nil {
		if err == sql.ErrNoRows {
			return nil, storages.ErrExchangeRateNotFound
		}
		log.Error("failed to execute SQL query", "error", err)
		return nil, err
	}

	return &rate, nil
}

func (s *SQLiteDB) AllRates(ctx context.Context) (map[string]float32, error) {
	const op = "func AllRates"
	log := s.log.With(
		slog.String("operation", op),
		slog.Any("calling context", ctx),
	)
	log.Debug("call of the AllRates SQL method")

	answer := map[string]float32{}
	query := `
		SELECT 
			BaseCurrency.code AS baseCurrencyCode,
			ToCurrency.code AS toCurrencyCode,
			Rates.rate
		FROM 
			Rates
		JOIN 
			Currencies AS BaseCurrency ON Rates.baseCurrencyID = BaseCurrency.id
		JOIN 
			Currencies AS ToCurrency ON Rates.toCurrencyID = ToCurrency.id;`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		log.Error("failed to execute SQL query", "error", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var baseCurrencyCode, toCurrencyCode string
		var rate float32

		if err := rows.Scan(&baseCurrencyCode, &toCurrencyCode, &rate); err != nil {
			log.Error("failed to scan data received from the database", "error", err)
			return nil, err
		}

		key := fmt.Sprintf("%s/%s", baseCurrencyCode, toCurrencyCode)

		answer[key] = rate
	}

	if err := rows.Err(); err != nil {
		log.Error("error when iterating over rows", "error", err)
		return nil, err
	}

	return answer, nil
}

func (s *SQLiteDB) RatesDownloadFromExternalAPI() error {
	op := "func RatesDownloadFromExternalAPI"
	log := s.log.With(slog.String("operation", op))

	log.Debug("call of the RatesDownloadFromExternalAPI SQL method")

	data, err := s.allCurrencies()
	if err != nil {
		log.Error("failed to get all currencies from the database", "error", err)
		return err
	}

	log.Debug("all currencies were successfully retrieved from the database", slog.Any("data", data))

	pairs, err := utils.GenerateCurrencyPairs(data)
	if err != nil {
		log.Error("failed to get currency pairs", "error", err)
		return err
	}

	log.Debug("currency pairs for API request successfully received", slog.Any("data", pairs))

	var rates []storages.Rate
	for baseCurrencie, toCurrencies := range pairs {
		dataFromApi, err := currencyrateapi.DownloadExchangeRateData(baseCurrencie, toCurrencies)
		if err != nil {
			log.Error("failed to get data from third-party API", "error", err)
			return err
		}

		if len(dataFromApi) == 0 {
			log.Error("the request was successful, but no data", "error", err)
			return err
		}

		rates = append(rates, dataFromApi...)
	}

	log.Debug("data was successfully received from the API", "data", rates, slog.Int("number of exchange rates", len(rates)))

	today := time.Now().Format(time.DateTime)
	query := `INSERT INTO Rates (baseCurrencyID, toCurrencyID, rate, date)
				VALUES
				((SELECT id FROM Currencies WHERE code = ?), (SELECT id FROM Currencies WHERE code = ?), ?, ?);`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		log.Error("failed to prepare SQL query", "error", err)
		return err
	}
	defer stmt.Close()

	for _, el := range rates {
		if _, err := stmt.Exec(el.BaseCurrency, el.ToCurrency, el.Rate, today); err != nil {
			log.Error("SQL query failed", "error", err)
			return err
		}
	}

	return nil
}

func (s *SQLiteDB) LoadDefaultRates() error {
	op := "func LoadDefaultRates"
	log := s.log.With(slog.String("operation", op))

	log.Debug("call of the RatesDownloadFromExternalAPI SQL method")

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
			s.log.Error("failed to execute sql query", "error", err)
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

func (s *SQLiteDB) allCurrencies() ([]string, error) {
	op := "func allCurrencies"
	log := s.log.With(slog.String("operation", op))

	log.Debug("call of the allCurrencies SQL method")

	var currencyCodes []string

	rows, err := s.db.Query("SELECT code FROM Currencies")
	if err != nil {
		log.Error("failed to execute sql query", "error", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			log.Error("failed to scan data received from the database", "error", err)
			return nil, err
		}
		currencyCodes = append(currencyCodes, code)
	}

	if err := rows.Err(); err != nil {
		log.Error("error when iterating over rows", "error", err)
		return nil, err
	}

	return currencyCodes, nil
}
