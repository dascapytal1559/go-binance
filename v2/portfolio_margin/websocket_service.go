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
	Event           UserDataEventType `json:"e"`
	Time            int64             `json:"E"`
	TransactionTime int64             `json:"T"`
	BusinessUnit    string            `json:"fs"`
	AccountUpdate   WsAccountUpdate   `json:"a"`
	// CrossWalletBalance  string                `json:"cw"`
	// MarginCallPositions []WsPosition          `json:"p"`
	// OrderTradeUpdate    WsOrderTradeUpdate    `json:"o"`
	// AccountConfigUpdate WsAccountConfigUpdate `json:"ac"`
}

// WsAccountUpdate define account update
type WsAccountUpdate struct {
	Reason    string       `json:"m"`
	Balances  []WsBalance  `json:"B"`
	Positions []WsPosition `json:"P"`
}

// WsBalance define balance
type WsBalance struct {
	Asset              string `json:"a"`
	Balance            string `json:"wb"`
	CrossWalletBalance string `json:"cw"`
	ChangeBalance      string `json:"bc"`
}

// WsPosition define position
type WsPosition struct {
	Symbol              string           `json:"s"`
	Amount              string           `json:"pa"`
	EntryPrice          string           `json:"ep"`
	AccumulatedRealized string           `json:"cr"`
	UnrealizedPnL       string           `json:"up"`
	Side                PositionSideType `json:"ps"`
	BreakevenPrice      string           `json:"bep"`
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
