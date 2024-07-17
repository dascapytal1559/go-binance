package pmargin

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/adshao/go-binance/v2/futures"
)

// NewCancelUMOrderService init canceling UM order service
func (c *Client) NewCancelUMOrderService() *CancelUMOrderService {
	return &CancelUMOrderService{c: c}
}

// CancelUMOrderService cancel an order
type CancelUMOrderService struct {
	c                 *Client
	symbol            string
	orderID           *int64
	origClientOrderID *string
}

// Symbol set symbol
func (s *CancelUMOrderService) Symbol(symbol string) *CancelUMOrderService {
	s.symbol = symbol
	return s
}

// OrderID set orderID
func (s *CancelUMOrderService) OrderID(orderID int64) *CancelUMOrderService {
	s.orderID = &orderID
	return s
}

// OrigClientOrderID set origClientOrderID
func (s *CancelUMOrderService) OrigClientOrderID(origClientOrderID string) *CancelUMOrderService {
	s.origClientOrderID = &origClientOrderID
	return s
}

// Do send request
func (s *CancelUMOrderService) Do(ctx context.Context, opts ...RequestOption) (res *CancelUMOrderResponse, err error) {
	r := &request{
		method:   http.MethodDelete,
		endpoint: "/papi/v1/um/order",
		secType:  secTypeSigned,
	}
	r.setFormParams(params{
		"symbol":            s.symbol,
		"orderId":           s.orderID,
		"origClientOrderId": s.origClientOrderID,
	})
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(CancelUMOrderResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// OrderResponse define order info
type CancelUMOrderResponse struct {
	AvgPrice                string                   `json:"avgPrice"`
	ClientOrderId           string                   `json:"clientOrderId"`
	CumQty                  string                   `json:"cumQty"`
	CumQuote                string                   `json:"cumQuote"`
	ExecutedQty             string                   `json:"executedQty"`
	OrderId                 int64                    `json:"orderId"`
	OrigQty                 string                   `json:"origQty"`
	Price                   string                   `json:"price"`
	ReduceOnly              bool                     `json:"reduceOnly"`
	Side                    futures.SideType         `json:"side"`
	PositionSide            futures.PositionSideType `json:"positionSide"`
	Status                  futures.OrderStatusType  `json:"status"`
	Symbol                  string                   `json:"symbol"`
	TimeInForce             futures.TimeInForceType  `json:"timeInForce"`
	Type                    futures.OrderType        `json:"type"`
	UpdateTime              int64                    `json:"updateTime"`
	SelfTradePreventionMode string                   `json:"selfTradePreventionMode"`
	GoodTillDate            int64                    `json:"goodTillDate"`
}
