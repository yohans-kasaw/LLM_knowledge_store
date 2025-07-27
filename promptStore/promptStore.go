package promptStore

import (
	"github.com/tmc/langchaingo/prompts"
)

var BusinessReviewPrompt = prompts.NewPromptTemplate(
	`You are an expert business document reviewer:
		1. **Clarity**: Assess readability, conciseness, and audience alignment
		2. **Structure**: Evaluate organization, flow, and logical coherence
		3. **Content**: Check evidence strength, relevance, and actionable insights
		4. **Tone**: Ensure professionalism, consistency, and appropriateness
		Document to review:
		File Name is {{.file_name}}

		{{.doc}}
		Provide specific, actionable feedback. For each issue:
		- Explain WHY it's a problem
		- Suggest HOW to fix it with examples
		- Rate severity: Critical, Warning, Suggestion
		Focus on the most impactful issues first.`,
	[]string{"doc", "file_name"})

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
