package portfolio_margin

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/adshao/go-binance/v2"
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
	UETypeStreamExpired       UserDataEventType = "listenKeyExpired"
	UETypeRiskLevelChange     UserDataEventType = "riskLevelChange"
	UETypeOpenOrderLossUpdate UserDataEventType = "openOrderLoss"

	UETypeFuturesAccountUpdate   UserDataEventType = "ACCOUNT_UPDATE"
	UETypeFuturesLeverageUpdate  UserDataEventType = "ACCOUNT_CONFIG_UPDATE"
	UETypeFuturesOrderUpdate     UserDataEventType = "ORDER_TRADE_UPDATE"
	UETypeFuturesCondOrderUpdate UserDataEventType = "CONDITIONAL_ORDER_TRADE_UPDATE"

	UETypeMarginAccountUpdate   UserDataEventType = "outboundAccountPosition"
	UETypeMarginBalanceUpdate   UserDataEventType = "balanceUpdate"
	UETypeMarginOrderUpdate     UserDataEventType = "executionReport"
	UETypeMarginLiabilityUpdate UserDataEventType = "liabilityChange"

	BusinessUnitTypeUM BusinessUnitType = "UM"
	BusinessUnitTypeCM BusinessUnitType = "CM"
)

type WsUserDataFuturesAccountUpdateEvent struct {
	Event           UserDataEventType `json:"e"`
	Time            int64             `json:"E"`
	TransactionTime int64             `json:"T"`
	BusinessUnit    BusinessUnitType  `json:"fs"`
	AccountAlias    string            `json:"i"`
	AccountUpdate   struct {
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

type WsUserDataFuturesLeverageUpdateEvent struct {
	Event           UserDataEventType `json:"e"`
	Time            int64             `json:"E"`
	TransactionTime int64             `json:"T"`
	BusinessUnit    BusinessUnitType  `json:"fs"`
	LeverageUpdate  struct {
		Symbol   string `json:"s"`
		Leverage int64  `json:"l"`
	} `json:"ac"`
}

type WsUserDataFuturesOrderUpdateEvent struct {
	Event           UserDataEventType `json:"e"`
	Time            int64             `json:"E"`
	TransactionTime int64             `json:"T"`
	BusinessUnit    BusinessUnitType  `json:"fs"`
	AccountAlias    string            `json:"i"`
	OrderUpdate     struct {
		Symbol               string                     `json:"s"`            // Symbol
		ClientOrderId        string                     `json:"c"`            // Client order ID
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
		LastFilledPrice      string                     `json:"L"`            // Last Filled Price
		AccumulatedFilledQty string                     `json:"z"`            // Order Filled Accumulated Quantity
		CommissionAsset      string                     `json:"N"`            // Commission Asset, will not push if no commission
		Commission           string                     `json:"n"`            // Commission, will not push if no commission
		TradeTime            int64                      `json:"T"`            // Order Trade Time
		TradeId              int64                      `json:"t"`            // Trade ID
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

type WsUserDataMarginLiabilityUpdateEvent struct {
	Event          UserDataEventType `json:"e"`
	Time           int64             `json:"E"`
	Asset          string            `json:"a"`
	Type           string            `json:"t"`
	TransactionId  int64             `json:"tx"`
	Principal      string            `json:"p"`
	Interest       string            `json:"i"`
	TotalLiability string            `json:"l"`
}

type WsUserDataMarginAccountUpdateEvent struct {
	Event          UserDataEventType `json:"e"`
	Time           int64             `json:"E"`
	LastUpdateTime int64             `json:"u"`
	UpdateId       int64             `json:"U"`
	Balances       []struct {
		Asset  string `json:"a"`
		Free   string `json:"f"`
		Locked string `json:"l"`
	} `json:"B"`
}

type WsUserDataMarginBalanceUpdateEvent struct {
	Event           UserDataEventType `json:"e"`
	Time            int64             `json:"E"`
	TransactionTime int64             `json:"T"`
	UpdateId        int64             `json:"U"`
	Asset           string            `json:"a"`
	BalanceDelta    string            `json:"d"`
}

type WsUserDataMarginOrderUpdateEvent struct {
	Event                   UserDataEventType       `json:"e"` // Event type
	Time                    int64                   `json:"E"` // Event time
	Symbol                  string                  `json:"s"` // Symbol
	ClientOrderId           string                  `json:"c"` // Client order ID
	Side                    string                  `json:"S"` // Side
	Type                    string                  `json:"o"` // Order type
	TimeInForce             binance.TimeInForceType `json:"f"` // Time in force
	Quantity                string                  `json:"q"` // Order quantity
	Price                   string                  `json:"p"` // Order price
	StopPrice               string                  `json:"P"` // Stop price
	IcebergQuantity         string                  `json:"F"` // Iceberg quantity
	OrderListId             int64                   `json:"g"` // OrderListId
	OrigCustomOrderId       string                  `json:"C"` // Original client order ID; This is the ID of the order being canceled
	ExecutionType           string                  `json:"x"` /// Current execution type
	Status                  string                  `json:"X"` // Current order status
	RejectReason            string                  `json:"r"` // Order reject reason; will be an error code.
	Id                      int64                   `json:"i"` // Order ID
	LastFilledQuantity      string                  `json:"l"` // Last executed quantity
	LastFilledPrice         string                  `json:"L"` // Last executed price
	FilledQuantity          string                  `json:"z"` // Cumulative filled quantity
	FeeAsset                string                  `json:"N"` // Commission asset
	FeeCost                 string                  `json:"n"` // Commission amount
	TransactionTime         int64                   `json:"T"` // Transaction time
	TradeId                 int64                   `json:"t"` // Trade ID
	IgnoreI                 int64                   `json:"I"` // Ignore
	IsInOrderBook           bool                    `json:"w"` // Is the order on the book?
	IsMaker                 bool                    `json:"m"` // Is this trade the maker side?
	IgnoreM                 bool                    `json:"M"` // Ignore
	CreateTime              int64                   `json:"O"` // Order creation time
	FilledQuoteQuantity     string                  `json:"Z"` // Cumulative quote asset transacted quantity
	LastQuoteQuantity       float64                 `json:"Y"` // Last quote asset transacted quantity (i.e. lastPrice * lastQty)
	QuoteQuantity           string                  `json:"Q"` // Quote Order Quantity
	SelfTradePreventionMode string                  `json:"V"` // selfTradePreventionMode

	//These are fields that appear in the payload only if certain conditions are met.
	TrailingDelta         int64  `json:"d"` // Trailing Delta; This is only visible if the order was a trailing stop order.
	PreventedMatchId      int64  `json:"v"` // Prevented Match Id; This is only visible if the order expire due to STP trigger.
	TrailingTime          int64  `json:"D"` // Trailing Time; This is only visible if the trailing stop order has been activated.
	StrategyId            int64  `json:"j"` // Strategy ID; This is only visible if the strategyId parameter was provided upon order placement
	StrategyType          int64  `json:"J"` // Strategy Type; This is only visible if the strategyType parameter was provided upon order placement
	WorkingTime           int64  `json:"W"` // Working Time; This is only visible if the order has been placed on the book.
	TradeGroupId          int64  `json:"u"` // TradeGroupId; This is only visible if the account is part of a trade group and the order expired due to STP trigger.
	CounterOrderId        int64  `json:"U"` // CounterOrderId; This is only visible if the order expired due to STP trigger.
	PreventedQuantity     string `json:"A"` // Prevented Quantity; This is only visible if the order expired due to STP trigger.
	LastPreventedQuantity string `json:"B"` // Last Prevented Quantity; This is only visible if the order expired due to STP trigger.
}

type WsUserDataEvent struct {
	Event                      UserDataEventType `json:"e"`
	Time                       int64             `json:"E"`
	FuturesAccountUpdateEvent  *WsUserDataFuturesAccountUpdateEvent
	FuturesLeverageUpdateEvent *WsUserDataFuturesLeverageUpdateEvent
	FuturesOrderUpdateEvent    *WsUserDataFuturesOrderUpdateEvent
	MarginLiabilityUpdateEvent *WsUserDataMarginLiabilityUpdateEvent
	MarginAccountUpdateEvent   *WsUserDataMarginAccountUpdateEvent
	MarginBalanceUpdateEvent   *WsUserDataMarginBalanceUpdateEvent
	MarginOrderUpdateEvent     *WsUserDataMarginOrderUpdateEvent
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
			errHandler(fmt.Errorf("WsUserDataEvent: %v: %s", err, message))
			return
		}

		switch event.Event {

		case UETypeStreamExpired:

		case UETypeFuturesAccountUpdate:
			subEvent := new(WsUserDataFuturesAccountUpdateEvent)
			if err := json.Unmarshal(message, subEvent); err != nil {
				errHandler(fmt.Errorf("WsUserDataFuturesAccountUpdateEvent: %v: %s", err, message))
				return
			} else {
				event.FuturesAccountUpdateEvent = subEvent
			}

		case UETypeFuturesLeverageUpdate:
			subEvent := new(WsUserDataFuturesLeverageUpdateEvent)
			if err := json.Unmarshal(message, subEvent); err != nil {
				errHandler(fmt.Errorf("WsUserDataFuturesLeverageUpdateEvent: %v: %s", err, message))
				return
			} else {
				event.FuturesLeverageUpdateEvent = subEvent
			}

		case UETypeFuturesOrderUpdate:
			subEvent := new(WsUserDataFuturesOrderUpdateEvent)
			if err := json.Unmarshal(message, subEvent); err != nil {
				errHandler(fmt.Errorf("WsUserDataFuturesOrderUpdateEvent: %v: %s", err, message))
				return
			} else {
				event.FuturesOrderUpdateEvent = subEvent
			}

		case UETypeMarginLiabilityUpdate:
			subEvent := new(WsUserDataMarginLiabilityUpdateEvent)
			if err := json.Unmarshal(message, subEvent); err != nil {
				errHandler(fmt.Errorf("WsUserDataMarginLiabilityUpdateEvent: %v: %s", err, message))
				return
			} else {
				event.MarginLiabilityUpdateEvent = subEvent
			}

		case UETypeMarginAccountUpdate:
			subEvent := new(WsUserDataMarginAccountUpdateEvent)
			if err := json.Unmarshal(message, subEvent); err != nil {
				errHandler(fmt.Errorf("WsUserDataMarginAccountUpdateEvent: %v: %s", err, message))
				return
			} else {
				event.MarginAccountUpdateEvent = subEvent
			}

		case UETypeMarginBalanceUpdate:
			subEvent := new(WsUserDataMarginBalanceUpdateEvent)
			if err := json.Unmarshal(message, subEvent); err != nil {
				errHandler(fmt.Errorf("WsUserDataMarginBalanceUpdateEvent: %v: %s", err, message))
				return
			} else {
				event.MarginBalanceUpdateEvent = subEvent
			}

		case UETypeMarginOrderUpdate:
			subEvent := new(WsUserDataMarginOrderUpdateEvent)
			if err := json.Unmarshal(message, subEvent); err != nil {
				errHandler(fmt.Errorf("WsUserDataMarginOrderUpdateEvent: %v: %s", err, message))
				return
			} else {
				event.MarginOrderUpdateEvent = subEvent
			}
		}

		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}
