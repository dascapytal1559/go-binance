package pmargin

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/adshao/go-binance/v2/futures"
)

// Order define order info
type UMQueryOrder struct {
	AveragePrice            string                   `json:"avgPrice"`
	ClientOrderID           string                   `json:"clientOrderId"`
	CumulativeQuote         string                   `json:"cumQuote"`
	ExecutedQuantity        string                   `json:"executedQty"`
	OrderID                 int64                    `json:"orderId"`
	OriginalQuantity        string                   `json:"origQty"`
	OriginalType            futures.OrderType        `json:"origType"`
	Price                   string                   `json:"price"`
	ReduceOnly              bool                     `json:"reduceOnly"`
	Side                    futures.SideType         `json:"side"`
	PositionSide            futures.PositionSideType `json:"positionSide"`
	Status                  futures.OrderStatusType  `json:"status"`
	Symbol                  string                   `json:"symbol"`
	Time                    int64                    `json:"time"`
	TimeInForce             futures.TimeInForceType  `json:"timeInForce"`
	Type                    futures.OrderType        `json:"type"`
	UpdateTime              int64                    `json:"updateTime"`
	SelfTradePreventionMode string                   `json:"selfTradePreventionMode"`
	GoodTillDate            int64                    `json:"goodTillDate"`
}

// ListUMOpenOrdersService list opened orders
type ListUMOpenOrdersService struct {
	c      *Client
	symbol string
}

// Symbol set symbol
func (s *ListUMOpenOrdersService) Symbol(symbol string) *ListUMOpenOrdersService {
	s.symbol = symbol
	return s
}

// Do send request
func (s *ListUMOpenOrdersService) Do(ctx context.Context, opts ...RequestOption) (res []*UMQueryOrder, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/papi/v1/um/openOrders",
		secType:  secTypeSigned,
	}
	if s.symbol != "" {
		r.setParam("symbol", s.symbol)
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []*UMQueryOrder{}, err
	}
	res = make([]*UMQueryOrder, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*UMQueryOrder{}, err
	}
	return res, nil
}
