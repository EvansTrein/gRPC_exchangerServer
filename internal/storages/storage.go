package storages

import "context"

const TableNameForCurrencyRates = "Rates"

type Rate struct {
	BaseCurrency string
	ToCurrency   string
	Rate         float32
}

type Database interface {
	AllRates(ctx context.Context) (map[string]float32, error)
	Rate(ctx context.Context, currency string) (*Rate, error)
	RatesDownloadFromExternalAPI(TableNameForCurrencyRates string) error
	LoadDefaultRates(TableNameForCurrencyRates string) error
	IsTableEmpty(tableName string) (bool, error)
	Close() error
}
