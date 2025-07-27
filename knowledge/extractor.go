package knowledge

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/prompts"
	"google.golang.org/genai"

	chatclient "starter/go_starter/chatClient"
	"starter/go_starter/promptStore"
)

const (
	MAX_CHUNK_WORDS    = 20
	MIN_SENTENCE_WORDS = 3
)

type Extractor struct {
	chat   *chatclient.ChatClient
	prompt *prompts.PromptTemplate
}

func NewExtractor(chatClient *chatclient.ChatClient) *Extractor {

	return &Extractor{
		chat:   chatClient,
		prompt: &promptStore.KnowledgeExtractionPrompt,
	}
}

func (e *Extractor) ExtractFromUserInput(input string) *[]string {
	filteredInput := heuristicFilter(input)
	if filteredInput == "" {
		log.Println("User input discarded by heuristic filter (too short or empty).")
		return nil
	}

	chunks, err := e.processWithGenAI(filteredInput)
	if err != nil {
		log.Printf("failed to process user input with GenAI: %e\n", err)
		return nil
	}

	if len(chunks) > 0 {
		log.Printf("Successfully extracted and stored %d knowledge chunks from user input.", len(chunks))
		return &chunks
	} else {
		log.Println("No useful knowledge chunks extracted from user input.")
	}

	return nil
}

func heuristicFilter(text string) string {
	sentences := splitIntoSentences(text)
	fmt.Println("extracted sentences", sentences)
	var filteredSentences []string
	for _, s := range sentences {
		wordCount := len(strings.Fields(s))
		if wordCount >= MIN_SENTENCE_WORDS {
			filteredSentences = append(filteredSentences, s)
		}
	}
	return strings.Join(filteredSentences, " ")
}

func (e *Extractor) processWithGenAI(text string) ([]string, error) {

	msg, err := e.prompt.Format(map[string]any{
		"text":            text,
		"max_chunk_words": MAX_CHUNK_WORDS,
	})

	if err != nil {
		log.Fatal(err)
	}

	genResp := e.chat.SendMessage(msg)

	res := extractChunks(genResp)
	return res, nil
}

func extractChunks(response string) []string {
	type DataWrapper struct {
		Data []string `json:"data"`
	}
	var wrapper DataWrapper

	cleanResponse := response
	if strings.HasPrefix(strings.TrimSpace(response), "```") {
		startIndex := strings.Index(response, "\n")
		if startIndex != -1 {
			cleanResponse = response[startIndex+1:]
		}

		endIndex := strings.LastIndex(cleanResponse, "```")
		if endIndex != -1 {
			cleanResponse = cleanResponse[:endIndex]
		}
	}

	cleanResponse = strings.TrimSpace(cleanResponse)

	err := json.Unmarshal([]byte(cleanResponse), &wrapper)
	if err != nil {
		fmt.Printf("Error unmarshaling JSON: %v\n", err)
		fmt.Printf("Attempted to unmarshal: %s\n", cleanResponse)
		return []string{}
	}
	return wrapper.Data
}

func splitIntoSentences(text string) []string {
	re := regexp.MustCompile(`([.!?]+)`)
	parts := re.Split(text, -1)
	punctuationMarks := re.FindAllString(text, -1)
	var result []string
	puncIndex := 0
	for _, part := range parts {
		sentence := strings.TrimSpace(part)
		if sentence != "" {
			words := strings.Fields(sentence)
			for i := 0; i < len(words); i += MAX_CHUNK_WORDS {
				end := i + MAX_CHUNK_WORDS
				if end > len(words) {
					end = len(words)
				}
				result = append(result, strings.Join(words[i:end], " "))
			}
		}

		if puncIndex < len(punctuationMarks) {
			punc := strings.TrimSpace(punctuationMarks[puncIndex])
			if punc != "" { // Ensure it's not just an empty string from trimming
				result = append(result, punc)
			}
			puncIndex++
		}
	}
	return result
}

func TestExtractor() {

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
		"You are Bilbo Baggens from middle earth. limit your response 10 words max",
	)

	e := NewExtractor(chat_client)
	res := e.ExtractFromUserInput("My name is Yohans.")
	fmt.Println(res)
}
