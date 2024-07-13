package pmargin

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/adshao/go-binance/v2/futures"
)

// UMPositionRisk define position risk info
type UMPositionRisk struct {
	EntryPrice       string                   `json:"entryPrice"`
	Leverage         string                   `json:"leverage"`
	MarkPrice        string                   `json:"markPrice"`
	MaxNotionalValue string                   `json:"maxNotionalValue"`
	PositionAmt      string                   `json:"positionAmt"`
	Notional         string                   `json:"notional"`
	Symbol           string                   `json:"symbol"`
	UnRealizedProfit string                   `json:"unRealizedProfit"`
	LiquidationPrice string                   `json:"liquidationPrice"`
	PositionSide     futures.PositionSideType `json:"positionSide"`
	UpdateTime       int64                    `json:"updateTime"`
}

// GetUMPositionRiskService get account balance
type GetUMPositionRiskService struct {
	c      *Client
	symbol string
}

// Symbol set symbol
func (s *GetUMPositionRiskService) Symbol(symbol string) *GetUMPositionRiskService {
	s.symbol = symbol
	return s
}

// Do send request
func (s *GetUMPositionRiskService) Do(ctx context.Context, opts ...RequestOption) (res []*UMPositionRisk, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/papi/v1/um/positionRisk",
		secType:  secTypeSigned,
	}
	if s.symbol != "" {
		r.setParam("symbol", s.symbol)
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []*UMPositionRisk{}, err
	}
	res = make([]*UMPositionRisk, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*UMPositionRisk{}, err
	}
	return res, nil
}
