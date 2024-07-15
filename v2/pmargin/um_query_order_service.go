package pmargin

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/adshao/go-binance/v2/futures"
)

// ListUMOpenOrdersService list opened orders
type ListUMOpenOrdersService struct {
	c      *Client
	symbol *string
}

// Symbol set symbol
func (s *ListUMOpenOrdersService) Symbol(symbol string) *ListUMOpenOrdersService {
	s.symbol = &symbol
	return s
}

// Do send request
func (s *ListUMOpenOrdersService) Do(ctx context.Context, opts ...RequestOption) (res []*UMOpenOrder, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/papi/v1/um/openOrders",
		secType:  secTypeSigned,
	}

	r.setParams(params{
		"symbol": s.symbol,
	})

	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []*UMOpenOrder{}, err
	}
	res = make([]*UMOpenOrder, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*UMOpenOrder{}, err
	}
	return res, nil
}

// OrderResponse define order info
type UMOpenOrder struct {
	AvgPrice                string                   `json:"avgPrice"`
	ClientOrderId           string                   `json:"clientOrderId"`
	CumQuote                string                   `json:"cumQuote"`
	ExecutedQty             string                   `json:"executedQty"`
	OrderId                 int64                    `json:"orderId"`
	OrigQty                 string                   `json:"origQty"`
	OrigType                futures.OrderType        `json:"origType"`
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
