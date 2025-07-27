package main

import (
	"context"
	"github.com/joho/godotenv"
	"google.golang.org/genai"
	"log"
	chatclient "starter/go_starter/chatClient"
	"starter/go_starter/cli"
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

	cli := cli.New(chat_client)
	cli.Run()
}
