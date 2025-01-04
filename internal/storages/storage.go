package storages

import "context"

type Rate struct {
	Currency string
	Value    float32
}

type Database interface {
	Rates(ctx context.Context) ([]Rate, error)
	Rate(ctx context.Context, currency string) (*Rate, error)
}
