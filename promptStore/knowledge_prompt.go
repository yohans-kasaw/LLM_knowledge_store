package promptStore

import (
	"github.com/tmc/langchaingo/prompts"
)

var KnowledgeExtractionPrompt = prompts.NewPromptTemplate(
	`
You are an expert knowledge extractor. Your task is to extract valuable information from the provided text.

Break down any long information into multiple, shorter, related chunks.

Paraphrase and summarize the extracted information concisely.

Ensure each output chunk is a short sentence, with a maximum of {{.max_chunk_words}} words.


CRITICAL: Output MUST be pure JSON format with a "data" key containing an array of strings.
Each string should be a distinct knowledge chunk.


Example 1:
Input: "The capital of France is Paris. It is known for the Eiffel Tower and its rich history."
Output: {"data": ["Paris is the capital of France.", "The Eiffel Tower is located in Paris.", "Paris boasts a rich history."]}

Example 2:
Input: "Albert Einstein developed the theory of relativity. He won the Nobel Prize in Physics in 1921 for his explanation of the photoelectric effect."
Output: {"data": ["Einstein developed the theory of relativity.", "Einstein won the Nobel Prize in Physics in 1921.", "The Nobel Prize was for explaining the photoelectric effect."]}

Input Text:
{{.text}}
	`,
	[]string{"text", "max_chunk_words"})

var KnowledgeEmbeddingPrompt = prompts.NewPromptTemplate(
	`The following information is retrieved from a knowledge base and may be relevant to your query:

{{.retrieved_knowledge}}

Based on the above knowledge and your original request, please respond to the following query:

User Query: {{.user_query}}

Please note: The provided 'Knowledge Points' are extracted from a database and may not be exhaustive or perfectly align with the current context. Treat them as supplementary information.`,
	[]string{"retrieved_knowledge", "user_query"},
)
