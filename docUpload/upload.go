package docUpload

import (
	"log"
	"os"
	"starter/go_starter/chatClient"
	"starter/go_starter/knowledge"
	"starter/go_starter/promptStore"

	"github.com/tmc/langchaingo/prompts"
)

type Uploader struct {
	chat      *chatclient.ChatClient
	prompt    *prompts.PromptTemplate
	knowledge *knowledge.Knowledge
}

func New(c *chatclient.ChatClient, kg *knowledge.Knowledge) *Uploader {
	return &Uploader{
		chat:      c,
		prompt:    &promptStore.BusinessReviewPrompt,
		knowledge: kg,
	}
}

func (r *Uploader) UploadAndReviewDoc(file_name string) string {
	// open the file here and send it as a doc
	content, err := os.ReadFile(file_name)
	if err != nil {
		log.Printf("err opening file please try again %v\n", err)
		return ""
	}

	msg, err := r.prompt.Format(map[string]any{
		"doc":       string(content),
		"file_name": file_name,
	})

	r.knowledge.AddInputToKnowledge(string(content))

	if err != nil {
		log.Fatal(err)
	}
	res := r.chat.SendMessage(msg)
	return res
}
