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
	healerAgentInstance mu.Agent
	healerAgentOnce     sync.Once
)

// GetHealerAgent returns the singleton instance of the healer agent
func GetHealerAgent(ctx context.Context, client openai.Client) mu.Agent {
	healerAgentOnce.Do(func() {
		healerAgentInstance = createHealerAgent(ctx, client)
	})
	return healerAgentInstance
}

func createHealerAgent(ctx context.Context, client openai.Client) mu.Agent {

	name := helpers.GetEnvOrDefault("HEALER_NAME", "Seraphina")
	model := helpers.GetEnvOrDefault("HEALER_MODEL", "ai/qwen2.5:1.5B-F16")
	temperature := helpers.StringToFloat(helpers.GetEnvOrDefault("HEALER_MODEL_TEMPERATURE", "0.0"))

	// [RAG]  Initialize the vector store for the agent
	errEmbedding := GenerateEmbeddings(ctx, &client, name, helpers.GetEnvOrDefault("HEALER_CONTEXT_PATH", ""))
	if errEmbedding != nil {
		fmt.Println("🔶 Error generating embeddings for healer agent:", errEmbedding)
	}

	// ---------------------------------------------------------
	// System Instructions
	// ---------------------------------------------------------
	var systemInstructions openai.ChatCompletionMessageParamUnion

	systemInstructionsContentPath := helpers.GetEnvOrDefault("HEALER_SYSTEM_INSTRUCTIONS_PATH", "")
	if systemInstructionsContentPath == "" {
		fmt.Println("🔶 No HEALER_SYSTEM_INSTRUCTIONS_PATH provided, using default instructions.")
	}

	// Read the content of the file at systemInstructionsContentPath
	systemInstructionsContent, err := helpers.ReadTextFile(systemInstructionsContentPath)

	if err != nil {
		fmt.Println("🔶 Error reading the file, using default instructions:", err)
		systemInstructions = openai.SystemMessage("You are a wise healer in a fantasy world.")
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
		fmt.Println("🔶 Error creating healer agent, creating ghost agent instead:", err)
		return NewGhostAgent("[Ghost] " + name)
	}

	return chatAgent
}