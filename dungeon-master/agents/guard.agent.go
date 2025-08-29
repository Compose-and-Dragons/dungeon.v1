package agents

import (
	"context"
	"dungeon-master/helpers"
	"fmt"
	"sync"

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
	model := helpers.GetEnvOrDefault("NON_PLAYER_CHARACTER_MODEL", "ai/qwen2.5:1.5B-F16")
	temperature := helpers.StringToFloat(helpers.GetEnvOrDefault("NON_PLAYER_CHARACTER_TEMPERATURE", "0.0"))

	systemInstructions := openai.SystemMessage(helpers.GetEnvOrDefault(
		"GUARD_SYSTEM_INSTRUCTIONS",
		`You are a guard at the entrance of a medieval castle.`,
	))

	contextInstructions := openai.SystemMessage(helpers.GetEnvOrDefault(
		"GUARD_CONTEXT",
		`TODO`,
	))

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

// IMPORTANT: QUESTION: how to handle the guard's memory?
// I would say keep it simple...
