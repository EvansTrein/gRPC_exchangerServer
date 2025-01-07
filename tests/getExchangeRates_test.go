package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pb "github.com/EvansTrein/proto-exchange/exchange"
)

func TestGetExchangeRates(t *testing.T) {
	ctx, suite := NewAppTest(t)

	resp, err := suite.ExchangeClient.GetExchangeRates(ctx, &pb.Empty{})
	
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Greater(t, len(resp.Rates), 0)
}
