package executor

import (
	"fmt"
	"time"

	"llmtrader/internal/types"
)

// Executor is responsible for actually placing/modifying orders
// on Bybit testnet. Right now it's a dry-run simulator that just
// returns a mock fill.
//
// TODO: wire real Bybit testnet REST calls here.
type Executor struct {
	// here you'd store API keys, http client, etc.
}

func New() *Executor {
	return &Executor{}
}

func (e *Executor) Execute(dec types.Decision) (types.ExecutionResult, error) {
	now := time.Now().UTC()

	// HOLD should never even call us, but just in case:
	if dec.Action == types.ActionHold {
		return types.ExecutionResult{
			Status:       "SKIPPED",
			Action:       dec.Action,
			Symbol:       dec.Symbol,
			TS:           now,
		}, nil
	}

	// This is where you'd map:
	//  - OPEN_LONG -> /private/order/create?side=Buy ...
	//  - OPEN_SHORT -> side=Sell ...
	//  - CLOSE -> reduce-only order
	//
	// For now just simulate.
	res := types.ExecutionResult{
		Status:            "FILLED_SIMULATION",
		Action:            dec.Action,
		Symbol:            dec.Symbol,
		RequestedSizeUSDT: dec.SizeUSDT,
		ExecutedSizeUSDT:  dec.SizeUSDT * 0.99,
		AvgFillPrice:      1234.56,
		BybitOrderID:      "SIM-"+fmt.Sprint(time.Now().UnixNano()),
		TS:                now,
	}

	return res, nil
}
