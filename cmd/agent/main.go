package main

import (
	"log"
	"time"

	"llmtrader/internal/collector"
	"llmtrader/internal/agent"
	"llmtrader/internal/risk"
	"llmtrader/internal/executor"
	"llmtrader/internal/logger"
	"llmtrader/internal/types"
)

func main() {
	log.Println("[main] LLM Trading Agent starting up")

	// Simple in-memory logger impl
	l := logger.New()

	// Create module instances (in real world you'd inject API keys/configs)
	coll := collector.New()
	agt := agent.New()
	riskMgr := risk.New()
	exec := executor.New()

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		// 1. snapshot from market/account
		snap, err := coll.FetchSnapshot()
		if err != nil {
			log.Printf("[main] collector error: %v\n", err)
			continue
		}

		// 2. ask LLM for decision
		rawDecision, err := agt.GetLLMDecision(snap)
		if err != nil {
			log.Printf("[main] agent error: %v\n", err)
			continue
		}

		// 3. risk filter
		safeDecision := riskMgr.FilterDecision(rawDecision, snap)

		var execResult *types.ExecutionResult
		// 4. execute if not HOLD
		if safeDecision.Action != types.ActionHold {
			r, err := exec.Execute(safeDecision)
			if err != nil {
				log.Printf("[main] executor error: %v\n", err)
			} else {
				execResult = &r
			}
		}

		// 5. log cycle
		l.LogCycle(snap, rawDecision, safeDecision, execResult)

		<-ticker.C
	}
}
