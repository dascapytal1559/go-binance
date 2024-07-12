package futures

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/adshao/go-binance/v2/common"
)

// Endpoints
const (
	baseWsMainUrl          = "wss://fstream.binance.com/ws"
	baseWsTestnetUrl       = "wss://stream.binancefuture.com/ws"
	baseCombinedMainURL    = "wss://fstream.binance.com/stream?streams="
	baseCombinedTestnetURL = "wss://stream.binancefuture.com/stream?streams="
)

var (
	// WebsocketTimeout is an interval for sending ping/pong messages if WebsocketKeepalive is enabled
	WebsocketTimeout = time.Second * 60
	// WebsocketKeepalive enables sending ping/pong messages to check the connection stability
	WebsocketKeepalive = false
	// UseTestnet switch all the WS streams from production to the testnet
	UseTestnet = false
	ProxyUrl   = ""
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
	if UseTestnet {
		return baseWsTestnetUrl
	}
	return baseWsMainUrl
}

// getCombinedEndpoint return the base endpoint of the combined stream according the UseTestnet flag
func getCombinedEndpoint() string {
	if UseTestnet {
		return baseCombinedTestnetURL
	}
	return baseCombinedMainURL
}

// WsAggTradeEvent define websocket aggTrde event.
type WsAggTradeEvent struct {
	Event            string `json:"e"`
	Time             int64  `json:"E"`
	Symbol           string `json:"s"`
	AggregateTradeID int64  `json:"a"`
	Price            string `json:"p"`
	Quantity         string `json:"q"`
	FirstTradeID     int64  `json:"f"`
	LastTradeID      int64  `json:"l"`
	TradeTime        int64  `json:"T"`
	Maker            bool   `json:"m"`
}

// WsAggTradeHandler handle websocket that push trade information that is aggregated for a single taker order.
type WsAggTradeHandler func(event *WsAggTradeEvent)

