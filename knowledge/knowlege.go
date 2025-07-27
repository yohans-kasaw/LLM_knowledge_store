package knowledge

import (
	"context"
	chatclient "starter/go_starter/chatClient"

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
