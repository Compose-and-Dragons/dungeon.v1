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
	sorcererAgentInstance mu.Agent
	sorcererAgentOnce     sync.Once
)

// GetBossAgent returns the singleton instance of the sorcerer agent
func GetBossAgent(ctx context.Context, client openai.Client) mu.Agent {
	sorcererAgentOnce.Do(func() {
		sorcererAgentInstance = createBossAgent(ctx, client)
	})
	return sorcererAgentInstance
}

// Huey, Dewey, and Louie
func createBossAgent(ctx context.Context, client openai.Client) mu.Agent {

	name := helpers.GetEnvOrDefault("BOSS_NAME", "Louie")
	description := helpers.GetEnvOrDefault("BOSS_DESCRIPTION", "A wise and powerful Sphinx in a fantasy world.")
	model := helpers.GetEnvOrDefault("BOSS_MODEL", "ai/qwen2.5:1.5B-F16")
	temperature := helpers.StringToFloat(helpers.GetEnvOrDefault("BOSS_MODEL_TEMPERATURE", "0.0"))

	// [RAG]  Initialize the vector store for the agent
	errEmbedding := GenerateEmbeddings(ctx, &client, name, helpers.GetEnvOrDefault("BOSS_CONTEXT_PATH", ""))
	if errEmbedding != nil {
		fmt.Println("🔶 Error generating embeddings for sorcerer agent:", errEmbedding)
	}

	// ---------------------------------------------------------
	// System Instructions
	// ---------------------------------------------------------
	var systemInstructions openai.ChatCompletionMessageParamUnion

	systemInstructionsContentPath := helpers.GetEnvOrDefault("BOSS_SYSTEM_INSTRUCTIONS_PATH", "")
	if systemInstructionsContentPath == "" {
		fmt.Println("🔶 No BOSS_SYSTEM_INSTRUCTIONS_PATH provided, using default instructions.")
	}

	// Read the content of the file at systemInstructionsContentPath
	systemInstructionsContent, err := helpers.ReadTextFile(systemInstructionsContentPath)

	if err != nil {
		fmt.Println("🔶 Error reading the file, using default instructions:", err)
		systemInstructions = openai.SystemMessage("You are a wise and powerful Sphinx in a fantasy world.")
	} else {
		systemInstructions = openai.SystemMessage(systemInstructionsContent)
	}

	chatAgent, err := mu.NewAgentWithDescription(ctx, name, description,
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
		fmt.Println("🔶 Error creating boss agent, creating ghost agent instead:", err)
		return NewGhostAgent("[Ghost] " + name)
	}

	return chatAgent

}
