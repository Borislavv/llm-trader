package agent

import (
	"encoding/json"
	"fmt"

	"llmtrader/internal/types"
)

// Agent is responsible for:
//  - constructing prompt for Qwen3-Max
//  - sending it to the model
//  - parsing model JSON answer into types.RawLLMDecision
//
// For now, we return a mock decision to OPEN_LONG ETHUSDT.
type Agent struct {
	// here you would hold LLM API credentials, http client, etc.
}

func New() *Agent {
	return &Agent{}
}

func (a *Agent) GetLLMDecision(snap types.Snapshot) (types.RawLLMDecision, error) {
	// TODO: build prompt and query Qwen3-Max.
	// We'll mock it now by generating JSON that looks like LLM output.

	mock := `{
		"action": "OPEN_LONG",
		"symbol": "ETHUSDT",
		"size_usdt": 1500,
		"leverage": 12,
		"stop_loss_pct": -3.5,
		"take_profit_pct": 4.0,
		"reason": "Momentum breakout, acceptable funding"
	}`

	var d types.RawLLMDecision
	if err := json.Unmarshal([]byte(mock), &d); err != nil {
		return types.RawLLMDecision{}, fmt.Errorf("failed to unmarshal mock llm decision: %w", err)
	}
	return d, nil
}
