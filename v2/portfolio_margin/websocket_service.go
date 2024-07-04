package portfolio_margin

import (
	"encoding/json"
	"fmt"
	"time"
)

// Endpoints
const (
	baseWsMainUrl = "wss://fstream.binance.com/ws"
)

var (
	// WebsocketTimeout is an interval for sending ping/pong messages if WebsocketKeepalive is enabled
	WebsocketTimeout = time.Second * 60
	// WebsocketKeepalive enables sending ping/pong messages to check the connection stability
	WebsocketKeepalive = false
	ProxyUrl           = ""
)

func getWsProxyUrl() *string {
	if ProxyUrl == "" {
		return nil
	}
	return &ProxyUrl
}

func SetWsProxyUrl(url string) {
	ProxyUrl = url
}

// getWsEndpoint return the base endpoint of the WS according the UseTestnet flag
func getWsEndpoint() string {
	return baseWsMainUrl
}

// WsUserDataEvent define user data event
type WsUserDataEvent struct {
	Event                   UserDataEventType         `json:"e"`
	Time                    int64                     `json:"E"`
	TransactionTime         int64                     `json:"T"`
	BusinessUnit            string                    `json:"fs"`
	FuturesAccountUpdate    WsFuturesAccountUpdate    `json:"a"`
	FuturesOrderTradeUpdate WsFuturesOrderTradeUpdate `json:"o"`
}

// WsAccountUpdate define account update
type WsFuturesAccountUpdate struct {
	Reason    string              `json:"m"`
	Balances  []WsFuturesBalance  `json:"B"`
	Positions []WsFuturesPosition `json:"P"`
}

// WsBalance define balance
type WsFuturesBalance struct {
	Asset              string `json:"a"`
	Balance            string `json:"wb"`
	CrossWalletBalance string `json:"cw"`
	ChangeBalance      string `json:"bc"`
}

// WsPosition define position
type WsFuturesPosition struct {
	Symbol              string           `json:"s"`
	Amount              string           `json:"pa"`
	EntryPrice          string           `json:"ep"`
	AccumulatedRealized string           `json:"cr"`
	UnrealizedPnL       string           `json:"up"`
	Side                PositionSideType `json:"ps"`
	BreakevenPrice      string           `json:"bep"`
}

// WsOrderTradeUpdate define order trade update
type WsFuturesOrderTradeUpdate struct {
	Symbol               string             `json:"s"`   // Symbol
	ClientOrderID        string             `json:"c"`   // Client order ID
	Side                 SideType           `json:"S"`   // Side
	Type                 OrderType          `json:"o"`   // Order type
	TimeInForce          TimeInForceType    `json:"f"`   // Time in force
	OriginalQty          string             `json:"q"`   // Original quantity
	OriginalPrice        string             `json:"p"`   // Original price
	AveragePrice         string             `json:"ap"`  // Average price
	StopPrice            string             `json:"sp"`  // Stop price. Please ignore with TRAILING_STOP_MARKET order
	ExecutionType        OrderExecutionType `json:"x"`   // Execution type
	Status               OrderStatusType    `json:"X"`   // Order status
	ID                   int64              `json:"i"`   // Order ID
	LastFilledQty        string             `json:"l"`   // Order Last Filled Quantity
	AccumulatedFilledQty string             `json:"z"`   // Order Filled Accumulated Quantity
	LastFilledPrice      string             `json:"L"`   // Last Filled Price
	CommissionAsset      string             `json:"N"`   // Commission Asset, will not push if no commission
	Commission           string             `json:"n"`   // Commission, will not push if no commission
	TradeTime            int64              `json:"T"`   // Order Trade Time
	TradeID              int64              `json:"t"`   // Trade ID
	BidsNotional         string             `json:"b"`   // Bids Notional
	AsksNotional         string             `json:"a"`   // Asks Notional
	IsMaker              bool               `json:"m"`   // Is this trade the maker side?
	IsReduceOnly         bool               `json:"R"`   // Is this reduce only
	WorkingType          WorkingType        `json:"wt"`  // Stop Price Working Type
	OriginalType         OrderType          `json:"ot"`  // Original Order Type
	PositionSide         PositionSideType   `json:"ps"`  // Position Side
	IsClosingPosition    bool               `json:"cp"`  // If Close-All, pushed with conditional order
	ActivationPrice      string             `json:"AP"`  // Activation Price, only puhed with TRAILING_STOP_MARKET order
	CallbackRate         string             `json:"cr"`  // Callback Rate, only puhed with TRAILING_STOP_MARKET order
	PriceProtect         bool               `json:"pP"`  // If price protection is turned on
	RealizedPnL          string             `json:"rp"`  // Realized Profit of the trade
	STP                  string             `json:"V"`   // STP mode
	PriceMode            string             `json:"pm"`  // Price match mode
	GTD                  int64              `json:"gtd"` // TIF GTD order auto cancel time
}

// WsUserDataHandler handle WsUserDataEvent
type WsUserDataHandler func(event *WsUserDataEvent)

// WsUserDataServe serve user data handler with listen key
func WsUserDataServe(listenKey string, handler WsUserDataHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s", getWsEndpoint(), listenKey)
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsUserDataEvent)
		err := json.Unmarshal(message, event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}
