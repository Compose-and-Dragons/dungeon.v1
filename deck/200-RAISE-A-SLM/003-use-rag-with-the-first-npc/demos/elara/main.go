package main

import (
	"context"
	"elara/agents"
	"fmt"
	"strings"

	"github.com/micro-agent/micro-agent-go/agent/helpers"
	"github.com/micro-agent/micro-agent-go/agent/msg"
	"github.com/micro-agent/micro-agent-go/agent/ui"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

func main() {

	ctx := context.Background()
	llmURL := helpers.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")

	similaritySearchLimit := helpers.StringToFloat(helpers.GetEnvOrDefault("SIMILARITY_LIMIT", "0.5"))
	similaritySearchMaxResults := helpers.StringToInt(helpers.GetEnvOrDefault("SIMILARITY_MAX_RESULTS", "2"))

	client := openai.NewClient(
		option.WithBaseURL(llmURL),
		option.WithAPIKey(""),
	)

	// ---------------------------------------------------------
	// AGENT: This is the Sorcerer agent
	// ---------------------------------------------------------
	sorcererAgent := agents.GetSorcererAgent(ctx, client)

	for {
		promptText := "🤖 (/bye to exit) [" + sorcererAgent.GetName() + "]>"
		// PROMPT:
		content, _ := ui.SimplePrompt(promptText, "Type your command here...")

		// USER MESSAGE: content.Input

		// ---------------------------------------------------------
		// Bye [COMMAND]
		// ---------------------------------------------------------
		if strings.HasPrefix(content.Input, "/bye") {
			fmt.Println("👋 Goodbye! Thanks for playing!")
			break
		}

		// DEBUG:
		if strings.HasPrefix(content.Input, "/memory") {
			msg.DisplayHistory(sorcererAgent)
			continue
		}

		// ---------------------------------------------------------
		// TALK TO: AGENT:: **SORCERER** + [RAG]
		// ---------------------------------------------------------
		ui.Println(ui.Purple, "<", sorcererAgent.GetName(), "speaking...>")

		// ---------------------------------------------------------
		// [RAG] SIMILARITY SEARCH:
		// ---------------------------------------------------------
		sorcererAgentMessages, err := GeneratePromptMessagesWithSimilarities(ctx, &client, sorcererAgent.GetName(), content.Input, similaritySearchLimit, similaritySearchMaxResults)

		if err != nil {
			ui.Println(ui.Red, "Error:", err)
		}

		// NOTE: RunStreams adds the messages to the agent's memory
		_, err = sorcererAgent.RunStream(sorcererAgentMessages, func(content string) error {
			fmt.Print(content)
			return nil
		})

		if err != nil {
			ui.Println(ui.Red, "Error:", err)
		}

		fmt.Println()
		fmt.Println()

	}

}

func GeneratePromptMessagesWithSimilarities(ctx context.Context, client *openai.Client, agentName, input string, similarityLimit float64, maxResults int) ([]openai.ChatCompletionMessageParamUnion, error) {
	fmt.Printf("🔍 Searching for similar chunks to '%s'\n", input)

	similarities, err := agents.SearchSimilarities(ctx, client, agentName, input, similarityLimit, maxResults)
	if err != nil {
		fmt.Println("🔴 Error searching for similarities:", err)
		return []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(input),
		}, err
	}

	if len(similarities) > 0 {
		similaritiesMessage := "Here is some context that might be useful:\n"
		for _, similarity := range similarities {
			similaritiesMessage += fmt.Sprintf("- %s\n", similarity.Prompt)
		}
		return []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(similaritiesMessage),
			openai.UserMessage(input),
		}, nil
	} else {
		fmt.Println("📝 No similarities found.")
		return []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(input),
		}, nil
	}
}
