# Langport playground agent

A Gantry agent exposed over AG-UI at `POST /agui`, with two demo tools
(`get_weather`, `search_docs`) and a `researcher` sub-agent.

## Run locally

```bash
cd agent
OPENROUTER_API_KEY=sk-... go run -ldflags=-linkmode=external .
```

(The `-ldflags=-linkmode=external` flag works around a known macOS
Go-toolchain issue — see gantry's own README — and is a no-op on Linux.)

## Environment variables

| Var | Default | Purpose |
| --- | --- | --- |
| `LLM_PROVIDER` | `openrouter` | One of `anthropic`, `openai`, `ollama`, `openrouter`. |
| `OPENROUTER_API_KEY` / `OPENROUTER_MODEL` | — / `anthropic/claude-sonnet-4-6` | Required when provider is `openrouter`. |
| `ANTHROPIC_API_KEY` / `ANTHROPIC_MODEL` | — / `claude-sonnet-4-5` | Required when provider is `anthropic`. |
| `OPENAI_API_KEY` / `OPENAI_MODEL` | — / `gpt-4o` | Required when provider is `openai`. |
| `OLLAMA_HOST` / `OLLAMA_MODEL` | `http://localhost:11434` / `llama3.2` | Used when provider is `ollama`; no key needed. |
| `AGUI_ADDR` | `:8080` | Listen address. |
| `AGUI_ALLOWED_ORIGINS` | unset (CORS disabled) | Comma-separated allowed origins, or `*`. |

## Test

```bash
cd agent && go test -ldflags=-linkmode=external ./... -v
```

## Try it

```bash
curl -N -X POST http://localhost:8080/agui \
  -H 'Content-Type: application/json' \
  -d '{"messages":[{"role":"user","content":"What is the weather in Tokyo, and can you find docs about agui streaming?"}]}'
```
