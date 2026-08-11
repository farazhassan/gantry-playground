package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/farazhassan/gantry"
	"github.com/farazhassan/gantry/components/ask"
	"github.com/farazhassan/gantry/components/subagent"
	"github.com/farazhassan/gantry/components/systemprompt"
	"github.com/farazhassan/gantry/components/tool"
	"github.com/farazhassan/gantry/components/ui/agui"
)

// askUserNudge steers the demo agent toward ask_user when a request is
// ambiguous, instead of silently guessing — otherwise a small/local model
// has little reason to reach for a tool whose description alone may not be
// directive enough (see the ask-user design doc, §8).
const askUserNudge = "If a request is ambiguous or needs a user preference you can't infer, use the ask_user tool to ask rather than guessing."

// newHandler builds the AG-UI HTTP handler for an agent backed by llm. The
// LLM is a parameter so tests can inject a mock while main() wires the real,
// env-selected client. opts pass straight through to agui.Handler. ask_user
// is registered as a client-side tool (tool.Client): when the model asks a
// question, the run suspends over AG-UI (a tool call with no result) for
// the playground to render and resume
func newHandler(llm gantry.LLMClient, opts ...agui.Option) (http.Handler, error) {
	agent, err := gantry.NewAgent(gantry.WithLLM(llm))
	if err != nil {
		return nil, err
	}

	researcher, err := buildResearcher(llm)
	if err != nil {
		return nil, err
	}

	if err := agent.With(subagent.Component(2, weatherTool{}, searchDocsTool{}, researcher)); err != nil {
		return nil, err
	}
	if err := agent.With(tool.Client(ask.Definition())); err != nil {
		return nil, err
	}
	if err := agent.With(systemprompt.New(askUserNudge)); err != nil {
		return nil, err
	}

	return agui.Handler(agent, opts...), nil
}

func healthzHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func main() {
	provider := envOr("LLM_PROVIDER", "openrouter")
	llm, err := newLLMClient(provider, os.Getenv)
	if err != nil {
		log.Fatalf("configure LLM: %v", err)
	}

	addr := envOr("AGUI_ADDR", ":8080")
	var aguiOpts []agui.Option
	if origins := parseOrigins(os.Getenv("AGUI_ALLOWED_ORIGINS")); len(origins) > 0 {
		aguiOpts = append(aguiOpts, agui.WithAllowedOrigins(origins...))
	}

	handler, err := newHandler(llm, aguiOpts...)
	if err != nil {
		log.Fatalf("build handler: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/agui", handler)
	mux.HandleFunc("/healthz", healthzHandler)

	log.Printf("langport playground agent listening on %s (POST /agui); provider=%s", addr, provider)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// parseOrigins splits AGUI_ALLOWED_ORIGINS on commas, trimming whitespace
// and dropping empty entries — a bare comma-split would silently turn
// "http://a, http://b" into a never-matching " http://b" (leading space) or
// turn a trailing/doubled comma into an empty-string "origin".
func parseOrigins(s string) []string {
	var out []string
	for _, o := range strings.Split(s, ",") {
		if o = strings.TrimSpace(o); o != "" {
			out = append(out, o)
		}
	}
	return out
}
