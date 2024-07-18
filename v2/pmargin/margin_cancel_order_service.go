package pmargin

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/adshao/go-binance/v2"
)

// NewCancelMarginOrderService init creating UM order service
func (c *Client) NewCancelMarginOrderService() *CancelMarginOrderService {
	return &CancelMarginOrderService{c: c}
}

// CancelMarginOrderService create order
type CancelMarginOrderService struct {
	c                 *Client
	symbol            string
	orderId           *int64
	origClientOrderId *string
	newClientOrderID  *string
}

// Symbol set symbol
func (s *CancelMarginOrderService) Symbol(symbol string) *CancelMarginOrderService {
	s.symbol = symbol
	return s
}

// OrderID set orderID
func (s *CancelMarginOrderService) OrderID(orderId int64) *CancelMarginOrderService {
	s.orderId = &orderId
	return s
}

// OrigClientOrderID set origClientOrderID
func (s *CancelMarginOrderService) OrigClientOrderID(origClientOrderId string) *CancelMarginOrderService {
	s.origClientOrderId = &origClientOrderId
	return s
}

// NewClientOrderID set newClientOrderID
func (s *CancelMarginOrderService) NewClientOrderID(newClientOrderID string) *CancelMarginOrderService {
	s.newClientOrderID = &newClientOrderID
	return s
}

// Do send request
func (s *CancelMarginOrderService) Do(ctx context.Context, opts ...RequestOption) (res *CancelMarginOrderResponse, err error) {
	r := &request{
		method:   http.MethodDelete,
		endpoint: "/papi/v1/margin/order",
		secType:  secTypeSigned,
	}
	r.setFormParams(params{
		"symbol":            s.symbol,
		"orderId":           s.orderId,
		"origClientOrderId": s.origClientOrderId,
		"newClientOrderID":  s.newClientOrderID,
	})
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(CancelMarginOrderResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// OrderResponse define order info
type CancelMarginOrderResponse struct {
	Symbol                  string                  `json:"symbol"`
	OrderId                 int64                   `json:"orderId"`
	OrigClientOrderId       string                  `json:"origClientOrderId"`
	ClientOrderId           string                  `json:"clientOrderId"`
	Price                   string                  `json:"price"`
	OrigQty                 string                  `json:"origQty"`
	ExecutedQty             string                  `json:"executedQty"`
	CummulativeQuoteQty     string                  `json:"cummulativeQuoteQty"`
	Status                  binance.OrderStatusType `json:"status"`
	TimeInForce             binance.TimeInForceType `json:"timeInForce"`
	Type                    binance.OrderType       `json:"type"`
	Side                    binance.SideType        `json:"side"`
	SelfTradePreventionMode binance.STPModeType     `json:"selfTradePreventionMode"`
}
