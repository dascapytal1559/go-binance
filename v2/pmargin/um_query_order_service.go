package pmargin

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/adshao/go-binance/v2/futures"
)

// NewQueryUMOrderService init querying UM order service
func (c *Client) NewQueryUMOrderService() *QueryUMOrderService {
	return &QueryUMOrderService{c: c}
}

// QueryUMOrder list opened orders
type QueryUMOrderService struct {
	c                 *Client
	symbol            string
	orderId           *int64
	origClientOrderId *string
}

// Symbol set symbol
func (s *QueryUMOrderService) Symbol(symbol string) *QueryUMOrderService {
	s.symbol = symbol
	return s
}

// OrderID set orderID
func (s *QueryUMOrderService) OrderID(orderId int64) *QueryUMOrderService {
	s.orderId = &orderId
	return s
}

// OrigClientOrderID set origClientOrderID
func (s *QueryUMOrderService) OrigClientOrderID(origClientOrderId string) *QueryUMOrderService {
	s.origClientOrderId = &origClientOrderId
	return s
}

// Do send request
func (s *QueryUMOrderService) Do(ctx context.Context, opts ...RequestOption) (res *UMQueryOrderResponse, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/papi/v1/um/order",
		secType:  secTypeSigned,
	}
	r.setParams(params{
		"symbol":            s.symbol,
		"orderId":           s.orderId,
		"origClientOrderId": s.origClientOrderId,
	})
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// OrderResponse define order info
type UMQueryOrderResponse struct {
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

// NewListUMOpenOrdersService init list UM open orders service
func (c *Client) NewListUMOpenOrdersService() *ListUMOpenOrdersService {
	return &ListUMOpenOrdersService{c: c}
}

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
func (s *ListUMOpenOrdersService) Do(ctx context.Context, opts ...RequestOption) (res []*UMQueryOrderResponse, err error) {
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
		return nil, err
	}
	res = make([]*UMQueryOrderResponse, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
