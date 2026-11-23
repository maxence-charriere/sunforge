// Package sanctum provides AI agents that turn JSON input into structured output.
package sanctum

import "context"

// Magister answers prompts using input data and decodes the answer into out.
type Magister interface {
	// Ask encodes in as JSON and decodes the JSON answer into out, which must be
	// a non-nil pointer. The context controls request cancellation and deadlines.
	Ask(ctx context.Context, in any, out any, prompt string) error

	// Askf is Ask with a fmt.Sprintf-formatted prompt.
	Askf(ctx context.Context, in any, out any, promptFormat string, v ...any) error
}
