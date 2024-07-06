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

// Business Unit
type BusinessUnitType string

// UserDataEventType define user data event type
type UserDataEventType string

const (
	UE_StreamExpired       UserDataEventType = "listenKeyExpired"
	UE_RiskLevelChange     UserDataEventType = "riskLevelChange"
	UE_OpenOrderLossUpdate UserDataEventType = "openOrderLoss"
	UE_LiabilityUpdate     UserDataEventType = "liabilityChange"

	UE_MarginAccountUpdate UserDataEventType = "outboundAccountPosition"
	UE_MarginBalanceUpdate UserDataEventType = "balanceUpdate"
	UE_MarginOrderUpdate   UserDataEventType = "executionReport"

	UE_FuturesAccountUpdate       UserDataEventType = "ACCOUNT_UPDATE"
	UE_FuturesAccountConfigUpdate UserDataEventType = "ACCOUNT_CONFIG_UPDATE"
	UE_FuturesOrderUpdate         UserDataEventType = "ORDER_TRADE_UPDATE"
	UE_FuturesCondOrderUpdate     UserDataEventType = "CONDITIONAL_ORDER_TRADE_UPDATE"

	BusinessUnitTypeUM BusinessUnitType = "UM"
	BusinessUnitTypeCM BusinessUnitType = "CM"
)

type WsUserDataFuturesAccountUpdateEvent struct {
	Event                UserDataEventType `json:"e"`
	Time                 int64             `json:"E"`
	TransactionTime      int64             `json:"T"`
	BusinessUnit         BusinessUnitType  `json:"fs"`
	AccountAlias         string            `json:"i"`
	FuturesAccountUpdate struct {
		Reason   string `json:"m"`
		Balances []struct {
			Asset              string `json:"a"`
			Balance            string `json:"wb"`
			CrossWalletBalance string `json:"cw"`
			ChangeBalance      string `json:"bc"`
		} `json:"B"`
		Positions []struct {
			Symbol              string                   `json:"s"`
			Amount              string                   `json:"pa"`
			EntryPrice          string                   `json:"ep"`
			AccumulatedRealized string                   `json:"cr"`
			UnrealizedPnL       string                   `json:"up"`
			Side                futures.PositionSideType `json:"ps"`
			BreakevenPrice      float64                  `json:"bep"`
		} `json:"P"`
	} `json:"a"`
}

type WsUserDataFuturesOrderUpdateEvent struct {
	Event              UserDataEventType `json:"e"`
	Time               int64             `json:"E"`
	TransactionTime    int64             `json:"T"`
	BusinessUnit       BusinessUnitType  `json:"fs"`
	AccountAlias       string            `json:"i"`
	FuturesOrderUpdate struct {
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
	} `json:"o"`
}

type WsUserDataEvent struct {
	Event                     UserDataEventType `json:"e"`
	Time                      int64             `json:"E"`
	FuturesAccountUpdateEvent *WsUserDataFuturesAccountUpdateEvent
	FuturesOrderUpdateEvent   *WsUserDataFuturesOrderUpdateEvent
}

// WsUserDataServe serve user data handler with listen key
func WsUserDataServe(
	listenKey string,
	handler func(event *WsUserDataEvent),
	errHandler ErrHandler,
) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s", getWsEndpoint(), listenKey)
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsUserDataEvent)
		err := json.Unmarshal(message, event)
		if err != nil {
			errHandler(fmt.Errorf("error unmarshalling WsUserDataEvent: %v: %s", err, message))
			return
		}

		switch event.Event {

		case UE_StreamExpired:

		case UE_FuturesAccountUpdate:
			subEvent := new(WsUserDataFuturesAccountUpdateEvent)
			if err := json.Unmarshal(message, subEvent); err != nil {
				errHandler(fmt.Errorf("error unmarshalling WsUserDataFuturesAccountUpdateEvent: %v: %s", err, message))
				return
			} else {
				event.FuturesAccountUpdateEvent = subEvent
			}

		case UE_FuturesOrderUpdate:
			subEvent := new(WsUserDataFuturesOrderUpdateEvent)
			if err := json.Unmarshal(message, subEvent); err != nil {
				errHandler(fmt.Errorf("error unmarshalling WsUserDataFuturesOrderUpdateEvent: %v: %s", err, message))
				return
			} else {
				event.FuturesOrderUpdateEvent = subEvent
			}
		}

		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}
