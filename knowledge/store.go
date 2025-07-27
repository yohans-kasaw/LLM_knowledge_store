package knowledge

import (
	"context"
	"fmt"
	"log"
	"slices"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/qdrant/go-client/qdrant"
	"google.golang.org/genai"
)

const (
	EMBED_MODEL      = "gemini-embedding-001"
	COLLECTION_NAME  = "rag_knowledge_store_test_3"
	VECTOR_DIMENSION = 3072
)

type Store struct {
	qdrantClient *qdrant.Client
	genaiClient  *genai.Client
	ctx          context.Context
}

func NewStore(ctx context.Context, gClient *genai.Client) *Store {
	qClient, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6334,
	})

	if err != nil {
		log.Panicf("failed to create Qdrant client: %w", err)
		return nil
	}

	store := &Store{
		qdrantClient: qClient,
		genaiClient:  gClient,
		ctx:          ctx,
	}

	err = store.ensureCollection()
	if err != nil {
		log.Panicf("failed to ensure Qdrant collection: %w", err)
		return nil
	}

	return store
}

func (s *Store) AddKnowledge(sentences []string) error {
	if len(sentences) == 0 {
		return nil
	}

	var content []*genai.Content
	for _, t := range sentences {
		content = append(
			content,
			genai.NewContentFromText(t, genai.RoleUser),
		)
	}

	result, err := s.genaiClient.Models.EmbedContent(
		s.ctx,
		EMBED_MODEL,
		content,
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed to embed content: %w", err)
	}

	if len(result.Embeddings) != len(sentences) {
		return fmt.Errorf("mismatch between sentences and embeddings count")
	}

	points := make([]*qdrant.PointStruct, len(sentences))
	for i := range sentences {
		uuidV4 := uuid.New()
		points[i] = &qdrant.PointStruct{
			Id:      qdrant.NewIDUUID(uuidV4.String()),
			Vectors: qdrant.NewVectors(result.Embeddings[i].Values...),
			Payload: qdrant.NewValueMap(map[string]any{
				"text": sentences[i],
			}),
		}
	}

	_, err = s.qdrantClient.Upsert(s.ctx, &qdrant.UpsertPoints{
		CollectionName: COLLECTION_NAME,
		Points:         points,
	})

	if err != nil {
		return fmt.Errorf("failed to upsert points to Qdrant: %w", err)
	}
	return nil
}

func (s *Store) RetrieveKnowledge(query string, topK *uint64) ([]string, error) {
	queryContent := []*genai.Content{genai.NewContentFromText(query, genai.RoleUser)}

	queryEmbeddingResult, err := s.genaiClient.Models.EmbedContent(
		s.ctx,
		EMBED_MODEL,
		queryContent,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to embed query: %w", err)
	}

	if len(queryEmbeddingResult.Embeddings) == 0 {
		return nil, fmt.Errorf("no embedding generated for the query")
	}

	searchResult, err := s.qdrantClient.Query(s.ctx, &qdrant.QueryPoints{
		CollectionName: COLLECTION_NAME,
		Query:          qdrant.NewQuery(queryEmbeddingResult.Embeddings[0].Values...),
		WithPayload:    qdrant.NewWithPayload(true),
		Limit:          topK,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search Qdrant: %w", err)
	}

	return extractSentencesFromScoredPoints(searchResult), nil
}

func (s *Store) ensureCollection() error {
	collections, err := s.qdrantClient.ListCollections(s.ctx)
	if err != nil {
		return fmt.Errorf("failed to list Qdrant collections: %w", err)
	}

	if slices.Contains(collections, COLLECTION_NAME) {
		return nil
	}

	err = s.qdrantClient.CreateCollection(s.ctx, &qdrant.CreateCollection{
		CollectionName: COLLECTION_NAME,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     uint64(VECTOR_DIMENSION),
			Distance: qdrant.Distance_Cosine,
		}),
	})
	if err != nil {
		return fmt.Errorf("failed to create collection '%s': %w", COLLECTION_NAME, err)
	}
	return nil
}

func extractSentencesFromScoredPoints(points []*qdrant.ScoredPoint) []string {
	var sentences []string
	for _, point := range points {
		if point.Payload != nil {
			if textVal, ok := point.Payload["text"]; ok {
				if textStr := textVal.GetStringValue(); textStr != "" {
					sentences = append(sentences, textStr)
				}
			}
		}
	}
	return sentences
}

func TestStore() {
	godotenv.Load()

	ctx := context.Background()
	gClient, err := genai.NewClient(ctx, nil)

	store := NewStore(ctx, gClient)
	if err != nil {
		log.Fatalf("Failed to create knowledge store: %v", err)
	}

	sentencesToAdd := []string{
		"The quick brown fox jumps over the lazy dog.",
		"Artificial intelligence is rapidly advancing.",
		"Go is a statically typed, compiled programming language.",
		"Retrieval Augmented Generation combines retrieval with language models.",
		"The sun rises in the east and sets in the west.",
	}
	err = store.AddKnowledge(sentencesToAdd)
	if err != nil {
		log.Fatalf("Failed to add knowledge: %v", err)
	}
	fmt.Println("Knowledge added successfully.")

	// 2. Retrieve Knowledge
	query := "What is RAG?"
	topK := uint64(4)
	retrievedSentences, err := store.RetrieveKnowledge(query, &topK)
	if err != nil {
		log.Fatalf("Failed to retrieve knowledge: %v", err)
	}

	fmt.Printf("\nQuery: %s\n", query)
	fmt.Println("Retrieved Knowledge:")
	for i, sentence := range retrievedSentences {
		fmt.Printf("%d. %s\n", i+1, sentence)
	}

	query2 := "brown fox jumps over what?"

	topK = uint64(1)
	retrievedSentences2, err := store.RetrieveKnowledge(query2, &topK)
	if err != nil {
		log.Fatalf("Failed to retrieve knowledge: %v", err)
	}

	fmt.Printf("\nQuery: %s\n", query2)
	fmt.Println("Retrieved Knowledge:")
	for i, sentence := range retrievedSentences2 {
		fmt.Printf("%d. %s\n", i+1, sentence)
	}
}

func (s *Store) IsUserFirstTime() bool {
	collections, err := s.qdrantClient.ListCollections(s.ctx)

	if err == nil {
		return false 
	}

	if !slices.Contains(collections, COLLECTION_NAME) {
		return false
	}

	return true
}
