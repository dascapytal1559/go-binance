package pmargin

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIncomeHistoryTest(t *testing.T) {
	client := NewClient(
		os.Getenv("TEST_BINANCE_API_KEY"),
		os.Getenv("TEST_BINANCE_SECRET_KEY"),
	)

	trades, err := client.NewIncomeHistoryService().Symbol("1000PEPEUSDT").Do(context.Background())

	for _, t := range trades {
		fmt.Printf("%v, %v\n", t.IncomeType, t.Income)
	}

	assert.Nil(t, err)
	assert.NotNil(t, trades)
}
