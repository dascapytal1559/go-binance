package pmargin

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/adshao/go-binance/v2"
)

// NewCreateMarginOrderService init creating UM order service
func (c *Client) NewCreateMarginOrderService() *CreateMarginOrderService {
	return &CreateMarginOrderService{c: c}
}

// CreateMarginOrderService create order
type CreateMarginOrderService struct {
	c                       *Client
	symbol                  string
	side                    binance.SideType
	orderType               binance.OrderType
	quantity                *string
	quoteOrderQty           *string
	price                   *string
	stopPrice               *string
	newClientOrderID        *string
	newOrderRespType        *binance.NewOrderRespType
	icebergQty              *string
	sideEffectType          *binance.SideEffectType
	timeInForce             *binance.TimeInForceType
	selfTradePreventionMode *binance.STPModeType
	autoRepayAtCancel       *bool
}

// Symbol set symbol
func (s *CreateMarginOrderService) Symbol(symbol string) *CreateMarginOrderService {
	s.symbol = symbol
	return s
}

// Side set side
func (s *CreateMarginOrderService) Side(side binance.SideType) *CreateMarginOrderService {
	s.side = side
	return s
}

// Type set type
func (s *CreateMarginOrderService) Type(orderType binance.OrderType) *CreateMarginOrderService {
	s.orderType = orderType
	return s
}

// Quantity set quantity
func (s *CreateMarginOrderService) Quantity(quantity string) *CreateMarginOrderService {
	s.quantity = &quantity
	return s
}

// QuoteOrderQty set quoteOrderQty
func (s *CreateMarginOrderService) QuoteOrderQty(quoteOrderQty string) *CreateMarginOrderService {
	s.quoteOrderQty = &quoteOrderQty
	return s
}

// Price set price
func (s *CreateMarginOrderService) Price(price string) *CreateMarginOrderService {
	s.price = &price
	return s
}

// StopPrice set stopPrice
func (s *CreateMarginOrderService) StopPrice(stopPrice string) *CreateMarginOrderService {
	s.stopPrice = &stopPrice
	return s
}

// NewClientOrderID set newClientOrderID
func (s *CreateMarginOrderService) NewClientOrderID(newClientOrderID string) *CreateMarginOrderService {
	s.newClientOrderID = &newClientOrderID
	return s
}

// NewOrderResponseType set newOrderResponseType
func (s *CreateMarginOrderService) NewOrderResponseType(newOrderResponseType binance.NewOrderRespType) *CreateMarginOrderService {
	s.newOrderRespType = &newOrderResponseType
	return s
}

// IcebergQty set icebergQty
func (s *CreateMarginOrderService) IcebergQty(icebergQty string) *CreateMarginOrderService {
	s.icebergQty = &icebergQty
	return s
}

// SideEffectType set sideEffectType
func (s *CreateMarginOrderService) SideEffectType(sideEffectType binance.SideEffectType) *CreateMarginOrderService {
	s.sideEffectType = &sideEffectType
	return s
}

// TimeInForce set timeInForce
func (s *CreateMarginOrderService) TimeInForce(timeInForce binance.TimeInForceType) *CreateMarginOrderService {
	s.timeInForce = &timeInForce
	return s
}

// SelfTradePreventionMode set selfTradePreventionMode
func (s *CreateMarginOrderService) SelfTradePreventionMode(selfTradePreventionMode binance.STPModeType) *CreateMarginOrderService {
	s.selfTradePreventionMode = &selfTradePreventionMode
	return s
}

// AutoRepayAtCancel set autoRepayAtCancel
func (s *CreateMarginOrderService) AutoRepayAtCancel(autoRepayAtCancel bool) *CreateMarginOrderService {
	s.autoRepayAtCancel = &autoRepayAtCancel
	return s
}

// Do send request
func (s *CreateMarginOrderService) Do(ctx context.Context, opts ...RequestOption) (res *CreateMarginOrderResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/papi/v1/margin/order",
		secType:  secTypeSigned,
	}
	r.setFormParams(params{
		"symbol":                  s.symbol,
		"side":                    s.side,
		"orderType":               s.orderType,
		"quantity":                s.quantity,
		"quoteOrderQty":           s.quoteOrderQty,
		"price":                   s.price,
		"stopPrice":               s.stopPrice,
		"newClientOrderID":        s.newClientOrderID,
		"newOrderRespType":        s.newOrderRespType,
		"icebergQty":              s.icebergQty,
		"sideEffectType":          s.sideEffectType,
		"timeInForce":             s.timeInForce,
		"selfTradePreventionMode": s.selfTradePreventionMode,
		"autoRepayAtCancel":       s.autoRepayAtCancel,
	})
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(CreateMarginOrderResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// OrderResponse define order info
type CreateMarginOrderResponse struct {
	Symbol                  string                  `json:"symbol"`
	OrderID                 int64                   `json:"orderId"`
	ClientOrderID           string                  `json:"clientOrderId"`
	TransactTime            int64                   `json:"transactTime"`
	Price                   string                  `json:"price"`
	SelfTradePreventionMode binance.STPModeType     `json:"selfTradePreventionMode"`
	OrigQty                 string                  `json:"origQty"`
	ExecutedQty             string                  `json:"executedQty"`
	CummulativeQuoteQty     string                  `json:"cummulativeQuoteQty"`
	Status                  binance.OrderStatusType `json:"status"`
	TimeInForce             binance.TimeInForceType `json:"timeInForce"`
	Type                    binance.OrderType       `json:"type"`
	Side                    binance.SideType        `json:"side"`
	MarginBuyBorrowAmount   string                  `json:"marginBuyBorrowAmount,omitempty"`
	MarginBuyBorrowAsset    string                  `json:"marginBuyBorrowAsset,omitempty"`
	Fills                   []struct {
		Price           string `json:"price"`
		Qty             string `json:"qty"`
		Commission      string `json:"commission"`
		CommissionAsset string `json:"commissionAsset"`
	} `json:"fills"`
}
