package agents

import (
	"context"
	"fmt"
	"sync"

	"github.com/micro-agent/micro-agent-go/agent/helpers"
	"github.com/micro-agent/micro-agent-go/agent/mu"

	"github.com/openai/openai-go/v2"
)

var (
	merchantAgentInstance mu.Agent
	merchantAgentOnce     sync.Once
)

// GetMerchantAgent returns the singleton instance of the merchant agent
func GetMerchantAgent(ctx context.Context, client openai.Client) mu.Agent {
	merchantAgentOnce.Do(func() {
		merchantAgentInstance = createMerchantAgent(ctx, client)
	})
	return merchantAgentInstance
}

func createMerchantAgent(ctx context.Context, client openai.Client) mu.Agent {

	name := helpers.GetEnvOrDefault("MERCHANT_NAME", "Thorin")
	model := helpers.GetEnvOrDefault("MERCHANT_MODEL", "ai/qwen2.5:1.5B-F16")
	temperature := helpers.StringToFloat(helpers.GetEnvOrDefault("MERCHANT_MODEL_TEMPERATURE", "0.0"))

	// [RAG]  Initialize the vector store for the agent
	errEmbedding := GenerateEmbeddings(ctx, &client, name, helpers.GetEnvOrDefault("MERCHANT_CONTEXT_PATH", ""))
	if errEmbedding != nil {
		fmt.Println("🔶 Error generating embeddings for merchant agent:", errEmbedding)
	}

	// ---------------------------------------------------------
	// System Instructions
	// ---------------------------------------------------------
	var systemInstructions openai.ChatCompletionMessageParamUnion

	systemInstructionsContentPath := helpers.GetEnvOrDefault("MERCHANT_SYSTEM_INSTRUCTIONS_PATH", "")
	if systemInstructionsContentPath == "" {
		fmt.Println("🔶 No MERCHANT_SYSTEM_INSTRUCTIONS_PATH provided, using default instructions.")
	}

	// Read the content of the file at systemInstructionsContentPath
	systemInstructionsContent, err := helpers.ReadTextFile(systemInstructionsContentPath)

	if err != nil {
		fmt.Println("🔶 Error reading the file, using default instructions:", err)
		systemInstructions = openai.SystemMessage("You are a shrewd merchant in a fantasy world.")
	} else {
		systemInstructions = openai.SystemMessage(systemInstructionsContent)
	}

	chatAgent, err := mu.NewAgent(ctx, name,
		mu.WithClient(client),
		mu.WithParams(openai.ChatCompletionNewParams{
			Model:       model,
			Temperature: openai.Opt(temperature),
			Messages: []openai.ChatCompletionMessageParamUnion{
				systemInstructions,
			},
		}),
	)
	if err != nil {
		fmt.Println("🔶 Error creating merchant agent, creating ghost agent instead:", err)
		return NewGhostAgent("[Ghost] " + name)
	}

	return chatAgent
}