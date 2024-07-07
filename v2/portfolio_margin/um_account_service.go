package portfolio_margin

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/adshao/go-binance/v2/futures"
)

type UMAccount struct {
	TradeGroupId int64               `json:"tradeGroupId"`
	Assets       []UMAccountAsset    `json:"assets"`
	Positions    []UMAccountPosition `json:"positions"`
}

type UMAccountAsset struct {
	Asset                  string `json:"asset"`
	CrossWalletBalance     string `json:"crossWalletBalance"`
	CrossUnPnl             string `json:"crossUnPnl"`
	MaintMargin            string `json:"maintMargin"`
	InitialMargin          string `json:"initialMargin"`
	PositionInitialMargin  string `json:"positionInitialMargin"`
	OpenOrderInitialMargin string `json:"openOrderInitialMargin"`
	UpdateTime             int64  `json:"updateTime"`
}

type UMAccountPosition struct {
	Symbol                 string                   `json:"symbol"`
	InitialMargin          string                   `json:"initialMargin"`
	MaintMargin            string                   `json:"maintMargin"`
	UnrealizedProfit       string                   `json:"unrealizedProfit"`
	PositionInitialMargin  string                   `json:"positionInitialMargin"`
	OpenOrderInitialMargin string                   `json:"openOrderInitialMargin"`
	Leverage               string                   `json:"leverage"`
	EntryPrice             string                   `json:"entryPrice"`
	MaxNotional            string                   `json:"maxNotional"`
	BidNotional            string                   `json:"bidNotional"`
	AskNotional            string                   `json:"askNotional"`
	PositionSide           futures.PositionSideType `json:"positionSide"`
	PositionAmt            string                   `json:"positionAmt"`
	UpdateTime             int64                    `json:"updateTime"`
	BreakEvenPrice         string                   `json:"breakEvenPrice"`
}

// GetAccountService get account info
type GetUMAccountService struct {
	c *Client
}

// Do send request
func (s *GetUMAccountService) Do(ctx context.Context, opts ...RequestOption) (res *UMAccount, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/papi/v1/um/account",
		secType:  secTypeSigned,
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(UMAccount)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
