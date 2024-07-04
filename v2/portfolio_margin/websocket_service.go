package portfolio_margin

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/adshao/go-binance/v2/futures"
)

// Endpoints
const (
	baseWsMainUrl = "wss://fstream.binance.com/pm/ws"
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
	BusinessUnit            BusinessUnitType          `json:"fs"`
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
	Symbol              string                   `json:"s"`
	Amount              string                   `json:"pa"`
	EntryPrice          string                   `json:"ep"`
	AccumulatedRealized string                   `json:"cr"`
	UnrealizedPnL       string                   `json:"up"`
	Side                futures.PositionSideType `json:"ps"`
	BreakevenPrice      float64                  `json:"bep"`
}

// WsOrderTradeUpdate define order trade update
type WsFuturesOrderTradeUpdate struct {
	Symbol               string                     `json:"s"`            // Symbol
	ClientOrderID        string                     `json:"c"`            // Client order ID
	Side                 futures.SideType           `json:"S"`            // Side
	Type                 futures.OrderType          `json:"o"`            // Order type
	TimeInForce          futures.TimeInForceType    `json:"f"`            // Time in force
	OriginalQty          string                     `json:"q"`            // Original quantity
	OriginalPrice        string                     `json:"p"`            // Original price
	AveragePrice         string                     `json:"ap"`           // Average price
	StopPrice            string                     `json:"sp"`           // Stop price. Please ignore with TRAILING_STOP_MARKET order
	ExecutionType        futures.OrderExecutionType `json:"x"`            // Execution type
	Status               futures.OrderStatusType    `json:"X"`            // Order status
	ID                   int64                      `json:"i"`            // Order ID
	LastFilledQty        string                     `json:"l"`            // Order Last Filled Quantity
	AccumulatedFilledQty string                     `json:"z"`            // Order Filled Accumulated Quantity
	LastFilledPrice      string                     `json:"L"`            // Last Filled Price
	CommissionAsset      string                     `json:"N"`            // Commission Asset, will not push if no commission
	Commission           string                     `json:"n"`            // Commission, will not push if no commission
	TradeTime            int64                      `json:"T"`            // Order Trade Time
	TradeID              int64                      `json:"t"`            // Trade ID
	BidsNotional         string                     `json:"b"`            // Bids Notional
	AsksNotional         string                     `json:"a"`            // Asks Notional
	IsMaker              bool                       `json:"m"`            // Is this trade the maker side?
	IsReduceOnly         bool                       `json:"R"`            // Is this reduce only
	PositionSide         futures.PositionSideType   `json:"ps"`           // Position Side
	RealizedPnL          string                     `json:"rp"`           // Realized Profit of the trade
	StrategyType         string                     `json:"st,omitempty"` // Strategy type, only pushed with conditional order triggered
	StrategyId           int64                      `json:"si,omitempty"` // StrategyId, only pushed with conditional order triggered
	STP                  string                     `json:"V"`            // STP mode
	GTD                  int64                      `json:"gtd"`          // TIF GTD order auto cancel time
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
			errHandler(fmt.Errorf("%v: %s", err, message))
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}
