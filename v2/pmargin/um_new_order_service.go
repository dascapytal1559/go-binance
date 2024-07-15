package pmargin

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/adshao/go-binance/v2/futures"
)

// CreateUMOrderService create order
type CreateUMOrderService struct {
	c                *Client
	symbol           string
	side             futures.SideType
	orderType        futures.OrderType
	positionSide     *futures.PositionSideType
	timeInForce      *futures.TimeInForceType
	quantity         *string
	reduceOnly       *bool
	price            *string
	newClientOrderID *string
	newOrderRespType *futures.NewOrderRespType
}

// Symbol set symbol
func (s *CreateUMOrderService) Symbol(symbol string) *CreateUMOrderService {
	s.symbol = symbol
	return s
}

// Side set side
func (s *CreateUMOrderService) Side(side futures.SideType) *CreateUMOrderService {
	s.side = side
	return s
}

// Type set type
func (s *CreateUMOrderService) Type(orderType futures.OrderType) *CreateUMOrderService {
	s.orderType = orderType
	return s
}

// PositionSide set side
func (s *CreateUMOrderService) PositionSide(positionSide futures.PositionSideType) *CreateUMOrderService {
	s.positionSide = &positionSide
	return s
}

// TimeInForce set timeInForce
func (s *CreateUMOrderService) TimeInForce(timeInForce futures.TimeInForceType) *CreateUMOrderService {
	s.timeInForce = &timeInForce
	return s
}

// Quantity set quantity
func (s *CreateUMOrderService) Quantity(quantity string) *CreateUMOrderService {
	s.quantity = &quantity
	return s
}

// ReduceOnly set reduceOnly
func (s *CreateUMOrderService) ReduceOnly(reduceOnly bool) *CreateUMOrderService {
	s.reduceOnly = &reduceOnly
	return s
}

// Price set price
func (s *CreateUMOrderService) Price(price string) *CreateUMOrderService {
	s.price = &price
	return s
}

// NewClientOrderID set newClientOrderID
func (s *CreateUMOrderService) NewClientOrderID(newClientOrderID string) *CreateUMOrderService {
	s.newClientOrderID = &newClientOrderID
	return s
}

// NewOrderResponseType set newOrderResponseType
func (s *CreateUMOrderService) NewOrderResponseType(newOrderResponseType futures.NewOrderRespType) *CreateUMOrderService {
	s.newOrderRespType = &newOrderResponseType
	return s
}

func (s *CreateUMOrderService) createOrder(ctx context.Context, endpoint string, opts ...RequestOption) (data []byte, header *http.Header, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: endpoint,
		secType:  secTypeSigned,
	}

	r.setFormParams(params{
		"symbol":           s.symbol,
		"side":             s.side,
		"type":             s.orderType,
		"quantity":         s.quantity,
		"positionSide":     s.positionSide,
		"timeInForce":      s.timeInForce,
		"reduceOnly":       s.reduceOnly,
		"price":            s.price,
		"newClientOrderId": s.newClientOrderID,
		"newOrderRespType": s.newOrderRespType,
	})

	data, header, err = s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []byte{}, &http.Header{}, err
	}
	return data, header, nil
}

// Do send request
func (s *CreateUMOrderService) Do(ctx context.Context, opts ...RequestOption) (res *CreateUMOrderResponse, err error) {
	data, header, err := s.createOrder(ctx, "/papi/v1/um/order", opts...)
	if err != nil {
		return nil, err
	}

	res = new(CreateUMOrderResponse)

	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}

	res.RateLimitOrder10s = header.Get("X-Mbx-Order-Count-10s")
	res.RateLimitOrder1m = header.Get("X-Mbx-Order-Count-1m")
	return res, nil
}

// OrderResponse define order info
type CreateUMOrderResponse struct {
	ClientOrderId           string                   `json:"clientOrderId"`
	CumQty                  string                   `json:"cumQty"`
	CumQuote                string                   `json:"cumQuote"`
	ExecutedQty             string                   `json:"executedQty"`
	OrderId                 int64                    `json:"orderId"`
	AvgPrice                string                   `json:"avgPrice"`
	OrigQty                 string                   `json:"origQty"`
	Price                   string                   `json:"price"`
	ReduceOnly              bool                     `json:"reduceOnly"`
	Side                    string                   `json:"side"`
	PositionSide            futures.PositionSideType `json:"positionSide"`
	Status                  futures.OrderStatusType  `json:"status"`
	Symbol                  string                   `json:"symbol"`
	TimeInForce             futures.TimeInForceType  `json:"timeInForce"`
	Type                    futures.OrderType        `json:"type"`
	SelfTradePreventionMode string                   `json:"selfTradePreventionMode"`
	GoodTillDate            int64                    `json:"goodTillDate"`
	UpdateTime              int64                    `json:"updateTime"`
	RateLimitOrder10s       string                   `json:"rateLimitOrder10s,omitempty"`
	RateLimitOrder1m        string                   `json:"rateLimitOrder1m,omitempty"`
}
