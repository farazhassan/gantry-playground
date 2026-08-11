package main

import (
	"github.com/farazhassan/gantry"
	"github.com/farazhassan/gantry/components/subagent"
	"github.com/farazhassan/gantry/components/tool"
)

// buildResearcher constructs the "researcher" sub-agent and wraps it as a
// delegate tool.Tool the parent agent can call. It shares the caller's LLM
// client — a sub-agent is just another *gantry.Agent, so no separate
// provider config is needed — and gets its own search_docs tool so it has
// something to actually do when delegated to (and something to show in the
// playground's Subagents timeline besides a single LLM turn).
//
// WithEventPassthrough is required: without it, the child run's own
// phase/tool/text events are scoped away from the ambient EventSink and
// never reach the AG-UI stream, so the playground's Subagents tab would
// stay empty even though delegation itself still works.
func buildResearcher(llm gantry.LLMClient) (tool.Tool, error) {
	researcher, err := gantry.NewAgent(gantry.WithLLM(llm), gantry.WithName("researcher"))
	if err != nil {
		return nil, err
	}
	if err := researcher.With(tool.FromTools(1, searchDocsTool{})); err != nil {
		return nil, err
	}
	return subagent.New(
		"researcher",
		"Delegates focused research to a specialist sub-agent with its own search_docs tool. Give it a specific goal.",
		researcher,
		subagent.WithEventPassthrough(),
	), nil
}
