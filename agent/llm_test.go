package main

import (
	"testing"

	"github.com/farazhassan/gantry/components/llm/anthropic"
	"github.com/farazhassan/gantry/components/llm/ollama"
	"github.com/farazhassan/gantry/components/llm/openai"
	"github.com/farazhassan/gantry/components/llm/openrouter"
)

func fakeEnv(vals map[string]string) func(string) string {
	return func(key string) string { return vals[key] }
}

func TestNewLLMClientOpenRouterDefault(t *testing.T) {
	client, err := newLLMClient("openrouter", fakeEnv(map[string]string{"OPENROUTER_API_KEY": "test-key"}))
	if err != nil {
		t.Fatalf("newLLMClient returned error: %v", err)
	}
	if _, ok := client.(*openrouter.Client); !ok {
		t.Errorf("client type = %T, want *openrouter.Client", client)
	}
}

func TestNewLLMClientEmptyProviderDefaultsToOpenRouter(t *testing.T) {
	client, err := newLLMClient("", fakeEnv(map[string]string{"OPENROUTER_API_KEY": "test-key"}))
	if err != nil {
		t.Fatalf("newLLMClient returned error: %v", err)
	}
	if _, ok := client.(*openrouter.Client); !ok {
		t.Errorf("client type = %T, want *openrouter.Client", client)
	}
}

func TestNewLLMClientOpenRouterMissingKeyErrors(t *testing.T) {
	if _, err := newLLMClient("openrouter", fakeEnv(nil)); err == nil {
		t.Fatal("newLLMClient with no OPENROUTER_API_KEY: want error, got nil")
	}
}

func TestNewLLMClientAnthropic(t *testing.T) {
	client, err := newLLMClient("anthropic", fakeEnv(map[string]string{"ANTHROPIC_API_KEY": "test-key"}))
	if err != nil {
		t.Fatalf("newLLMClient returned error: %v", err)
	}
	if _, ok := client.(*anthropic.Client); !ok {
		t.Errorf("client type = %T, want *anthropic.Client", client)
	}
}

func TestNewLLMClientAnthropicMissingKeyErrors(t *testing.T) {
	if _, err := newLLMClient("anthropic", fakeEnv(nil)); err == nil {
		t.Fatal("newLLMClient with no ANTHROPIC_API_KEY: want error, got nil")
	}
}

func TestNewLLMClientOpenAI(t *testing.T) {
	client, err := newLLMClient("openai", fakeEnv(map[string]string{"OPENAI_API_KEY": "test-key"}))
	if err != nil {
		t.Fatalf("newLLMClient returned error: %v", err)
	}
	if _, ok := client.(*openai.Client); !ok {
		t.Errorf("client type = %T, want *openai.Client", client)
	}
}

func TestNewLLMClientOllamaNoKeyRequired(t *testing.T) {
	client, err := newLLMClient("ollama", fakeEnv(nil))
	if err != nil {
		t.Fatalf("newLLMClient returned error: %v", err)
	}
	if _, ok := client.(*ollama.Client); !ok {
		t.Errorf("client type = %T, want *ollama.Client", client)
	}
}

func TestNewLLMClientUnknownProviderErrors(t *testing.T) {
	if _, err := newLLMClient("carrier-pigeon", fakeEnv(nil)); err == nil {
		t.Fatal("newLLMClient with unknown provider: want error, got nil")
	}
}
