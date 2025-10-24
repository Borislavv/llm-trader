package risk

import (
	"fmt"

	"llmtrader/internal/types"
)

// RiskManager validates / clamps the LLM decision
// according to risk rules before we try to execute.

type RiskManager struct {
	MaxLeverage         int
	MaxPositionUSDT     float64
	MaxStopLossPct      float64 // e.g. -2.0 means we won't allow wider than -2%
	MinTakeProfitPct    float64
	MaxTakeProfitPct    float64
}

func New() *RiskManager {
	return &RiskManager{
		MaxLeverage:      10,
		MaxPositionUSDT:  2000,
		MaxStopLossPct:   -2.0,
		MinTakeProfitPct: 1.0,
		MaxTakeProfitPct: 5.0,
	}
}

func (r *RiskManager) FilterDecision(raw types.RawLLMDecision, snap types.Snapshot) types.Decision {
	// default HOLD if nonsense
	if raw.Action == "" {
		return types.Decision{
			Action:  types.ActionHold,
			Reason:  "empty action from LLM",
			RiskNotes: "LLM returned no action, forcing HOLD",
		}
	}

	// COPY raw as base
	dec := types.Decision{
		Action:         raw.Action,
		Symbol:         raw.Symbol,
		SizeUSDT:       raw.SizeUSDT,
		Leverage:       raw.Leverage,
		StopLossPct:    raw.StopLossPct,
		TakeProfitPct:  raw.TakeProfitPct,
		PercentToClose: raw.PercentToClose,
		Reason:         raw.Reason,
		RiskNotes:      "",
	}

	// If HOLD -> no further checks.
	if dec.Action == types.ActionHold {
		return dec
	}

	// BASIC sanity: allowed actions only
	switch dec.Action {
	case types.ActionOpenLong, types.ActionOpenShort, types.ActionClose:
	default:
		// Unknown action => HOLD
		return types.Decision{
			Action:    types.ActionHold,
			Reason:    "Invalid action from LLM",
			RiskNotes: fmt.Sprintf("LLM tried invalid action %q", dec.Action),
		}
	}

	// Clamp leverage
	if dec.Leverage > r.MaxLeverage {
		dec.RiskNotes += fmt.Sprintf("leverage clamped %d->%d; ", dec.Leverage, r.MaxLeverage)
		dec.Leverage = r.MaxLeverage
	}

	// Clamp size
	if dec.SizeUSDT > r.MaxPositionUSDT {
		dec.RiskNotes += fmt.Sprintf("size_usdt clamped %.2f->%.2f; ", dec.SizeUSDT, r.MaxPositionUSDT)
		dec.SizeUSDT = r.MaxPositionUSDT
	}

	// If OPEN_* we require stop_loss and take_profit sanity
	if dec.Action == types.ActionOpenLong || dec.Action == types.ActionOpenShort {
		// stop_loss can't be looser than MaxStopLossPct
		if dec.StopLossPct < r.MaxStopLossPct {
			dec.RiskNotes += fmt.Sprintf("stop_loss_pct clamped %.2f->%.2f; ", dec.StopLossPct, r.MaxStopLossPct)
			dec.StopLossPct = r.MaxStopLossPct
		}
		// TP must be within [MinTakeProfitPct, MaxTakeProfitPct]
		if dec.TakeProfitPct < r.MinTakeProfitPct {
			dec.RiskNotes += fmt.Sprintf("take_profit_pct clamped %.2f->%.2f; ", dec.TakeProfitPct, r.MinTakeProfitPct)
			dec.TakeProfitPct = r.MinTakeProfitPct
		}
		if dec.TakeProfitPct > r.MaxTakeProfitPct {
			dec.RiskNotes += fmt.Sprintf("take_profit_pct clamped %.2f->%.2f; ", dec.TakeProfitPct, r.MaxTakeProfitPct)
			dec.TakeProfitPct = r.MaxTakeProfitPct
		}
	}

	// TODO: daily loss cap, margin checks, side duplication checks, etc.

	return dec
}
