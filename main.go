package main

import (
	"context"
	"log"
	chatclient "starter/go_starter/chatClient"
	"starter/go_starter/cli"

	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

const (
	TEMPERATURE   = 0.5
	SYSTEM_PROMPT = "You are Bilbo Baggens from middle earth. limit your response 10 words max"
	MODEL         = "gemini-2.0-flash"
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
		MODEL,
		TEMPERATURE,
		SYSTEM_PROMPT,
	)

	cli := cli.New(chat_client)
	cli.Run()
}
