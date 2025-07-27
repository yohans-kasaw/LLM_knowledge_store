package main

import (
	"context"
	"log"
	chatclient "starter/go_starter/chatClient"
	"starter/go_starter/cli"
	"starter/go_starter/docUpload"
	"starter/go_starter/knowledge"
	"starter/go_starter/promptStore"

	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

func main() {
	godotenv.Load()
	ctx := context.Background()

	ANNA_SYSTEM_PROMPT, _ := promptStore.AnnaSystemPrompt.Prompt.Format(map[string]any{})

	gClient, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	knowledge := knowledge.New(ctx, gClient, chatclient.New(
		ctx,
		gClient,
		"gemini-2.0-flash",
		0.5,
		ANNA_SYSTEM_PROMPT,
	))

	uploader := docUpload.New(
		chatclient.New(
			ctx,
			gClient,
			"gemini-2.0-flash",
			0.5,
			ANNA_SYSTEM_PROMPT,
		),
		knowledge,
	)
	cli := cli.New(
		chatclient.New(
			ctx,
			gClient,
			"gemini-2.0-flash",
			0.5,
			ANNA_SYSTEM_PROMPT,
		),
		uploader,
		knowledge,
	)
	cli.Run()
}
