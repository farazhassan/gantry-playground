# gantry-playground

A playground for testing [Gantry](https://github.com/farazhassan/gantry)-made
agents end-to-end: a small demo agent exposed over
[AG-UI](https://ag-ui.com/), paired with
[langport](https://github.com/langport-dev/langport), a hosted UI that talks
to any AG-UI-compatible agent over SSE.

Two services, run together with Docker Compose:

- **`agent`** ([agent/](agent/)) — a Gantry agent built from this repo,
  serving `POST /agui`. Has two demo tools (`get_weather`, `search_docs`), a
  `researcher` sub-agent, and a client-side `ask_user` tool for
  human-in-the-loop prompts. Supports Anthropic, OpenAI, Ollama, and
  OpenRouter as LLM backends.
- **`langport`** — the playground UI (`ghcr.io/langport-dev/langport`),
  which renders the agent's AG-UI stream: messages, tool calls, and the
  Subagents timeline.

## Quickstart

```bash
cp .env.example .env
# edit .env and set at least one provider's API key (OPENROUTER_API_KEY by default)
docker compose up --build
```

Then open http://localhost:3000 for the langport UI. The agent itself listens
on http://localhost:8080 (`POST /agui`, `GET /healthz`).

## Configuration

All variables below go in the root `.env` (used by `docker compose up`); see
[.env.example](.env.example) for the full list with defaults.

| Var | Default | Purpose |
| --- | --- | --- |
| `AGUI_PORT` | `8080` | Port the agent listens on and langport connects to. |
| `LLM_PROVIDER` | `openrouter` | One of `anthropic`, `openai`, `ollama`, `openrouter`. |
| `OPENROUTER_API_KEY` / `OPENROUTER_MODEL` | — / `anthropic/claude-sonnet-4-6` | Required when provider is `openrouter`. |
| `ANTHROPIC_API_KEY` / `ANTHROPIC_MODEL` | — / `claude-sonnet-4-5` | Required when provider is `anthropic`. |
| `OPENAI_API_KEY` / `OPENAI_MODEL` | — / `gpt-4o` | Required when provider is `openai`. |
| `OLLAMA_HOST` / `OLLAMA_MODEL` | `http://localhost:11434` / `llama3.2` | Used when provider is `ollama`; no key needed. |
| `AGUI_ALLOWED_ORIGINS` | `*` | Comma-separated CORS origins the agent accepts. |

## Running the agent without Docker

See [agent/README.md](agent/README.md) for running and testing the agent
directly with `go run` / `go test`.

## Project layout

```
.
├── docker-compose.yml   # agent + langport playground stack
├── .env.example          # config for docker compose up
└── agent/                 # the Gantry agent itself (Go module)
```

## License

[MIT](LICENSE)
