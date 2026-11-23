package sanctum

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/maxence-charriere/go-app/v11/pkg/errors"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

type openAIAgent struct {
	Client          openai.Client
	Model           string
	MaxToken        int
	ReasoningEffort openai.ReasoningEffort
}

var _ Magister = (*openAIAgent)(nil)

func (a *openAIAgent) Ask(ctx context.Context, in any, out any, prompt string) error {
	msg, err := json.Marshal(in)
	if err != nil {
		return errors.New("encoding openai message failed").Wrap(err)
	}

	params := responses.ResponseNewParams{
		Model:        a.Model,
		Instructions: openai.String(prompt),
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: responses.ResponseInputParam{
				responses.ResponseInputItemUnionParam{
					OfMessage: &responses.EasyInputMessageParam{
						Role: "user",
						Content: responses.EasyInputMessageContentUnionParam{
							OfString: openai.String("Return JSON only.\n" + string(msg)),
						},
					},
				},
			},
		},
		Text: responses.ResponseTextConfigParam{
			Format: responses.ResponseFormatTextConfigUnionParam{
				OfJSONObject: &responses.ResponseFormatJSONObjectParam{},
			},
		},
	}
	if a.ReasoningEffort != "" {
		params.Reasoning = openai.ReasoningParam{Effort: a.ReasoningEffort}
	}
	if a.MaxToken > 0 {
		params.MaxOutputTokens = openai.Int(int64(a.MaxToken))
	}

	resp, err := a.Client.Responses.New(ctx, params)
	if err != nil {
		return errors.New("openai response failed").Wrap(err)
	}
	if resp.Status == responses.ResponseStatusIncomplete {
		return errors.New("openai response is incomplete").
			WithTag("reason", resp.IncompleteDetails.Reason)
	}

	text := strings.TrimSpace(resp.OutputText())
	if text == "" {
		return errors.New("openai response has no output text").WithTag("raw", resp.RawJSON())
	}
	if err := json.Unmarshal([]byte(text), out); err != nil {
		return errors.New("decoding openai response failed").Wrap(err)
	}
	return nil
}

func (a *openAIAgent) Askf(ctx context.Context, in any, out any, promptFormat string, v ...any) error {
	return a.Ask(ctx, in, out, fmt.Sprintf(promptFormat, v...))
}
