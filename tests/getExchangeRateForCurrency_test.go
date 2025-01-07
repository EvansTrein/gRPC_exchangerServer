package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/EvansTrein/gRPC_exchangerServer/internal/server"
	"github.com/EvansTrein/gRPC_exchangerServer/internal/storages"
	pb "github.com/EvansTrein/proto-exchange/exchange"
)

type testCase struct {
	name           string
	fromCurrency   string
	toCurrency     string
	expectedError  error
}

func TestGetExchangeRateForCurrency(t *testing.T) {
	ctx, suite := NewAppTest(t)

	testCases := []testCase{
		{
			name:          "valid USD to EUR",
			fromCurrency:  "USD",
			toCurrency:    "EUR",
			expectedError: nil,
		},
		{
			name:          "valid EUR to USD",
			fromCurrency:  "EUR",
			toCurrency:    "USD",
			expectedError: nil,
		},
		{
			name:          "valid CNY to RUB",
			fromCurrency:  "CNY",
			toCurrency:    "RUB",
			expectedError: nil,
		},
		{
			name:          "valid USD to RUB",
			fromCurrency:  "USD",
			toCurrency:    "RUB",
			expectedError: nil,
		},
		{
			name:          "valid EUR to RUB",
			fromCurrency:  "EUR",
			toCurrency:    "RUB",
			expectedError: nil,
		},
		{
			name:          "invalid no single currency",
			fromCurrency:  "USD",
			toCurrency:    "AED",
			expectedError: status.Error(codes.NotFound, storages.ErrExchangeRateNotFound.Error()),
		},
		{
			name:          "invalid both currencies no",
			fromCurrency:  "CAD",
			toCurrency:    "AED",
			expectedError: status.Error(codes.NotFound, storages.ErrExchangeRateNotFound.Error()),
		},
		{
			name:          "invalid one currency not transferred",
			fromCurrency:  "USD",
			toCurrency:    "",
			expectedError: status.Error(codes.InvalidArgument, server.ErrinValidData.Error()),
		},
		{
			name:          "invalid both currencies have not been transferred",
			fromCurrency:  "",
			toCurrency:    "",
			expectedError: status.Error(codes.InvalidArgument, server.ErrinValidData.Error()),
		},
		{
			name:          "invalid same currencies",
			fromCurrency:  "USD",
			toCurrency:    "USD",
			expectedError: status.Error(codes.InvalidArgument, server.ErrinValidData.Error()),
		},
		{
			name:          "invalid non-currency format currency code",
			fromCurrency:  "USDUSDUSD",
			toCurrency:    "RUB",
			expectedError: status.Error(codes.InvalidArgument, server.ErrinValidData.Error()),
		},
		{
			name:          "invalid non-currency format currency code",
			fromCurrency:  "R",
			toCurrency:    "EUR",
			expectedError: status.Error(codes.InvalidArgument, server.ErrinValidData.Error()),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := suite.ExchangeClient.GetExchangeRateForCurrency(ctx, &pb.CurrencyRequest{
				FromCurrency: tc.fromCurrency,
				ToCurrency:   tc.toCurrency,
			})

			if tc.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.fromCurrency, resp.FromCurrency)
				assert.Equal(t, tc.toCurrency, resp.ToCurrency)
				assert.Greater(t, resp.Rate, float32(0))
			}
		})
	}
}
