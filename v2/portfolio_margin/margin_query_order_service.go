package portfolio_margin

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/adshao/go-binance/v2"
)

type MarginQueryOrder struct {
	ClientOrderID           string                  `json:"clientOrderId"`
	ExecutedQuantity        string                  `json:"executedQty"`
	CumulativeQuoteQuantity string                  `json:"cummulativeQuoteQty"`
	IcebergQuantity         string                  `json:"icebergQty"`
	IsWorking               bool                    `json:"isWorking"`
	OrderID                 int64                   `json:"orderId"`
	OriginalQuantity        string                  `json:"origQty"`
	Price                   string                  `json:"price"`
	Side                    binance.SideType        `json:"side"`
	Status                  binance.OrderStatusType `json:"status"`
	StopPrice               string                  `json:"stopPrice"`
	Symbol                  string                  `json:"symbol"`
	Time                    int64                   `json:"time"`
	TimeInForce             binance.TimeInForceType `json:"timeInForce"`
	Type                    binance.OrderType       `json:"type"`
	UpdateTime              int64                   `json:"updateTime"`
	AccountId               int64                   `json:"accountId"`
	SelfTradePreventionMode string                  `json:"selfTradePreventionMode"`
	PreventedMatchId        int64                   `json:"preventedMatchId"`
	PreventedQuantity       string                  `json:"preventedQuantity"`
}

// ListMarginOpenOrdersService list opened orders
type ListMarginOpenOrdersService struct {
	c      *Client
	symbol string
}

// Symbol set symbol
func (s *ListMarginOpenOrdersService) Symbol(symbol string) *ListMarginOpenOrdersService {
	s.symbol = symbol
	return s
}

// Do send request
func (s *ListMarginOpenOrdersService) Do(ctx context.Context, opts ...RequestOption) (res []*MarginQueryOrder, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/papi/v1/margin/openOrders",
		secType:  secTypeSigned,
	}
	if s.symbol != "" {
		r.setParam("symbol", s.symbol)
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []*MarginQueryOrder{}, err
	}
	res = make([]*MarginQueryOrder, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*MarginQueryOrder{}, err
	}
	return res, nil
}
