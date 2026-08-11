package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/farazhassan/gantry"
	"github.com/farazhassan/gantry/eval"
)

func TestNewHandlerServesRunOverSSE(t *testing.T) {
	llm := eval.NewMockLLMClient(gantry.LLMResponse{
		Content:    "hello from the agent",
		StopReason: gantry.StopReasonEnd,
	})
	handler, err := newHandler(llm)
	if err != nil {
		t.Fatalf("newHandler returned error: %v", err)
	}

	body, err := json.Marshal(map[string]any{
		"messages": []map[string]string{{"role": "user", "content": "hi"}},
	})
	if err != nil {
		t.Fatalf("json.Marshal returned error: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/agui", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}
	out := rec.Body.String()
	if !strings.Contains(out, `"RUN_STARTED"`) {
		t.Errorf("response missing RUN_STARTED frame: %s", out)
	}
	if !strings.Contains(out, `"RUN_FINISHED"`) {
		t.Errorf("response missing RUN_FINISHED frame: %s", out)
	}
	// eval.MockLLMClient streams Content as a sequence of fixed-size rune
	// chunks (see gantry/eval.chunkSize = 6), so the full string never
	// appears contiguously in the SSE body — reconstruct it from the
	// TEXT_MESSAGE_CONTENT deltas instead of asserting a raw substring.
	if got := assistantText(t, out); got != "hello from the agent" {
		t.Errorf("reconstructed assistant content = %q, want %q (body: %s)", got, "hello from the agent", out)
	}
}

func TestNewHandlerSuspendsOnAskUser(t *testing.T) {
	llm := eval.NewMockLLMClient(gantry.LLMResponse{
		ToolCalls: []gantry.ToolCall{
			{ID: "q1", Name: "ask_user", Input: json.RawMessage(`{"questions":[{"header":"time","text":"What time works?","options":["6pm","7pm"]}]}`)},
		},
		StopReason: gantry.StopReasonToolUse,
	})
	handler, err := newHandler(llm)
	if err != nil {
		t.Fatalf("newHandler returned error: %v", err)
	}

	body, err := json.Marshal(map[string]any{
		"messages": []map[string]string{{"role": "user", "content": "book a table"}},
	})
	if err != nil {
		t.Fatalf("json.Marshal returned error: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/agui", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}
	out := rec.Body.String()
	if !strings.Contains(out, `"TOOL_CALL_START"`) {
		t.Errorf("response missing TOOL_CALL_START frame: %s", out)
	}
	if !strings.Contains(out, `"RUN_FINISHED"`) {
		t.Errorf("response missing RUN_FINISHED frame: %s", out)
	}
	if strings.Contains(out, `"TOOL_CALL_RESULT"`) {
		t.Errorf("ask_user is a client tool — the run must suspend with no TOOL_CALL_RESULT: %s", out)
	}
}

// assistantText concatenates every TEXT_MESSAGE_CONTENT delta in an SSE
// response body, in order, to reconstruct the streamed assistant message.
func assistantText(t *testing.T, sse string) string {
	t.Helper()
	var sb strings.Builder
	for _, line := range strings.Split(sse, "\n") {
		line = strings.TrimPrefix(line, "data: ")
		if !strings.Contains(line, `"TEXT_MESSAGE_CONTENT"`) {
			continue
		}
		var frame struct {
			Delta string `json:"delta"`
		}
		if err := json.Unmarshal([]byte(line), &frame); err != nil {
			t.Fatalf("unmarshal SSE frame: %v (line: %q)", err, line)
		}
		sb.WriteString(frame.Delta)
	}
	return sb.String()
}

func TestHealthzReturnsOK(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	healthzHandler(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestParseOrigins(t *testing.T) {
	got := parseOrigins(" http://a.example , ,http://b.example,")
	want := []string{"http://a.example", "http://b.example"}
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("parseOrigins = %v, want %v", got, want)
	}
}
