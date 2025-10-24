# LLM Trading Agent (Prototype)

## Overview

Minimal PoC of an autonomous trading agent:
- pulls snapshot from Bybit testnet (`collector`)
- asks Qwen3-Max for a decision (`agent`)
- validates it (`risk`)
- executes order (`executor`)
- logs everything (`logger`)

**This is an offline skeleton**:
- `collector` has stubs where you must integrate Bybit testnet REST/WebSocket SDK.
- `agent` currently returns a hardcoded mock decision. Replace `MockLLM()` with real Qwen3-Max API call.
- `executor` is a dry-run executor. It only prints/logs instead of placing a real order. Hook Bybit testnet here.
- The loop in `cmd/agent/main.go` runs every 60s.

This repo is structured so you can:
1. drop in your Bybit client,
2. drop in your Qwen3-Max client,
3. instantly get an end-to-end loop.

### Directory structure

- `cmd/agent/main.go`  
  main loop / orchestration.

- `internal/types/`  
  shared data structures and JSON contracts between modules.

- `internal/collector/collector.go`  
  fetch account/market snapshot from Bybit testnet.

- `internal/agent/agent.go`  
  build system prompt and call Qwen3-Max to get trading decision.

- `internal/risk/risk.go`  
  risk filter to clamp leverage, size, etc.

- `internal/executor/executor.go`  
  submit orders (currently dry-run).

- `internal/logger/logger.go`  
  logs every cycle to stdout and in-memory history.

## How to run

Requirements:
- Go >= 1.23 (or recent Go)
- Internet access OFF is okay for now because we have no external deps.

```bash
cd cmd/agent
go run .
```

What you'll see:
- on each tick, snapshot is generated (mock data),
- agent proposes a trade decision,
- risk filter adjusts it,
- executor "executes" (prints),
- logger outputs the cycle summary.

## TODO (where you wire real world)

### Connect Bybit testnet
In `internal/collector/collector.go`:
- implement `FetchSnapshot()`:
  - call Bybit testnet REST API
  - fill `types.Snapshot`

In `internal/executor/executor.go`:
- implement `Execute()`:
  - translate `types.Decision` into REST order:
    - open long / short (perp)
    - close
    - attach stop-loss / take-profit
  - return `types.ExecutionResult`.

### Connect Qwen3-Max
In `internal/agent/agent.go`:
- replace `MockLLMDecision()` with an actual call to Qwen3-Max.
- build the system prompt from `snapshot`.
- enforce that model responds with valid JSON matching `types.RawLLMDecision`.

### Persistent logging
Right now `logger` prints to stdout.
You probably want:
- append to a JSONL file,
- push to Postgres / Grafana later.

## Safety Notes

- `risk.FilterDecision` is **critical**.  
  It clamps leverage, position size, and stop-loss distance to avoid suicidal trades.
- You should expand it with:
  - daily max loss cutoff
  - cooldown after liquidation events
  - per-symbol position uniqueness

## License
This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0).

Commercial use of this software requires a separate commercial license.
Please contact glazunov2142@gmail.com for details.
