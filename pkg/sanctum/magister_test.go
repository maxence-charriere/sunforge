package sanctum

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewScribe(t *testing.T) {
	const prompt = "Read the HTML and return an article as JSON."
	in := map[string]string{"body": "<h2>Mage</h2><p>Damage increased by 5%.</p>"}
	called := false
	client := testOpenAIClient(func(r *http.Request) (*http.Response, error) {
		called = true
		var body struct {
			Model           string `json:"model"`
			MaxOutputTokens int    `json:"max_output_tokens"`
			Instructions    string `json:"instructions"`
			Reasoning       struct {
				Effort string `json:"effort"`
			} `json:"reasoning"`
			Input []struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"input"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		require.Equal(t, "gpt-6.2-astra", body.Model)
		require.Equal(t, "medium", body.Reasoning.Effort)
		require.Equal(t, 32768, body.MaxOutputTokens)
		require.Equal(t, prompt, body.Instructions)
		require.Len(t, body.Input, 1)
		require.Equal(t, "user", body.Input[0].Role)
		encoded, err := json.Marshal(in)
		require.NoError(t, err)
		require.Equal(t, "Return JSON only.\n"+string(encoded), body.Input[0].Content)
		return openAIResponse(http.StatusOK, `{"output":[{"type":"message","content":[{"type":"output_text","text":"{\"relevant\":true,\"title\":\"Mage tuning\",\"text\":\"Damage increased by 5%.\"}"}]}]}`), nil
	})

	var out struct {
		Relevant bool   `json:"relevant"`
		Title    string `json:"title"`
		Text     string `json:"text"`
	}
	err := NewScribe(client).Ask(context.Background(), in, &out, prompt)
	require.NoError(t, err)
	require.True(t, called)
	require.True(t, out.Relevant)
	require.Equal(t, "Mage tuning", out.Title)
	require.Equal(t, "Damage increased by 5%.", out.Text)
}
