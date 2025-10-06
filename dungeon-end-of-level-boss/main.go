package main

import (
	"context"
	"end-of-level-boss/agents"
	"fmt"
	"strings"

	"github.com/micro-agent/micro-agent-go/agent/experimental/a2a"
	"github.com/micro-agent/micro-agent-go/agent/helpers"
	"github.com/micro-agent/micro-agent-go/agent/msg"
	"github.com/micro-agent/micro-agent-go/agent/ui"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

func main() {

	ctx := context.Background()

	llmURL := helpers.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")
	mcpHost := helpers.GetEnvOrDefault("MCP_HOST", "http://localhost:9011/mcp")

	fmt.Println("🌍 LLM URL:", llmURL)
	fmt.Println("🌍 MCP Host:", mcpHost)

	similaritySearchLimit := helpers.StringToFloat(helpers.GetEnvOrDefault("SIMILARITY_LIMIT", "0.5"))
	similaritySearchMaxResults := helpers.StringToInt(helpers.GetEnvOrDefault("SIMILARITY_MAX_RESULTS", "2"))
	httpPort := helpers.GetEnvOrDefault("BOSS_REMOTE_AGENT_HTTP_PORT", "8080")

	client := openai.NewClient(
		option.WithBaseURL(llmURL),
		option.WithAPIKey(""),
	)
	// [End Of Level Boss Agent]
	bossAgent := agents.GetBossAgent(ctx, client)
	// 👀 `GetBossAgent` 
	// [Agent Card] (for A2A server registration)
	agentCard := a2a.AgentCard{
		Name:        bossAgent.GetName(),
		Description: bossAgent.GetDescription(),
		URL:         "http://localhost:" + httpPort,
		Version:     "1.0.0",
		// Define skills of the remote agent
		Skills: []map[string]any{
			{
				"id":          "ask_for_something",
				"name":        "Ask for something",
				"description": bossAgent.GetName() + " is using a small language model to answer questions",
			},
		},
	}
	// Call remotely the stream completion method of the agent
	// [Agent Stream Callback] (for /stream endpoint)
	agentStreamCallback := func(taskRequest a2a.TaskRequest, streamFunc func(content string) error) error {

		fmt.Printf("🟢 Processing streaming task request: %s\n", taskRequest.ID)
		// Extract user message
		userMessage := taskRequest.Params.Message.Parts[0].Text
		fmt.Printf("🔵 UserMessage: %s\n", userMessage)
		fmt.Printf("🟡 TaskRequest Metadata: %v\n", taskRequest.Params.MetaData)

		var userPrompt string

		switch taskRequest.Params.MetaData["skill"] {
		case "ask_for_something":
			userPrompt = userMessage

		default:
			userPrompt = "Be nice, and explain that " + fmt.Sprintf("%v", taskRequest.Params.MetaData["skill"]) + " is not a valid task ID."
		}

		ui.Println(ui.Green, "<", bossAgent.GetName(), "speaking...>")

		// ---------------------------------------------------------
		// [RAG] SIMILARITY SEARCH:
		// ---------------------------------------------------------
		bossAgentMessages, err := GeneratePromptMessagesWithSimilarities(ctx, &client, bossAgent.GetName(), userPrompt, similaritySearchLimit, similaritySearchMaxResults)

		if err != nil {
			ui.Println(ui.Red, "Error:", err)
		}

		// Generate a [Streaming Chat Completion]
		_, err = bossAgent.RunStream(
			bossAgentMessages,
			func(content string) error {
				if content != "" {
					fmt.Print(content)         // Print to console for debugging
					return streamFunc(content) // Stream to client
				}
				return nil // Continue streaming
			})

		fmt.Println() // Ensure the output ends with a newline
		if err != nil {
			fmt.Printf("❌ Error during streaming chat completion: %v\n", err)
			return err
		}

		// DEBUG: display the messages history
		if strings.HasPrefix(userPrompt, "/debug") {
			msg.DisplayHistory(bossAgent)
		}

		return nil
	}

	a2aServer := a2a.NewA2AServerWithStreaming(helpers.StringToInt(httpPort), agentCard, agentStreamCallback)
	fmt.Println("🚀 Starting A2A server with streaming support on port:", httpPort)
	if err := a2aServer.Start(); err != nil {
		fmt.Printf("❌ Failed to start A2A server: %v\n", err)
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
