package pmargin

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var client = NewClient(
	os.Getenv("TEST_BINANCE_API_KEY"),
	os.Getenv("TEST_BINANCE_SECRET_KEY"),
)

func TestTradeServiceHistorical(t *testing.T) {
	trades, err := client.NewHistoricalTradesService().FromID(0).Do(context.TODO())

	assert.Nil(t, err)
	assert.NotNil(t, trades)
}
