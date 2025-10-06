package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	"github.com/openai/openai-go/v2/shared"
)

func main() {
	ctx := context.Background()

	// Docker Model Runner base URL
	chatURL := os.Getenv("MODEL_RUNNER_BASE_URL")
	model := os.Getenv("MODEL_RUNNER_LLM_TOOLS")

	client := openai.NewClient(
		option.WithBaseURL(chatURL),
		option.WithAPIKey(""),
	)

	// TOOLS: Define the tools available to the model
	// Each tool must have a name, description, and parameters schema

	// TOOL:
	castSpellTool := openai.ChatCompletionFunctionTool(shared.FunctionDefinitionParam{
		Name:        "cast_spell",
		Description: openai.String("Cast a magical spell on the target adventurer"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]interface{}{
				"target": map[string]string{
					"type": "string",
				},
			},
			"required": []string{"target"},
		},
	})

	// TOOL:
	bardInspireTool := openai.ChatCompletionFunctionTool(shared.FunctionDefinitionParam{
		Name:        "bardic_inspiration",
		Description: openai.String("Grant bardic inspiration to boost an ally's courage and abilities"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]interface{}{
				"ally": map[string]string{
					"type": "string",
				},
			},
			"required": []string{"ally"},
		},
	})

	// TOOLS: used by the parameters request
	tools := []openai.ChatCompletionToolUnionParam{
		castSpellTool,
		bardInspireTool,
	}

	// USER MESSAGE:
	userQuestion := openai.UserMessage(`
		Cast a spell on the brave warrior Thorin
		and grant bardic inspiration to the wise wizard Gandalf
		and cast another spell on the elven archer Legolas
	`)

	
	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			userQuestion,
		},
		// IMPORTANT: try this:
		//ParallelToolCalls: openai.Bool(false), // default value
		ParallelToolCalls: openai.Bool(true), // Sequential tool calls
		Tools:             tools,
		Model:             model,
		Temperature:       openai.Opt(0.0),
	}

	// Make [COMPLETION] request
	completion, err := client.Chat.Completions.New(ctx, params)
	if err != nil {
		panic(err)
	}

	// TOOL CALLS: Extract tool calls from the response
	toolCalls := completion.Choices[0].Message.ToolCalls

	// Return early if there are no tool calls
	if len(toolCalls) == 0 {
		fmt.Println("😡 No function call")
		fmt.Println()
		return
	}

	fmt.Println(strings.Repeat("=", 80))

	// Display the tool calls
	for _, toolCall := range toolCalls {
		fmt.Println(toolCall.Function.Name, toolCall.Function.Arguments)
	}

	fmt.Println(strings.Repeat("=", 80))
}
