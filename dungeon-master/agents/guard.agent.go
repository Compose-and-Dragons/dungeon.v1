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
	guardAgentInstance mu.Agent
	guardAgentOnce     sync.Once
)

// GetGuardAgent returns the singleton instance of the guard agent
func GetGuardAgent(ctx context.Context, client openai.Client) mu.Agent {
	guardAgentOnce.Do(func() {
		guardAgentInstance = createGuardAgent(ctx, client)
	})
	return guardAgentInstance
}

// Huey, Dewey, and Louie
func createGuardAgent(ctx context.Context, client openai.Client) mu.Agent {

	name := helpers.GetEnvOrDefault("GUARD_NAME", "Huey")
	model := helpers.GetEnvOrDefault("GUARD_MODEL", "ai/qwen2.5:1.5B-F16")
	temperature := helpers.StringToFloat(helpers.GetEnvOrDefault("GUARD_MODEL_TEMPERATURE", "0.0"))

	// [RAG]  Initialize the vector store for the agent
	errEmbedding := GenerateEmbeddings(ctx, &client, name, helpers.GetEnvOrDefault("GUARD_CONTEXT_PATH", ""))
	if errEmbedding != nil {
		fmt.Println("🔶 Error generating embeddings for guard agent:", errEmbedding)
	}

	// ---------------------------------------------------------
	// System Instructions
	// ---------------------------------------------------------
	var systemInstructions openai.ChatCompletionMessageParamUnion

	systemInstructionsContentPath := helpers.GetEnvOrDefault("GUARD_SYSTEM_INSTRUCTIONS_PATH", "")
	if systemInstructionsContentPath == "" {
		fmt.Println("🔶 No GUARD_SYSTEM_INSTRUCTIONS_PATH provided, using default instructions.")
		//systemInstructions = openai.SystemMessage("You are an elf guard in a fantasy world.")
	}

	// Read the content of the file at systemInstructionsContentPath
	systemInstructionsContent, err := helpers.ReadTextFile(systemInstructionsContentPath)

	if err != nil {
		fmt.Println("🔶 Error reading the file, using default instructions:", err)
		systemInstructions = openai.SystemMessage("You are an elf guard in a fantasy world.")
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
		fmt.Println("🔶 Error creating sorcerer agent, creating ghost agent instead:", err)
		return NewGhostAgent("[Ghost] " + name)
	}
	return chatAgent
}
