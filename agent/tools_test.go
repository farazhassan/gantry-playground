package main

import (
	"context"
	"encoding/json"
	"testing"
)

func TestToolDefinitions(t *testing.T) {
	if got := (weatherTool{}).Definition().Name; got != "get_weather" {
		t.Errorf("weatherTool name = %q, want %q", got, "get_weather")
	}
	if got := (searchDocsTool{}).Definition().Name; got != "search_docs" {
		t.Errorf("searchDocsTool name = %q, want %q", got, "search_docs")
	}
}

func TestWeatherToolKnownCity(t *testing.T) {
	out, err := (weatherTool{}).Invoke(context.Background(), json.RawMessage(`{"location":"Tokyo"}`))
	if err != nil {
		t.Fatalf("Invoke returned error: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("Invoke output is not valid JSON: %v", err)
	}
	if got["condition"] != "clear" {
		t.Errorf("condition = %v, want %q", got["condition"], "clear")
	}
}

func TestWeatherToolUnknownCityFallsBackToDefault(t *testing.T) {
	out, err := (weatherTool{}).Invoke(context.Background(), json.RawMessage(`{"location":"Nowhereville"}`))
	if err != nil {
		t.Fatalf("Invoke returned error: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("Invoke output is not valid JSON: %v", err)
	}
	if got["condition"] != "sunny" {
		t.Errorf("condition = %v, want fallback %q", got["condition"], "sunny")
	}
}

func TestWeatherToolMissingLocationErrors(t *testing.T) {
	if _, err := (weatherTool{}).Invoke(context.Background(), json.RawMessage(`{}`)); err == nil {
		t.Fatal("Invoke with no location: want error, got nil")
	}
}

func TestSearchDocsToolMatchesKeyword(t *testing.T) {
	out, err := (searchDocsTool{}).Invoke(context.Background(), json.RawMessage(`{"query":"how does the agui stream work"}`))
	if err != nil {
		t.Fatalf("Invoke returned error: %v", err)
	}
	var got struct {
		Results []docHit `json:"results"`
	}
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("Invoke output is not valid JSON: %v", err)
	}
	if len(got.Results) == 0 {
		t.Fatal("Results is empty, want at least one hit for \"agui\"")
	}
}

func TestSearchDocsToolNoMatch(t *testing.T) {
	out, err := (searchDocsTool{}).Invoke(context.Background(), json.RawMessage(`{"query":"unrelated banana bread recipe"}`))
	if err != nil {
		t.Fatalf("Invoke returned error: %v", err)
	}
	var got struct {
		Results []docHit `json:"results"`
	}
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("Invoke output is not valid JSON: %v", err)
	}
	if len(got.Results) != 0 {
		t.Errorf("Results = %v, want empty for an unrelated query", got.Results)
	}
}
