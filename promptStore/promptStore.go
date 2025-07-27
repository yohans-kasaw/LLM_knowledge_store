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
