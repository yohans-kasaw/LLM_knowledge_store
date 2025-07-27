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
	Store     *knowledge.Store
	Extractor *knowledge.Extractor
}

func New(c *chatclient.ChatClient, store *knowledge.Store, extractor *knowledge.Extractor) *Uploader {
	return &Uploader{
		chat:   c,
		prompt: &promptStore.BusinessReviewPrompt,
		Store: store,
		Extractor: extractor,
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

	r.addInputToKnowledge(string(content))

	if err != nil {
		log.Fatal(err)
	}
	res := r.chat.SendMessage(msg)
	return res
}

func (r *Uploader) addInputToKnowledge(input string) {
	chunks := r.Extractor.ExtractFromUserInput(input)
	r.Store.AddKnowledge(*chunks)
}
