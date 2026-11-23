// Package sanctum provides AI agents that turn JSON input into structured output.
package sanctum

import (
	"context"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/shared"
)

// Magister answers prompts using input data and decodes the answer into out.
type Magister interface {
	// Ask encodes in as JSON and decodes the JSON answer into out, which must be
	// a non-nil pointer. The context controls request cancellation and deadlines.
	Ask(ctx context.Context, in any, out any, prompt string) error

	// Askf is Ask with a fmt.Sprintf-formatted prompt.
	Askf(ctx context.Context, in any, out any, promptFormat string, v ...any) error
}

// NewScribe returns a Magister for source analysis and writing, backed by
// GPT-6.2 Sol with medium reasoning effort. Supply editorial instructions in
// the prompt passed to Ask or Askf.
// Output is capped at 32,768 tokens, including reasoning tokens.
// The caller's context controls request deadlines.
func NewScribe(client openai.Client) Magister {
	return &openAIAgent{
		Client:          client,
		Model:           string(shared.ChatModelGPT6_2Sol),
		MaxToken:        32768,
		ReasoningEffort: openai.ReasoningEffortMedium,
	}
}
