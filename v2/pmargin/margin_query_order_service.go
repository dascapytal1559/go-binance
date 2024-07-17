package pmargin

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/adshao/go-binance/v2"
)

// NewListMarginOpenOrdersService init list margin open orders service
func (c *Client) NewListMarginOpenOrdersService() *ListMarginOpenOrdersService {
	return &ListMarginOpenOrdersService{c: c}
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
func (s *ListMarginOpenOrdersService) Do(ctx context.Context, opts ...RequestOption) (res []*MarginQueryOrderResponse, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/papi/v1/margin/openOrders",
		secType:  secTypeSigned,
	}
	r.setParam("symbol", s.symbol)
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = make([]*MarginQueryOrderResponse, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type MarginQueryOrderResponse struct {
	ClientOrderID           string                  `json:"clientOrderId"`
	CumulativeQuoteQuantity string                  `json:"cummulativeQuoteQty"`
	ExecutedQuantity        string                  `json:"executedQty"`
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
