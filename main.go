package main

import (
	"context"
	"log"
	chatclient "starter/go_starter/chatClient"
	"starter/go_starter/cli"
	"starter/go_starter/docUpload"
	"starter/go_starter/knowledge"

	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

const (
	SYSTEM_PROMPT = "You are Bilbo Baggens from middle earth. limit your response 10 words max"
)

func main() {
	godotenv.Load()
	ctx := context.Background()

	gClient, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	knowledge := knowledge.New(ctx, gClient, chatclient.New(
		ctx,
		gClient,
		"gemini-2.0-flash",
		0.5,
		SYSTEM_PROMPT,
	))

	uploader := docUpload.New(
		chatclient.New(
			ctx,
			gClient,
			"gemini-2.0-flash",
			0.5,
			SYSTEM_PROMPT,
		),
		knowledge,
	)
	cli := cli.New(
		chatclient.New(
			ctx,
			gClient,
			"gemini-2.0-flash",
			0.5,
			SYSTEM_PROMPT,
		),
		uploader,
		knowledge,
	)
	cli.Run()
}