// WsAggTradeServe serve websocket that push trade information that is aggregated for a single taker order.
func WsAggTradeServe(symbol string, handler WsAggTradeHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s@aggTrade", getWsEndpoint(), strings.ToLower(symbol))
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsAggTradeEvent)
		err := json.Unmarshal(message, &event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsCombinedAggTradeServe is similar to WsAggTradeServe, but it handles multiple symbols
func WsCombinedAggTradeServe(symbols []string, handler WsAggTradeHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := getCombinedEndpoint()
	for _, s := range symbols {
		endpoint += fmt.Sprintf("%s@aggTrade", strings.ToLower(s)) + "/"
	}
	endpoint = endpoint[:len(endpoint)-1]
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		j, err := newJSON(message)
		if err != nil {
			errHandler(err)
			return
		}

		stream := j.Get("stream").MustString()
		data := j.Get("data").MustMap()

		symbol := strings.Split(stream, "@")[0]

		jsonData, _ := json.Marshal(data)

		event := new(WsAggTradeEvent)
		err = json.Unmarshal(jsonData, event)
		if err != nil {
			errHandler(err)
			return
		}
		event.Symbol = strings.ToUpper(symbol)

		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsMarkPriceEvent define websocket markPriceUpdate event.
type WsMarkPriceEvent struct {
	Event                string `json:"e"`
	Time                 int64  `json:"E"`
	Symbol               string `json:"s"`
	MarkPrice            string `json:"p"`
	IndexPrice           string `json:"i"`
	EstimatedSettlePrice string `json:"P"`
	FundingRate          string `json:"r"`
	NextFundingTime      int64  `json:"T"`
}

// WsMarkPriceHandler handle websocket that pushes price and funding rate for a single symbol.
type WsMarkPriceHandler func(event *WsMarkPriceEvent)

func wsMarkPriceServe(endpoint string, handler WsMarkPriceHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsMarkPriceEvent)
		err := json.Unmarshal(message, &event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsMarkPriceServe serve websocket that pushes price and funding rate for a single symbol.
func WsMarkPriceServe(symbol string, handler WsMarkPriceHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s@markPrice", getWsEndpoint(), strings.ToLower(symbol))
	return wsMarkPriceServe(endpoint, handler, errHandler)
}

// WsMarkPriceServeWithRate serve websocket that pushes price and funding rate for a single symbol and rate.
func WsMarkPriceServeWithRate(symbol string, rate time.Duration, handler WsMarkPriceHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	var rateStr string
	switch rate {
	case 3 * time.Second:
		rateStr = ""
	case 1 * time.Second:
		rateStr = "@1s"
	default:
		return nil, nil, errors.New("invalid rate")
	}
	endpoint := fmt.Sprintf("%s/%s@markPrice%s", getWsEndpoint(), strings.ToLower(symbol), rateStr)
	return wsMarkPriceServe(endpoint, handler, errHandler)
}

func wsCombinedMarkPriceServe(endpoint string, handler WsMarkPriceHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		j, err := newJSON(message)
		if err != nil {
			errHandler(err)
			return
		}

		data := j.Get("data").MustMap()
		jsonData, _ := json.Marshal(data)

		event := new(WsMarkPriceEvent)
		err = json.Unmarshal(jsonData, event)
		if err != nil {
			errHandler(err)
			return
		}

		handler(event)
	}

	return wsServe(cfg, wsHandler, errHandler)
}

// WsCombinedMarkPriceServe is similar to WsMarkPriceServe, but it handles multiple symbols
func WsCombinedMarkPriceServe(symbols []string, handler WsMarkPriceHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := getCombinedEndpoint()
	for _, s := range symbols {
		endpoint += fmt.Sprintf("%s@markPrice", strings.ToLower(s)) + "/"
	}
	endpoint = endpoint[:len(endpoint)-1]

	return wsCombinedMarkPriceServe(endpoint, handler, errHandler)
}

// WsCombinedMarkPriceServeWithRate is similar to WsMarkPriceServeWithRate, but it for multiple symbols
func WsCombinedMarkPriceServeWithRate(symbolLevels map[string]time.Duration, handler WsMarkPriceHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := getCombinedEndpoint()
	for symbol, rate := range symbolLevels {
		var rateStr string
		switch rate {
		case 3 * time.Second:
			rateStr = ""
		case 1 * time.Second:
			rateStr = "@1s"
		default:
			return nil, nil, fmt.Errorf("invalid rate. Symbol %s (rate %d)", symbol, rate)
		}

		endpoint += fmt.Sprintf("%s@markPrice%s", strings.ToLower(symbol), rateStr) + "/"
	}

	endpoint = endpoint[:len(endpoint)-1]

	return wsCombinedMarkPriceServe(endpoint, handler, errHandler)
}

// WsAllMarkPriceEvent defines an array of websocket markPriceUpdate events.
type WsAllMarkPriceEvent []*WsMarkPriceEvent

// WsAllMarkPriceHandler handle websocket that pushes price and funding rate for all symbol.
type WsAllMarkPriceHandler func(event WsAllMarkPriceEvent)

func wsAllMarkPriceServe(endpoint string, handler WsAllMarkPriceHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		var event WsAllMarkPriceEvent
		err := json.Unmarshal(message, &event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsAllMarkPriceServe serve websocket that pushes price and funding rate for all symbol.
func WsAllMarkPriceServe(handler WsAllMarkPriceHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/!markPrice@arr", getWsEndpoint())
	return wsAllMarkPriceServe(endpoint, handler, errHandler)
}

// WsAllMarkPriceServeWithRate serve websocket that pushes price and funding rate for all symbol and rate.
func WsAllMarkPriceServeWithRate(rate time.Duration, handler WsAllMarkPriceHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	var rateStr string
	switch rate {
	case 3 * time.Second:
		rateStr = ""
	case 1 * time.Second:
		rateStr = "@1s"
	default:
		return nil, nil, errors.New("invalid rate")
	}
	endpoint := fmt.Sprintf("%s/!markPrice@arr%s", getWsEndpoint(), rateStr)
	return wsAllMarkPriceServe(endpoint, handler, errHandler)
}

// WsKlineEvent define websocket kline event
type WsKlineEvent struct {
	Event  string  `json:"e"`
	Time   int64   `json:"E"`
	Symbol string  `json:"s"`
	Kline  WsKline `json:"k"`
}

// WsKline define websocket kline
type WsKline struct {
	StartTime            int64  `json:"t"`
	EndTime              int64  `json:"T"`
	Symbol               string `json:"s"`
	Interval             string `json:"i"`
	FirstTradeID         int64  `json:"f"`
	LastTradeID          int64  `json:"L"`
	Open                 string `json:"o"`
	Close                string `json:"c"`
	High                 string `json:"h"`
	Low                  string `json:"l"`
	Volume               string `json:"v"`
	TradeNum             int64  `json:"n"`
	IsFinal              bool   `json:"x"`
	QuoteVolume          string `json:"q"`
	ActiveBuyVolume      string `json:"V"`
	ActiveBuyQuoteVolume string `json:"Q"`
}

// WsKlineHandler handle websocket kline event
type WsKlineHandler func(event *WsKlineEvent)

// WsKlineServe serve websocket kline handler with a symbol and interval like 15m, 30s
func WsKlineServe(symbol string, interval string, handler WsKlineHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s@kline_%s", getWsEndpoint(), strings.ToLower(symbol), interval)
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsKlineEvent)
		err := json.Unmarshal(message, event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsCombinedKlineServe is similar to WsKlineServe, but it handles multiple symbols with it interval
func WsCombinedKlineServe(symbolIntervalPair map[string]string, handler WsKlineHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := getCombinedEndpoint()
	for symbol, interval := range symbolIntervalPair {
		endpoint += fmt.Sprintf("%s@kline_%s", strings.ToLower(symbol), interval) + "/"
	}
	endpoint = endpoint[:len(endpoint)-1]
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		j, err := newJSON(message)
		if err != nil {
			errHandler(err)
			return
		}

		stream := j.Get("stream").MustString()
		data := j.Get("data").MustMap()

		symbol := strings.Split(stream, "@")[0]

		jsonData, _ := json.Marshal(data)

		event := new(WsKlineEvent)
		err = json.Unmarshal(jsonData, event)
		if err != nil {
			errHandler(err)
			return
		}
		event.Symbol = strings.ToUpper(symbol)

		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsContinuousKlineEvent define websocket continuous kline event
type WsContinuousKlineEvent struct {
	Event        string            `json:"e"`
	Time         int64             `json:"E"`
	PairSymbol   string            `json:"ps"`
	ContractType string            `json:"ct"`
	Kline        WsContinuousKline `json:"k"`
}

// WsContinuousKline define websocket continuous kline
type WsContinuousKline struct {
	StartTime            int64  `json:"t"`
	EndTime              int64  `json:"T"`
	Interval             string `json:"i"`
	FirstTradeID         int64  `json:"f"`
	LastTradeID          int64  `json:"L"`
	Open                 string `json:"o"`
	Close                string `json:"c"`
	High                 string `json:"h"`
	Low                  string `json:"l"`
	Volume               string `json:"v"`
	TradeNum             int64  `json:"n"`
	IsFinal              bool   `json:"x"`
	QuoteVolume          string `json:"q"`
	ActiveBuyVolume      string `json:"V"`
	ActiveBuyQuoteVolume string `json:"Q"`
}

// WsContinuousKlineSubscribeArgs used with WsContinuousKlineServe or WsCombinedContinuousKlineServe
type WsContinuousKlineSubscribeArgs struct {
	Pair         string
	ContractType string
	Interval     string
}

// WsContinuousKlineHandler handle websocket continuous kline event
type WsContinuousKlineHandler func(event *WsContinuousKlineEvent)

// WsContinuousKlineServe serve websocket continuous kline handler with a pair and contractType and interval like 15m, 30s
func WsContinuousKlineServe(subscribeArgs *WsContinuousKlineSubscribeArgs, handler WsContinuousKlineHandler,
	errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s_%s@continuousKline_%s", getWsEndpoint(), strings.ToLower(subscribeArgs.Pair),
		strings.ToLower(subscribeArgs.ContractType), subscribeArgs.Interval)
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsContinuousKlineEvent)
		err := json.Unmarshal(message, event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsCombinedContinuousKlineServe is similar to WsContinuousKlineServe, but it handles multiple pairs of different contractType with its interval
func WsCombinedContinuousKlineServe(subscribeArgsList []*WsContinuousKlineSubscribeArgs,
	handler WsContinuousKlineHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := getCombinedEndpoint()
	for _, val := range subscribeArgsList {
		endpoint += fmt.Sprintf("%s_%s@continuousKline_%s", strings.ToLower(val.Pair),
			strings.ToLower(val.ContractType), val.Interval) + "/"
	}
	endpoint = endpoint[:len(endpoint)-1]
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		j, err := newJSON(message)
		if err != nil {
			errHandler(err)
			return
		}

		data := j.Get("data").MustMap()

		jsonData, _ := json.Marshal(data)

		event := new(WsContinuousKlineEvent)
		err = json.Unmarshal(jsonData, event)
		if err != nil {
			errHandler(err)
			return
		}

		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsMiniMarketTickerEvent define websocket mini market ticker event.
type WsMiniMarketTickerEvent struct {
	Event       string `json:"e"`
	Time        int64  `json:"E"`
	Symbol      string `json:"s"`
	ClosePrice  string `json:"c"`
	OpenPrice   string `json:"o"`
	HighPrice   string `json:"h"`
	LowPrice    string `json:"l"`
	Volume      string `json:"v"`
	QuoteVolume string `json:"q"`
}

// WsMiniMarketTickerHandler handle websocket that pushes 24hr rolling window mini-ticker statistics for a single symbol.
type WsMiniMarketTickerHandler func(event *WsMiniMarketTickerEvent)

// WsMiniMarketTickerServe serve websocket that pushes 24hr rolling window mini-ticker statistics for a single symbol.
func WsMiniMarketTickerServe(symbol string, handler WsMiniMarketTickerHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s@miniTicker", getWsEndpoint(), strings.ToLower(symbol))
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsMiniMarketTickerEvent)
		err := json.Unmarshal(message, &event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsAllMiniMarketTickerEvent define an array of websocket mini market ticker events.
type WsAllMiniMarketTickerEvent []*WsMiniMarketTickerEvent

// WsAllMiniMarketTickerHandler handle websocket that pushes price and funding rate for all markets.
type WsAllMiniMarketTickerHandler func(event WsAllMiniMarketTickerEvent)

// WsAllMiniMarketTickerServe serve websocket that pushes price and funding rate for all markets.
func WsAllMiniMarketTickerServe(handler WsAllMiniMarketTickerHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/!miniTicker@arr", getWsEndpoint())
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		var event WsAllMiniMarketTickerEvent
		err := json.Unmarshal(message, &event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsMarketTickerEvent define websocket market ticker event.
type WsMarketTickerEvent struct {
	Event              string `json:"e"`
	Time               int64  `json:"E"`
	Symbol             string `json:"s"`
	PriceChange        string `json:"p"`
	PriceChangePercent string `json:"P"`
	WeightedAvgPrice   string `json:"w"`
	ClosePrice         string `json:"c"`
	CloseQty           string `json:"Q"`
	OpenPrice          string `json:"o"`
	HighPrice          string `json:"h"`
	LowPrice           string `json:"l"`
	BaseVolume         string `json:"v"`
	QuoteVolume        string `json:"q"`
	OpenTime           int64  `json:"O"`
	CloseTime          int64  `json:"C"`
	FirstID            int64  `json:"F"`
	LastID             int64  `json:"L"`
	TradeCount         int64  `json:"n"`
}

// WsMarketTickerHandler handle websocket that pushes 24hr rolling window mini-ticker statistics for a single symbol.
type WsMarketTickerHandler func(event *WsMarketTickerEvent)

// WsMarketTickerServe serve websocket that pushes 24hr rolling window mini-ticker statistics for a single symbol.
func WsMarketTickerServe(symbol string, handler WsMarketTickerHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s@ticker", getWsEndpoint(), strings.ToLower(symbol))
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsMarketTickerEvent)
		err := json.Unmarshal(message, &event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsAllMarketTickerEvent define an array of websocket mini ticker events.
type WsAllMarketTickerEvent []*WsMarketTickerEvent

// WsAllMarketTickerHandler handle websocket that pushes price and funding rate for all markets.
type WsAllMarketTickerHandler func(event WsAllMarketTickerEvent)

// WsAllMarketTickerServe serve websocket that pushes price and funding rate for all markets.
func WsAllMarketTickerServe(handler WsAllMarketTickerHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/!ticker@arr", getWsEndpoint())
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		var event WsAllMarketTickerEvent
		err := json.Unmarshal(message, &event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsBookTickerEvent define websocket best book ticker event.
type WsBookTickerEvent struct {
	Event           string `json:"e"`
	UpdateID        int64  `json:"u"`
	Time            int64  `json:"E"`
	TransactionTime int64  `json:"T"`
	Symbol          string `json:"s"`
	BestBidPrice    string `json:"b"`
	BestBidQty      string `json:"B"`
	BestAskPrice    string `json:"a"`
	BestAskQty      string `json:"A"`
}

// WsBookTickerHandler handle websocket that pushes updates to the best bid or ask price or quantity in real-time for a specified symbol.
type WsBookTickerHandler func(event *WsBookTickerEvent)

// WsBookTickerServe serve websocket that pushes updates to the best bid or ask price or quantity in real-time for a specified symbol.
func WsBookTickerServe(symbol string, handler WsBookTickerHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s@bookTicker", getWsEndpoint(), strings.ToLower(symbol))
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsBookTickerEvent)
		err := json.Unmarshal(message, &event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsAllBookTickerServe serve websocket that pushes updates to the best bid or ask price or quantity in real-time for all symbols.
func WsAllBookTickerServe(handler WsBookTickerHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/!bookTicker", getWsEndpoint())
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsBookTickerEvent)
		err := json.Unmarshal(message, &event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsLiquidationOrderEvent define websocket liquidation order event.
type WsLiquidationOrderEvent struct {
	Event            string             `json:"e"`
	Time             int64              `json:"E"`
	LiquidationOrder WsLiquidationOrder `json:"o"`
}

// WsLiquidationOrder define websocket liquidation order.
type WsLiquidationOrder struct {
	Symbol               string          `json:"s"`
	Side                 SideType        `json:"S"`
	OrderType            OrderType       `json:"o"`
	TimeInForce          TimeInForceType `json:"f"`
	OrigQuantity         string          `json:"q"`
	Price                string          `json:"p"`
	AvgPrice             string          `json:"ap"`
	OrderStatus          OrderStatusType `json:"X"`
	LastFilledQty        string          `json:"l"`
	AccumulatedFilledQty string          `json:"z"`
	TradeTime            int64           `json:"T"`
}

// WsLiquidationOrderHandler handle websocket that pushes force liquidation order information for specific symbol.
type WsLiquidationOrderHandler func(event *WsLiquidationOrderEvent)

// WsLiquidationOrderServe serve websocket that pushes force liquidation order information for specific symbol.
func WsLiquidationOrderServe(symbol string, handler WsLiquidationOrderHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s@forceOrder", getWsEndpoint(), strings.ToLower(symbol))
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsLiquidationOrderEvent)
		err := json.Unmarshal(message, &event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsAllLiquidationOrderServe serve websocket that pushes force liquidation order information for all symbols.
func WsAllLiquidationOrderServe(handler WsLiquidationOrderHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/!forceOrder@arr", getWsEndpoint())
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsLiquidationOrderEvent)
		err := json.Unmarshal(message, &event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsDepthEvent Partial and Diff. depth event

type WsDepthEvent struct {
	Event            string `json:"e"`
	Time             int64  `json:"E"`
	TransactionTime  int64  `json:"T"`
	Symbol           string `json:"s"`
	FirstUpdateID    int64  `json:"U"`
	LastUpdateID     int64  `json:"u"`
	PrevLastUpdateID int64  `json:"pu"`
	Bids             []Ask  `json:"b"`
	Asks             []Bid  `json:"a"`
}

// WsDepthHandler Partial and Diff. depth stream handler
type WsDepthHandler func(event *WsDepthEvent)

// getDepthStreamName builds the symbol string with levels and rate
func getDepthStreamName(symbol string, levels string, rate string) (string, error) {
	var rateFinal string
	if rate != "" {
		if rate != "100ms" && rate != "250ms" && rate != "500ms" {
			return "", fmt.Errorf("invalid rate")
		}
		rateFinal = fmt.Sprintf("@%s", rate)
	}

	var levelsFinal string
	if levels != "" {
		if levels != "5" && levels != "10" && levels != "20" {
			return "", fmt.Errorf("invalid levels")
		}
		levelsFinal = levels
	}

	return fmt.Sprintf("%s@depth%s%s", strings.ToLower(symbol), levelsFinal, rateFinal), nil
}

func getDepthEndpoint(symbol string, levels string, rate string) (string, error) {
	endpoint, err := getDepthStreamName(symbol, levels, rate)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", getWsEndpoint(), endpoint), nil
}

func getDepthCombinedEndpoint(symbolLevelsRates [][]string) (string, error) {
	streamNames := make([]string, len(symbolLevelsRates))
	for i, slr := range symbolLevelsRates {
		if len(slr) != 3 {
			return "", errors.New("invalid symbol-levels-rate")
		}

		partialEndpoint, err := getDepthStreamName(slr[0], slr[1], slr[2])
		if err != nil {
			return "", err
		}

		streamNames[i] = partialEndpoint
	}

	return getCombinedEndpoint() + strings.Join(streamNames, "/"), nil
}

// wsDepthServe Partial and Diff. depth handler
func wsDepthServe(endpoint string, combined bool, handler WsDepthHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	cfg := newWsConfig(endpoint)
	return wsServe(cfg,
		func(message []byte) {
			type receive struct {
				Event            string     `json:"e"`
				Time             int64      `json:"E"`
				TransactionTime  int64      `json:"T"`
				Symbol           string     `json:"s"`
				FirstUpdateID    int64      `json:"U"`
				LastUpdateID     int64      `json:"u"`
				PrevLastUpdateID int64      `json:"pu"`
				Bids             [][]string `json:"b"`
				Asks             [][]string `json:"a"`
			}

			var received receive
			if combined {
				type combinedWrapper struct {
					Stream string  `json:"stream"`
					Data   receive `json:"data"`
				}
				wrapper := new(combinedWrapper)
				err := json.Unmarshal(message, wrapper)
				if err != nil {
					errHandler(err)
					return
				}
				received = wrapper.Data
			} else {
				err := json.Unmarshal(message, &received)
				if err != nil {
					errHandler(err)
					return
				}
			}

			handler(&WsDepthEvent{
				Event:            received.Event,
				Time:             received.Time,
				TransactionTime:  received.TransactionTime,
				Symbol:           received.Symbol,
				FirstUpdateID:    received.FirstUpdateID,
				LastUpdateID:     received.LastUpdateID,
				PrevLastUpdateID: received.PrevLastUpdateID,
				Bids:             common.JsonToLevels(received.Bids),
				Asks:             common.JsonToLevels(received.Asks),
			})
		},
		errHandler,
	)
}

// WsPartialDepthServe serve websocket partial depth handler with a symbol
func WsPartialDepthServe(symbol string, levels string, handler WsDepthHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint, err := getDepthEndpoint(symbol, levels, "")
	if err != nil {
		return nil, nil, err
	}
	return wsDepthServe(endpoint, false, handler, errHandler)
}

// WsPartialDepthServeWithRate serve websocket partial depth handler with rate.
func WsPartialDepthServeWithRate(symbol string, levels string, rate string, handler WsDepthHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint, err := getDepthEndpoint(symbol, levels, rate)
	if err != nil {
		return nil, nil, err
	}
	return wsDepthServe(endpoint, false, handler, errHandler)
}

// WsCombinedPartialDepthServe is similar to WsPartialDepthServe, but it for multiple symbols
func WsCombinedPartialDepthServe(symbolLevelsRates [][]string, handler WsDepthHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint, err := getDepthCombinedEndpoint(symbolLevelsRates)
	if err != nil {
		return nil, nil, err
	}
	return wsDepthServe(endpoint, true, handler, errHandler)
}

// WsDiffDepthServe serve websocket diff. depth handler.
func WsDiffDepthServe(symbol string, handler WsDepthHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint, err := getDepthEndpoint(symbol, "", "")
	if err != nil {
		return nil, nil, err
	}
	return wsDepthServe(endpoint, false, handler, errHandler)
}

// WsDiffDepthServeWithRate serve websocket diff. depth handler with rate.
func WsDiffDepthServeWithRate(symbol string, rate string, handler WsDepthHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint, err := getDepthEndpoint(symbol, "", rate)
	if err != nil {
		return nil, nil, err
	}
	return wsDepthServe(endpoint, false, handler, errHandler)
}

// WsCombinedDiffDepthServe is similar to WsDiffDepthServe, but it for multiple symbols
func WsCombinedDiffDepthServe(symbolLevelsRates [][]string, handler WsDepthHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint, err := getDepthCombinedEndpoint(symbolLevelsRates)
	if err != nil {
		return nil, nil, err
	}
	return wsDepthServe(endpoint, true, handler, errHandler)
}

// WsBLVTInfoEvent define websocket BLVT info event
type WsBLVTInfoEvent struct {
	Event          string         `json:"e"`
	Time           int64          `json:"E"`
	Symbol         string         `json:"s"`
	Issued         float64        `json:"m"`
	Baskets        []WsBLVTBasket `json:"b"`
	Nav            float64        `json:"n"`
	Leverage       float64        `json:"l"`
	TargetLeverage int64          `json:"t"`
	FundingRate    float64        `json:"f"`
}

// WsBLVTBasket define websocket BLVT basket
type WsBLVTBasket struct {
	Symbol   string `json:"s"`
	Position int64  `json:"n"`
}

// WsBLVTInfoHandler handle websocket BLVT event
type WsBLVTInfoHandler func(event *WsBLVTInfoEvent)

// WsBLVTInfoServe serve BLVT info stream
func WsBLVTInfoServe(name string, handler WsBLVTInfoHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s@tokenNav", getWsEndpoint(), strings.ToUpper(name))
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsBLVTInfoEvent)
		err := json.Unmarshal(message, &event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsBLVTKlineEvent define BLVT kline event
type WsBLVTKlineEvent struct {
	Event  string      `json:"e"`
	Time   int64       `json:"E"`
	Symbol string      `json:"s"`
	Kline  WsBLVTKline `json:"k"`
}

// WsBLVTKline BLVT kline
type WsBLVTKline struct {
	StartTime       int64  `json:"t"`
	CloseTime       int64  `json:"T"`
	Symbol          string `json:"s"`
	Interval        string `json:"i"`
	FirstUpdateTime int64  `json:"f"`
	LastUpdateTime  int64  `json:"L"`
	OpenPrice       string `json:"o"`
	ClosePrice      string `json:"c"`
	HighPrice       string `json:"h"`
	LowPrice        string `json:"l"`
	Leverage        string `json:"v"`
	Count           int64  `json:"n"`
}

// WsBLVTKlineHandler BLVT kline handler
type WsBLVTKlineHandler func(event *WsBLVTKlineEvent)

// WsBLVTKlineServe serve BLVT kline stream
func WsBLVTKlineServe(name string, interval string, handler WsBLVTKlineHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s@nav_Kline_%s", getWsEndpoint(), strings.ToUpper(name), interval)
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsBLVTKlineEvent)
		err := json.Unmarshal(message, event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsCompositeIndexEvent websocket composite index event
type WsCompositeIndexEvent struct {
	Event       string          `json:"e"`
	Time        int64           `json:"E"`
	Symbol      string          `json:"s"`
	Price       string          `json:"p"`
	Composition []WsComposition `json:"c"`
}

// WsComposition websocket composite index event composition
type WsComposition struct {
	BaseAsset    string `json:"b"`
	WeightQty    string `json:"w"`
	WeighPercent string `json:"W"`
}

// WsCompositeIndexHandler websocket composite index handler
type WsCompositeIndexHandler func(event *WsCompositeIndexEvent)

// WsCompositiveIndexServe serve composite index information for index symbols
func WsCompositiveIndexServe(symbol string, handler WsCompositeIndexHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s@compositeIndex", getWsEndpoint(), strings.ToLower(symbol))
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsCompositeIndexEvent)
		err := json.Unmarshal(message, event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}

// WsUserDataEvent define user data event
type WsUserDataEvent struct {
	Event               UserDataEventType     `json:"e"`
	Time                int64                 `json:"E"`
	CrossWalletBalance  string                `json:"cw"`
	MarginCallPositions []WsPosition          `json:"p"`
	TransactionTime     int64                 `json:"T"`
	AccountUpdate       WsAccountUpdate       `json:"a"`
	OrderTradeUpdate    WsOrderTradeUpdate    `json:"o"`
	AccountConfigUpdate WsAccountConfigUpdate `json:"ac"`
}

func (e *WsUserDataEvent) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Event               UserDataEventType     `json:"e"`
		Time                interface{}           `json:"E"`
		CrossWalletBalance  string                `json:"cw"`
		MarginCallPositions []WsPosition          `json:"p"`
		TransactionTime     int64                 `json:"T"`
		AccountUpdate       WsAccountUpdate       `json:"a"`
		OrderTradeUpdate    WsOrderTradeUpdate    `json:"o"`
		AccountConfigUpdate WsAccountConfigUpdate `json:"ac"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	e.Event = tmp.Event
	switch v := tmp.Time.(type) {
	case float64:
		e.Time = int64(v)
	case string:
		parsedTime, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
		e.Time = parsedTime
	default:
		return fmt.Errorf("unexpected type for E: %T", tmp.Time)
	}
	e.CrossWalletBalance = tmp.CrossWalletBalance
	e.MarginCallPositions = tmp.MarginCallPositions
	e.TransactionTime = tmp.TransactionTime
	e.AccountUpdate = tmp.AccountUpdate
	e.OrderTradeUpdate = tmp.OrderTradeUpdate
	e.AccountConfigUpdate = tmp.AccountConfigUpdate
	return nil
}

// WsAccountUpdate define account update
type WsAccountUpdate struct {
	Reason    UserDataEventReasonType `json:"m"`
	Balances  []WsBalance             `json:"B"`
	Positions []WsPosition            `json:"P"`
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
	Symbol                    string           `json:"s"`
	Side                      PositionSideType `json:"ps"`
	Amount                    string           `json:"pa"`
	MarginType                MarginType       `json:"mt"`
	IsolatedWallet            string           `json:"iw"`
	EntryPrice                string           `json:"ep"`
	MarkPrice                 string           `json:"mp"`
	UnrealizedPnL             string           `json:"up"`
	AccumulatedRealized       string           `json:"cr"`
	MaintenanceMarginRequired string           `json:"mm"`
}

// WsOrderTradeUpdate define order trade update
type WsOrderTradeUpdate struct {
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

// WsAccountConfigUpdate define account config update
type WsAccountConfigUpdate struct {
	Symbol   string `json:"s"`
	Leverage int64  `json:"l"`
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
