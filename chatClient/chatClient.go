package chatclient

import (
	"context"
	"log"

	"google.golang.org/genai"
)

type ChatClient struct {
	Chat *genai.Chat
	Ctx  context.Context
}

func New(
	ctx context.Context,
	gClient *genai.Client,
	model string,
	temperature float32,
	systemPrompt string,
) *ChatClient {

	config := &genai.GenerateContentConfig{
		Temperature: genai.Ptr(temperature),
	}
	var history []*genai.Content
	history = append(history, &genai.Content{
		Role: "user",
		Parts: []*genai.Part{{
			Text: systemPrompt,
		}},
	})

	chat, err := gClient.Chats.Create(ctx, model, config, history)

	if err != nil {
		log.Fatal(err)
	}

	return &ChatClient{Chat: chat, Ctx: ctx}
}

func (c *ChatClient) SendMessage(msg string) string {
	res, err := c.Chat.SendMessage(c.Ctx, genai.Part{Text: msg})

	if err != nil {
		log.Fatal(err)
	}

	return res.Text()
}
