package main

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/farazhassan/gantry"
	"github.com/farazhassan/gantry/components/subagent"
	"github.com/farazhassan/gantry/eval"
)

func TestBuildResearcherDefinition(t *testing.T) {
	llm := eval.NewMockLLMClient(gantry.LLMResponse{Content: "unused", StopReason: gantry.StopReasonEnd})
	researcher, err := buildResearcher(llm)
	if err != nil {
		t.Fatalf("buildResearcher returned error: %v", err)
	}
	if got := researcher.Definition().Name; got != "researcher" {
		t.Errorf("Definition().Name = %q, want %q", got, "researcher")
	}
}

// TestParentDelegatesToResearcher scripts a full round trip using two
// separate mock LLM clients — one per agent, matching gantry's own
// examples/subagent convention — so response ordering isn't coupled to
// delegate-call internals: the parent delegates to the researcher, the
// researcher calls search_docs and reports a finding, and the parent
// reports that finding as its own final answer.
func TestParentDelegatesToResearcher(t *testing.T) {
	researcherLLM := eval.NewMockLLMClient(
		gantry.LLMResponse{
			ToolCalls: []gantry.ToolCall{{
				ID:    "call-docs",
				Name:  "search_docs",
				Input: json.RawMessage(`{"query":"agui stream"}`),
			}},
			StopReason: gantry.StopReasonToolUse,
		},
		gantry.LLMResponse{
			Content:    "the agui package maps events onto AG-UI's SSE protocol",
			StopReason: gantry.StopReasonEnd,
		},
	)
	parentLLM := eval.NewMockLLMClient(
		gantry.LLMResponse{
			ToolCalls: []gantry.ToolCall{{
				ID:    "call-delegate",
				Name:  "researcher",
				Input: json.RawMessage(`{"goal":"find docs about agui streaming"}`),
			}},
			StopReason: gantry.StopReasonToolUse,
		},
		gantry.LLMResponse{
			Content:    "Per the researcher: the agui package maps events onto AG-UI's SSE protocol",
			StopReason: gantry.StopReasonEnd,
		},
	)

	researcher, err := buildResearcher(researcherLLM)
	if err != nil {
		t.Fatalf("buildResearcher returned error: %v", err)
	}
	agent, err := gantry.NewAgent(gantry.WithLLM(parentLLM))
	if err != nil {
		t.Fatalf("NewAgent returned error: %v", err)
	}
	if err := agent.With(subagent.Component(1, weatherTool{}, searchDocsTool{}, researcher)); err != nil {
		t.Fatalf("agent.With returned error: %v", err)
	}

	state, err := agent.Run(context.Background(), "what does agui streaming do?")
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if !strings.Contains(state.FinalOutput, "SSE protocol") {
		t.Errorf("FinalOutput = %q, want it to mention the researcher's finding", state.FinalOutput)
	}
}
