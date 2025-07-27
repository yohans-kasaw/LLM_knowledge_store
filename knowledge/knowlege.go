package knowledge

import (
	"context"
	"fmt"
	chatclient "starter/go_starter/chatClient"
	"starter/go_starter/promptStore"
	"strings"

	"google.golang.org/genai"
)

type Knowledge struct {
	Store     *Store
	Extractor *Extractor
}

func New(ctx context.Context, gClient *genai.Client, chatClient *chatclient.ChatClient) *Knowledge {

	store := NewStore(ctx, gClient)
	extractor := NewExtractor(chatClient)

	return &Knowledge{
		Store:     store,
		Extractor: extractor,
	}
}

func (k *Knowledge) AddInputToKnowledge(input string) {
	chunks := k.Extractor.ExtractFromUserInput(input)
	if chunks != nil {
		k.Store.AddKnowledge(*chunks)
	}
}

func (k *Knowledge) EmbbedAdditonalKnowledge(input string) string {
	topK := uint64(4)
	retrievedKnowledge, err := k.Store.RetrieveKnowledge(input, &topK)
	if err != nil {
		return input
	}

	if len(retrievedKnowledge) == 0 {
		return input
	}

	var knowledgePointsBuilder strings.Builder
	for i, sentence := range retrievedKnowledge {
		knowledgePointsBuilder.WriteString(fmt.Sprintf("Knowledge Point %d: %s\n", i+1, sentence))
	}

	formattedPrompt, err := promptStore.KnowledgeEmbeddingPrompt.Format(
		map[string]any{
			"retrieved_knowledge": knowledgePointsBuilder.String(),
			"user_query":          input,
		},
	)
	if err != nil {
		return input
	}

	return formattedPrompt
}
