package main

import (
	"fmt"

	"github.com/farazhassan/gantry"
	"github.com/farazhassan/gantry/components/llm/anthropic"
	"github.com/farazhassan/gantry/components/llm/ollama"
	"github.com/farazhassan/gantry/components/llm/openai"
	"github.com/farazhassan/gantry/components/llm/openrouter"
)

// newLLMClient builds the configured gantry.LLMClient for provider, reading
// model/API key through getenv rather than relying on each adapter's
// internal os.Getenv fallback. That keeps provider selection testable with a
// fake env and lets us return a clean error instead of the adapter's panic
// when a required key is missing.
func newLLMClient(provider string, getenv func(string) string) (gantry.LLMClient, error) {
	switch provider {
	case "", "openrouter":
		key := getenv("OPENROUTER_API_KEY")
		if key == "" {
			return nil, fmt.Errorf("newLLMClient: OPENROUTER_API_KEY is required for provider %q", "openrouter")
		}
		model := getenvOr(getenv, "OPENROUTER_MODEL", "anthropic/claude-sonnet-4-6")
		return openrouter.New(model, openrouter.WithAPIKey(key)), nil

	case "anthropic":
		key := getenv("ANTHROPIC_API_KEY")
		if key == "" {
			return nil, fmt.Errorf("newLLMClient: ANTHROPIC_API_KEY is required for provider %q", "anthropic")
		}
		model := getenvOr(getenv, "ANTHROPIC_MODEL", "claude-sonnet-4-5")
		return anthropic.New(model, anthropic.WithAPIKey(key)), nil

	case "openai":
		key := getenv("OPENAI_API_KEY")
		if key == "" {
			return nil, fmt.Errorf("newLLMClient: OPENAI_API_KEY is required for provider %q", "openai")
		}
		model := getenvOr(getenv, "OPENAI_MODEL", "gpt-4o")
		return openai.New(model, openai.WithAPIKey(key)), nil

	case "ollama":
		model := getenvOr(getenv, "OLLAMA_MODEL", "llama3.2")
		var opts []ollama.Option
		if base := getenv("OLLAMA_HOST"); base != "" {
			opts = append(opts, ollama.WithBaseURL(base))
		}
		return ollama.New(model, opts...), nil

	default:
		return nil, fmt.Errorf("newLLMClient: unknown LLM_PROVIDER %q (want anthropic, openai, ollama, or openrouter)", provider)
	}
}

func getenvOr(getenv func(string) string, key, fallback string) string {
	if v := getenv(key); v != "" {
		return v
	}
	return fallback
}
