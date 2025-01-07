package storages

import (
	"context"
	"errors"
)

const TableNameForCurrencyRates = "Rates"

var ErrExchangeRateNotFound = errors.New("exchange rate not found")

type Rate struct {
	BaseCurrency string
	ToCurrency   string
	Rate         float32
}

type Database interface {
	AllRates(ctx context.Context) (map[string]float32, error)
	Rate(ctx context.Context, fromCurrency, toCurrency string) (*Rate, error)
	RatesDownloadFromExternalAPI() error
	LoadDefaultRates() error
	IsTableEmpty(tableName string) (bool, error)
	Close() error
}
