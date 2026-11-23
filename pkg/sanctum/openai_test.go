package sanctum

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/maxence-charriere/go-app/v11/pkg/errors"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/stretchr/testify/require"
)

type openAITransport func(*http.Request) (*http.Response, error)

func (f openAITransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func testOpenAIAgent(client openai.Client, model string, maxToken int) Magister {
	return &openAIAgent{
		Client:          client,
		Model:           model,
		MaxToken:        maxToken,
		ReasoningEffort: openai.ReasoningEffortNone,
	}
}

func testOpenAIClient(transport openAITransport) openai.Client {
	return openai.NewClient(
		option.WithAPIKey("test-key"),
		option.WithBaseURL("https://openai.test/v1/"),
		option.WithHTTPClient(&http.Client{Transport: transport}),
		option.WithMaxRetries(0),
	)
}

func openAIResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": {"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func TestOpenAIAgentAsk(t *testing.T) {
	for _, formatted := range []bool{false, true} {
		name := "Ask"
		if formatted {
			name = "Askf"
		}
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			called := false
			client := testOpenAIClient(func(r *http.Request) (*http.Response, error) {
				called = true
				require.Equal(t, http.MethodPost, r.Method)
				require.Equal(t, "/v1/responses", r.URL.Path)
				require.Equal(t, "Bearer test-key", r.Header.Get("Authorization"))
				deadline, ok := r.Context().Deadline()
				require.True(t, ok)
				wantDeadline, _ := ctx.Deadline()
				require.Equal(t, wantDeadline, deadline)
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				require.JSONEq(t, `{
					"model":"test-model", "max_output_tokens":512,
					"instructions":"Recommend 3 items for a mage",
					"reasoning":{"effort":"none"},
					"input":[{"role":"user","content":"Return JSON only.\n{\"class\":\"mage\"}"}],
					"text":{"format":{"type":"json_object"}}
				}`, string(body))
				return openAIResponse(http.StatusOK, `{"output":[{"type":"message","role":"assistant","content":[{"type":"output_text","text":" \n{\"answer\":\"fire\"}\n "}]}]}`), nil
			})
			mage := testOpenAIAgent(client, "test-model", 512)
			in := map[string]string{"class": "mage"}
			var out struct {
				Answer string `json:"answer"`
			}
			var err error
			if formatted {
				err = mage.Askf(ctx, in, &out, "Recommend %d items for a %s", 3, "mage")
			} else {
				err = mage.Ask(ctx, in, &out, "Recommend 3 items for a mage")
			}
			require.NoError(t, err)
			require.True(t, called)
			require.Equal(t, "fire", out.Answer)
		})
	}
}

func TestOpenAIAgentAskErrors(t *testing.T) {
	for _, scenario := range []struct {
		name      string
		status    int
		body      string
		wantError string
	}{
		{"api error", http.StatusBadRequest, `{"error":{"message":"invalid model","type":"invalid_request_error"}}`, "openai response failed"},
		{"missing text", http.StatusOK, `{"output":[]}`, "openai response has no output text"},
		{"blank text", http.StatusOK, `{"output":[{"type":"message","content":[{"type":"output_text","text":" \n "}]}]}`, "openai response has no output text"},
		{"invalid json", http.StatusOK, `{"output":[{"type":"message","content":[{"type":"output_text","text":"not json"}]}]}`, "decoding openai response failed"},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			client := testOpenAIClient(func(r *http.Request) (*http.Response, error) {
				var body map[string]any
				require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
				require.NotContains(t, body, "max_output_tokens")
				return openAIResponse(scenario.status, scenario.body), nil
			})
			var out map[string]any
			err := testOpenAIAgent(client, "test-model", 0).Ask(context.Background(), nil, &out, "Answer")
			require.ErrorContains(t, err, scenario.wantError)
			if strings.Contains(scenario.wantError, "no output text") {
				require.JSONEq(t, scenario.body, errors.Tag(err, "raw").(string))
			}
		})
	}
}

func TestOpenAIAgentAskInvalidInput(t *testing.T) {
	client := testOpenAIClient(func(r *http.Request) (*http.Response, error) {
		t.Fatal("invalid input must not make an API request")
		return nil, nil
	})
	var out any
	err := testOpenAIAgent(client, "test-model", 512).Ask(context.Background(), make(chan int), &out, "Answer")
	require.ErrorContains(t, err, "encoding openai message failed")
	var encodingError *json.UnsupportedTypeError
	require.ErrorAs(t, err, &encodingError)
}

func TestOpenAIAgentAskIncomplete(t *testing.T) {
	for _, reason := range []string{"max_output_tokens", "content_filter"} {
		t.Run(reason, func(t *testing.T) {
			client := testOpenAIClient(func(r *http.Request) (*http.Response, error) {
				// Even valid JSON must be rejected when the response is incomplete.
				return openAIResponse(http.StatusOK, `{"status":"incomplete","incomplete_details":{"reason":"`+reason+`"},"output":[{"type":"message","content":[{"type":"output_text","text":"{\"title\":\"Partial article\"}"}]}]}`), nil
			})
			out := map[string]string{"title": "Original article"}
			err := testOpenAIAgent(client, "test-model", 32768).Ask(context.Background(), nil, &out, "Write an article")
			require.ErrorContains(t, err, "openai response is incomplete")
			require.Equal(t, reason, errors.Tag(err, "reason"))
			require.Equal(t, map[string]string{"title": "Original article"}, out)
		})
	}
}

func TestOpenAIAgentAskCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	client := testOpenAIClient(func(r *http.Request) (*http.Response, error) {
		require.ErrorIs(t, r.Context().Err(), context.Canceled)
		return nil, r.Context().Err()
	})
	var out any
	err := testOpenAIAgent(client, "test-model", 512).Ask(ctx, nil, &out, "Answer")
	require.ErrorIs(t, err, context.Canceled)
}

func TestOpenAIAgentAskDefaultReasoning(t *testing.T) {
	called := false
	client := testOpenAIClient(func(r *http.Request) (*http.Response, error) {
		called = true
		var body map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		require.NotContains(t, body, "reasoning")
		return openAIResponse(http.StatusOK, `{"output":[{"type":"message","content":[{"type":"output_text","text":"{}"}]}]}`), nil
	})
	agent := &openAIAgent{Client: client, Model: "test-model"}
	var out map[string]any
	require.NoError(t, agent.Ask(context.Background(), nil, &out, "Answer"))
	require.True(t, called)
}
