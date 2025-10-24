package types

import "time"

// Snapshot is the "world state" we feed into the LLM.
type Snapshot struct {
	Timestamp time.Time  `json:"timestamp"`
	Account   Account    `json:"account"`
	Positions []Position `json:"positions"`
	Markets   []Market   `json:"markets"`
}

type Account struct {
	BalanceUSDT      float64 `json:"balance_usdt"`
	EquityUSDT       float64 `json:"equity_usdt"`
	FreeMarginUSDT   float64 `json:"free_margin_usdt"`
}

type Position struct {
	Symbol            string  `json:"symbol"`
	Side              string  `json:"side"` // LONG / SHORT
	SizeUSDT          float64 `json:"size_usdt"`
	EntryPrice        float64 `json:"entry_price"`
	Leverage          int     `json:"leverage"`
	UnrealizedPNLUSDT float64 `json:"unrealized_pnl_usdt"`
}

type Market struct {
	Symbol        string  `json:"symbol"`
	LastPrice     float64 `json:"last_price"`
	Change1mPct   float64 `json:"change_1m_pct"`
	Change5mPct   float64 `json:"change_5m_pct"`
	Change1hPct   float64 `json:"change_1h_pct"`
	FundingRate   float64 `json:"funding_rate"`
}

// ----- Decisions / Actions from the LLM -----

const (
	ActionOpenLong  = "OPEN_LONG"
	ActionOpenShort = "OPEN_SHORT"
	ActionClose     = "CLOSE"
	ActionHold      = "HOLD"
)

// RawLLMDecision is the direct JSON we expect from Qwen3-Max.
type RawLLMDecision struct {
	Action          string   `json:"action"`
	Symbol          string   `json:"symbol,omitempty"`
	SizeUSDT        float64  `json:"size_usdt,omitempty"`
	Leverage        int      `json:"leverage,omitempty"`
	StopLossPct     float64  `json:"stop_loss_pct,omitempty"`
	TakeProfitPct   float64  `json:"take_profit_pct,omitempty"`
	PercentToClose  float64  `json:"percent_to_close,omitempty"` // for CLOSE, [0..100]
	Reason          string   `json:"reason"`
}

// Decision is the sanitized, risk-approved command that executor will try to run.
type Decision struct {
	Action         string  `json:"action"`
	Symbol         string  `json:"symbol,omitempty"`
	SizeUSDT       float64 `json:"size_usdt,omitempty"`
	Leverage       int     `json:"leverage,omitempty"`
	StopLossPct    float64 `json:"stop_loss_pct,omitempty"`
	TakeProfitPct  float64 `json:"take_profit_pct,omitempty"`
	PercentToClose float64 `json:"percent_to_close,omitempty"`
	Reason         string  `json:"reason"`
	// RiskNotes can include "leverage clamped from 25->10"
	RiskNotes      string  `json:"risk_notes,omitempty"`
}

// ExecutionResult is what the exchange tells us after we attempt to act.
type ExecutionResult struct {
	Status            string    `json:"status"` // e.g. "FILLED", "REJECTED", "SIMULATED"
	Action            string    `json:"action"`
	Symbol            string    `json:"symbol,omitempty"`
	RequestedSizeUSDT float64   `json:"requested_size_usdt,omitempty"`
	ExecutedSizeUSDT  float64   `json:"executed_size_usdt,omitempty"`
	AvgFillPrice      float64   `json:"avg_fill_price,omitempty"`
	BybitOrderID      string    `json:"bybit_order_id,omitempty"`
	TS                time.Time `json:"ts"`
}
