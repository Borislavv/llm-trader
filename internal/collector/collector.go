package collector

import (
	"time"

	"llmtrader/internal/types"
)

// Collector is responsible for fetching current snapshot from Bybit testnet.
// In this skeleton, we mock data. Replace with real Bybit API calls.
type Collector struct {
	// here you would hold API keys / http client / ws client etc.
}

func New() *Collector {
	return &Collector{}
}

func (c *Collector) FetchSnapshot() (types.Snapshot, error) {
	// TODO: integrate Bybit testnet REST+WS.
	// Currently returns dummy snapshot.

	now := time.Now().UTC()

	snap := types.Snapshot{
		Timestamp: now,
		Account: types.Account{
			BalanceUSDT:    10000.0,
			EquityUSDT:     10040.5,
			FreeMarginUSDT: 7000.0,
		},
		Positions: []types.Position{
			{
				Symbol:            "ETHUSDT",
				Side:              "LONG",
				SizeUSDT:          1500.0,
				EntryPrice:        3280.5,
				Leverage:          8,
				UnrealizedPNLUSDT: 220.4,
			},
		},
		Markets: []types.Market{
			{
				Symbol:      "BTCUSDT",
				LastPrice:   66810.0,
				Change1mPct: 0.42,
				Change5mPct: 0.86,
				Change1hPct: 2.1,
				FundingRate: 0.009,
			},
			{
				Symbol:      "ETHUSDT",
				LastPrice:   3295.0,
				Change1mPct: 0.15,
				Change5mPct: 0.50,
				Change1hPct: 1.3,
				FundingRate: 0.012,
			},
		},
	}

	return snap, nil
}
