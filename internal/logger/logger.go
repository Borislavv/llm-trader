package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"llmtrader/internal/types"
)

// Logger collects all info each cycle.
// Right now: just prints prettified JSON to stdout and keeps it in memory.

type Logger struct {
	history []CycleLog
}

type CycleLog struct {
	TS            time.Time                `json:"ts"`
	Snapshot      types.Snapshot           `json:"snapshot"`
	RawDecision   types.RawLLMDecision     `json:"raw_decision"`
	SafeDecision  types.Decision           `json:"safe_decision"`
	ExecResult    *types.ExecutionResult   `json:"exec_result,omitempty"`
}

func New() *Logger {
	return &Logger{
		history: make([]CycleLog, 0, 1024),
	}
}

func (l *Logger) LogCycle(
	snap types.Snapshot,
	raw types.RawLLMDecision,
	safe types.Decision,
	execRes *types.ExecutionResult,
) {
	entry := CycleLog{
		TS:           time.Now().UTC(),
		Snapshot:     snap,
		RawDecision:  raw,
		SafeDecision: safe,
		ExecResult:   execRes,
	}

	l.history = append(l.history, entry)

	pretty, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		log.Printf("[logger] marshal error: %v\n", err)
		return
	}

	fmt.Println("=======================================================")
	fmt.Println(string(pretty))
	fmt.Println("=======================================================")
}
