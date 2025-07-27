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

	chat_client := chatclient.New(
		ctx,
		gClient,
		"gemini-2.0-flash",
		0.5,
		SYSTEM_PROMPT,
	)

	extractor_chat_client := chatclient.New(
		ctx,
		gClient,
		"gemini-2.0-flash",
		0.5,
		SYSTEM_PROMPT,
	)

	uploader_chat_client := chatclient.New(
		ctx,
		gClient,
		"gemini-2.0-flash",
		0.5,
		SYSTEM_PROMPT,
	)

	store := knowledge.NewStore(ctx, gClient)
	extractor := knowledge.NewExtractor(extractor_chat_client)
	uploader := docUpload.New(uploader_chat_client, store, extractor)
	cli := cli.New(chat_client, uploader, store, extractor)
	cli.Run()
}
