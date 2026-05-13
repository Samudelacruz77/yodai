package chat

import (
	"context"

	openai "github.com/sashabaranov/go-openai"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Client struct {
	api   *openai.Client
	model string
}

func NewClient(baseURL, model string) *Client {
	cfg := openai.DefaultConfig("tensorrt_llm")
	cfg.BaseURL = baseURL + "/v1"
	return &Client{
		api:   openai.NewClientWithConfig(cfg),
		model: model,
	}
}

func (c *Client) Chat(ctx context.Context, message string, history []Message, maxTokens int, temperature, topP float32) (string, error) {
	msgs := c.buildMessages(message, history)
	resp, err := c.api.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       c.model,
		Messages:    msgs,
		MaxTokens:   maxTokens,
		Temperature: temperature,
		TopP:        topP,
	})
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", nil
	}
	return resp.Choices[0].Message.Content, nil
}

func (c *Client) StreamChat(ctx context.Context, message string, history []Message, maxTokens int, temperature, topP float32) (<-chan string, <-chan error) {
	tokens := make(chan string)
	errs := make(chan error, 1)

	go func() {
		defer close(tokens)
		defer close(errs)

		msgs := c.buildMessages(message, history)
		stream, err := c.api.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
			Model:       c.model,
			Messages:    msgs,
			MaxTokens:   maxTokens,
			Temperature: temperature,
			TopP:        topP,
			Stream:      true,
		})
		if err != nil {
			errs <- err
			return
		}
		defer stream.Close()

		for {
			resp, err := stream.Recv()
			if err != nil {
				if err.Error() != "EOF" {
					errs <- err
				}
				return
			}
			if len(resp.Choices) > 0 {
				delta := resp.Choices[0].Delta.Content
				if delta != "" {
					select {
					case tokens <- delta:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	return tokens, errs
}

func (c *Client) buildMessages(message string, history []Message) []openai.ChatCompletionMessage {
	msgs := make([]openai.ChatCompletionMessage, 0, len(history)+2)
	msgs = append(msgs, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: YodaSystemPrompt,
	})
	for _, h := range history {
		msgs = append(msgs, openai.ChatCompletionMessage{
			Role:    h.Role,
			Content: h.Content,
		})
	}
	msgs = append(msgs, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: message,
	})
	return msgs
}
