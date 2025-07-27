package reviewer

import (
	"fmt"
	"log"
	"os"
	"starter/go_starter/chatClient"
	"starter/go_starter/promptStore"

	"github.com/tmc/langchaingo/prompts"
)

type Reviewer struct {
	chat   *chatclient.ChatClient
	prompt *prompts.PromptTemplate
}

func New(c *chatclient.ChatClient) *Reviewer {
	return &Reviewer{
		chat:   c,
		prompt: &promptStore.BusinessReviewPrompt,
	}
}

func (r *Reviewer) Review(doc *string) string {
	msg, err := r.prompt.Format(map[string]any{
		"doc": *doc,
	})
	if err != nil {
		log.Fatal(err)
	}
	res := r.chat.SendMessage(msg)
	return res
}

func (r *Reviewer) ReviewDoc(file_name string) string {
	// open the file here and send it as a doc
	content, err := os.ReadFile(file_name)
	if err != nil{
		fmt.Printf("err opening file please try again %v\n", err)
		return ""
	}

	msg, err := r.prompt.Format(map[string]any{
		"doc": content,
		"file_name": file_name,
	})
	if err != nil {
		log.Fatal(err)
	}
	res := r.chat.SendMessage(msg)
	return res
}
