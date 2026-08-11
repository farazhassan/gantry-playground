package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/farazhassan/gantry"
)

// weatherTool is a deterministic, no-network stand-in for a real weather
// API — enough to exercise TOOL_CALL_* frames over AG-UI without external
// dependencies or flakiness.
type weatherTool struct{}

type weatherReading struct {
	Condition string
	TempC     float64
}

var weatherByCity = map[string]weatherReading{
	"san francisco": {"foggy", 15.5},
	"new york":      {"clear", 22.0},
	"london":        {"rainy", 13.0},
	"tokyo":         {"clear", 27.5},
}

func (weatherTool) Definition() gantry.ToolDef {
	return gantry.ToolDef{
		Name:        "get_weather",
		Description: "Returns the current weather for a named city. Demo data only, not live.",
		Schema:      json.RawMessage(`{"type":"object","properties":{"location":{"type":"string","description":"City name, e.g. \"Tokyo\""}},"required":["location"]}`),
	}
}

func (weatherTool) Invoke(_ context.Context, in json.RawMessage) (json.RawMessage, error) {
	var args struct {
		Location string `json:"location"`
	}
	if err := json.Unmarshal(in, &args); err != nil {
		return nil, fmt.Errorf("get_weather: invalid input: %w", err)
	}
	if args.Location == "" {
		return nil, fmt.Errorf("get_weather: location is required")
	}
	w, ok := weatherByCity[strings.ToLower(strings.TrimSpace(args.Location))]
	if !ok {
		w = weatherReading{"sunny", 21.0}
	}
	return json.Marshal(map[string]any{
		"location":  args.Location,
		"condition": w.Condition,
		"temp_c":    w.TempC,
	})
}

// searchDocsTool is a deterministic, no-network stand-in for a real search/
// retrieval tool — a tiny fixed doc index matched by keyword.
type searchDocsTool struct{}

type docHit struct {
	Title   string `json:"title"`
	Snippet string `json:"snippet"`
}

type docEntry struct {
	Keywords []string
	Hit      docHit
}

var docIndex = []docEntry{
	{
		Keywords: []string{"gantry", "agent", "middleware"},
		Hit:      docHit{"Gantry core loop", "Gantry runs a fixed phase loop with onion-style middleware at every phase."},
	},
	{
		Keywords: []string{"agui", "ag-ui", "stream", "sse"},
		Hit:      docHit{"AG-UI streaming", "The agui package maps Gantry's Event stream onto AG-UI's SSE protocol."},
	},
	{
		Keywords: []string{"subagent", "delegate", "nested"},
		Hit:      docHit{"Delegating to sub-agents", "components/subagent wraps a child *gantry.Agent as a tool the parent can call."},
	},
	{
		Keywords: []string{"deploy", "docker", "compose"},
		Hit:      docHit{"Deploying with Docker", "docker compose up brings up the agent and playground together on one shared port."},
	},
}

func (searchDocsTool) Definition() gantry.ToolDef {
	return gantry.ToolDef{
		Name:        "search_docs",
		Description: "Searches a small fixed set of demo docs about this project by keyword. Demo data only.",
		Schema:      json.RawMessage(`{"type":"object","properties":{"query":{"type":"string"}},"required":["query"]}`),
	}
}

func (searchDocsTool) Invoke(_ context.Context, in json.RawMessage) (json.RawMessage, error) {
	var args struct {
		Query string `json:"query"`
	}
	if err := json.Unmarshal(in, &args); err != nil {
		return nil, fmt.Errorf("search_docs: invalid input: %w", err)
	}
	q := strings.ToLower(args.Query)
	var hits []docHit
	for _, entry := range docIndex {
		for _, kw := range entry.Keywords {
			if strings.Contains(q, kw) {
				hits = append(hits, entry.Hit)
				break
			}
		}
	}
	return json.Marshal(map[string]any{"query": args.Query, "results": hits})
}
