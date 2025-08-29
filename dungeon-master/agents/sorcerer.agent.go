package agents

import (
	"context"
	"dungeon-master/helpers"
	"fmt"
	"sync"

	"github.com/micro-agent/micro-agent-go/agent/mu"
	"github.com/openai/openai-go/v2"
)

// TODO:
// - use another model (reasoning model), try with Lucy
// - add rag data

var (
	sorcererAgentInstance mu.Agent
	sorcererAgentOnce     sync.Once
)

// GetSorcererAgent returns the singleton instance of the sorcerer agent
func GetSorcererAgent(ctx context.Context, client openai.Client) mu.Agent {
	sorcererAgentOnce.Do(func() {
		sorcererAgentInstance = createSorcererAgent(ctx, client)
	})
	return sorcererAgentInstance
}

// Huey, Dewey, and Louie
func createSorcererAgent(ctx context.Context, client openai.Client) mu.Agent {

	name := helpers.GetEnvOrDefault("SORCERER_NAME", "Dewey")
	model := helpers.GetEnvOrDefault("SORCERER_MODEL", "ai/qwen2.5:1.5B-F16")
	temperature := helpers.StringToFloat(helpers.GetEnvOrDefault("SORCERER_MODEL_TEMPERATURE", "0.0"))

	// ---------------------------------------------------------
	// System Instructions
	// ---------------------------------------------------------
	var systemInstructions openai.ChatCompletionMessageParamUnion

	systemInstructionsContentPath := helpers.GetEnvOrDefault("SORCERER_SYSTEM_INSTRUCTIONS_PATH", "")
	if systemInstructionsContentPath == "" {
		fmt.Println("🔶 No SORCERER_SYSTEM_INSTRUCTIONS_PATH provided, using default instructions.")
		systemInstructions = openai.SystemMessage("You are a wise and powerful sorcerer in a fantasy world.")
	}

	// Read the content of the file at systemInstructionsContentPath
	systemInstructionsContent, err := helpers.ReadTextFile(systemInstructionsContentPath)

	if err != nil {
		fmt.Println("🔶 Error reading the file, using default instructions:", err)
		systemInstructions = openai.SystemMessage("You are a wise and powerful sorcerer in a fantasy world.")
	} else {
		systemInstructions = openai.SystemMessage(systemInstructionsContent)
	}

	// ---------------------------------------------------------
	// Context Instructions
	// ---------------------------------------------------------
	var contextInstructions openai.ChatCompletionMessageParamUnion

	contextInstructionsContentPath := helpers.GetEnvOrDefault("SORCERER_CONTEXT_PATH", "")
	if contextInstructionsContentPath == "" {
		fmt.Println("🔶 No SORCERER_CONTEXT_PATH provided, using default instructions.")
		contextInstructions = openai.SystemMessage("You are in a fantasy world.")
	}

	// Read the content of the file at contextInstructionsContentPath
	contextInstructionsContent, err := helpers.ReadTextFile(contextInstructionsContentPath)
	if err != nil {
		fmt.Println("🔶 Error reading the file, using default instructions:", err)
		contextInstructions = openai.SystemMessage("You are in a fantasy world.")
	} else {
		contextInstructions = openai.SystemMessage(contextInstructionsContent)
	}

	chatAgent, err := mu.NewAgent(ctx, name,
		mu.WithClient(client),
		mu.WithParams(openai.ChatCompletionNewParams{
			Model:       model,
			Temperature: openai.Opt(temperature),
			Messages: []openai.ChatCompletionMessageParamUnion{
				systemInstructions,
				contextInstructions,
			},
		}),
	)
	if err != nil {
		fmt.Println("🔶 Error creating sorcerer agent, creating ghost agent instead:", err)
		return NewGhostAgent("[Ghost] " + name)
	}

	return chatAgent

}
