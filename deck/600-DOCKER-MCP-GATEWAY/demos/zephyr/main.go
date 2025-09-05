package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/micro-agent/micro-agent-go/agent/helpers"
	"github.com/micro-agent/micro-agent-go/agent/msg"
	"github.com/micro-agent/micro-agent-go/agent/mu"
	"github.com/micro-agent/micro-agent-go/agent/tools"
	"github.com/micro-agent/micro-agent-go/agent/ui"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

func main() {

	ctx := context.Background()
	llmURL := helpers.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")
	mcpHost := helpers.GetEnvOrDefault("MCP_HOST", "http://localhost:9011/mcp")
	dungeonMasterModel := helpers.GetEnvOrDefault("DUNGEON_MASTER_MODEL", "hf.co/menlo/jan-nano-gguf:q4_k_m")

	client := openai.NewClient(
		option.WithBaseURL(llmURL),
		option.WithAPIKey(""),
	)

	mcpClient, err := tools.NewStreamableHttpMCPClient(ctx, mcpHost)
	if err != nil {
		panic(fmt.Errorf("failed to create MCP client: %v", err))
	}

	ui.Println(ui.Orange, "MCP Client initialized successfully")

	// ---------------------------------------------------------
	// TOOLS CATALOG: get the list of tools from the [MCP] client
	// ---------------------------------------------------------
	toolsIndex := mcpClient.OpenAITools()

	DisplayToolsIndex(toolsIndex)

	// ---------------------------------------------------------
	// AGENT: This is the Dungeon Master agent using tools
	// ---------------------------------------------------------
	dungeonMasterToolsAgentName := helpers.GetEnvOrDefault("DUNGEON_MASTER_NAME", "Sam")
	dungeonMasterModeltemperature := helpers.StringToFloat(helpers.GetEnvOrDefault("DUNGEON_MASTER_MODEL_TEMPERATURE", "0.0"))

	dungeonMasterToolsAgent, err := mu.NewAgent(ctx, dungeonMasterToolsAgentName,
		mu.WithClient(client),
		mu.WithParams(openai.ChatCompletionNewParams{
			Model:       dungeonMasterModel,
			Temperature: openai.Opt(dungeonMasterModeltemperature),
			ToolChoice: openai.ChatCompletionToolChoiceOptionUnionParam{
				OfAuto: openai.String("auto"),
			},
			Tools:             toolsIndex,
			ParallelToolCalls: openai.Opt(false),
		}),
	)
	if err != nil {
		panic(err)
	}

	// SYSTEM MESSAGE:
	instructions := fmt.Sprintf(`Your name is "%s the Dungeon Master".`, dungeonMasterToolsAgentName) + "\n" + helpers.GetEnvOrDefault("DUNGEON_MASTER_SYSTEM_INSTRUCTIONS", dungeonMasterToolsAgentName)
	dungeonMasterSystemInstructions := openai.SystemMessage(instructions)

	for {
		promptText := "🤖 (/bye to exit) [" + dungeonMasterToolsAgent.GetName() + "]>"
		// PROMPT:
		content, _ := ui.SimplePrompt(promptText, "Type your command here...")

		// USER MESSAGE: content.Input
		userMessage := openai.UserMessage(content.Input)

		// ---------------------------------------------------------
		// Bye [COMMAND]
		// ---------------------------------------------------------
		if strings.HasPrefix(content.Input, "/bye") {
			fmt.Println("👋 Goodbye! Thanks for playing!")
			break
		}

		// DEBUG:
		if strings.HasPrefix(content.Input, "/memory") {
			msg.DisplayHistory(dungeonMasterToolsAgent)
			continue
		}

		// ---------------------------------------------------------
		// TALK TO: AGENT:: **SORCERER** + [RAG]
		// ---------------------------------------------------------
		ui.Println(ui.Purple, "<", dungeonMasterToolsAgent.GetName(), "speaking...>")

		thinkingCtrl := ui.NewThinkingController()
		thinkingCtrl.Start(ui.Cyan, "Tools detection.....")

		// Create executeFunction with MCP client option
		// Tool execution callback
		executeFn := ExecuteFunction(mcpClient, thinkingCtrl)

		dungeonMasterMessages := []openai.ChatCompletionMessageParamUnion{
			dungeonMasterSystemInstructions,
			userMessage,
		}
		// QUESTION: should I keep the last message?
		// QUESTION: should I reset the messages to only keep the last message + system?

		// TOOLS DETECTION:
		_, toolCallsResults, assistantMessage, err := dungeonMasterToolsAgent.DetectToolCalls(dungeonMasterMessages, executeFn)
		if err != nil {
			panic(err)
		}

		thinkingCtrl.Stop()

		if len(toolCallsResults) > 0 {
			// IMPORTANT: This is the answer from the [MCP] server
			DisplayMCPToolCallResult(toolCallsResults)
		}

		// ASSISTANT MESSAGE:
		// This is the final answer from the agent
		DisplayDMResponse(assistantMessage)

	}

}

func ExecuteFunction(mcpClient *tools.MCPClient, thinkingCtrl *ui.ThinkingController) func(string, string) (string, error) {

	return func(functionName string, arguments string) (string, error) {

		fmt.Printf("🟢 %s with arguments: %s\n", functionName, arguments)

		// WAITING: for [CONFIRMATION] function is detected, function execution confirmation
		thinkingCtrl.Pause()
		choice := ui.GetChoice(ui.Yellow, "Do you want to execute this function? (y)es (n)o (a)bort", []string{"y", "n", "a"}, "y")
		thinkingCtrl.Resume()

		switch choice {
		case "n":
			return `{"result": "Function not executed"}`, nil
		case "a": // abort
			return `{"result": "Function not executed"}`,
				&mu.ExitToolCallsLoopError{Message: "Tool execution aborted by user"}

		default: // [YES] if the user confirms (yes)

			switch functionName {

			// ---------------------------------------------------------
			// [MCP] TOOL CALLS: implementation
			// ---------------------------------------------------------
			default:
				// If MCP client is available, use it to execute the tool
				if mcpClient != nil {
					ctx := context.Background()
					result, err := mcpClient.CallTool(ctx, functionName, arguments)
					if err != nil {
						return "", fmt.Errorf("MCP tool execution failed: %v", err)
					}
					if len(result.Content) > 0 {
						// Take the first content item and return its text
						resultContent := result.Content[0].(mcp.TextContent).Text
						fmt.Println("✅ Tool executed successfully")
						// ✋ could be JSON or not
						return resultContent, nil

					}
					return `{"result": "Tool executed successfully but returned no content"}`, nil
				}
				return `{"result": "Function not executed"}`, nil
			}

		}
	}
}

func DisplayToolsIndex(toolsIndex []openai.ChatCompletionToolUnionParam) {
	for _, tool := range toolsIndex {
		ui.Printf(ui.Magenta, "Tool: %s - %s\n", tool.GetFunction().Name, tool.GetFunction().Description)
	}
	fmt.Println()
}

func DisplayMCPToolCallResult(results []string) {
	fmt.Println(strings.Repeat("-", 3) + "[MCP RESPONSE]" + strings.Repeat("-", 33))
	fmt.Println(results[0])
	fmt.Println(strings.Repeat("-", 50))
}

func DisplayDMResponse(assistantMessage string) {
	ui.Println(ui.Green, strings.Repeat("-", 3)+"[DM RESPONSE]"+strings.Repeat("-", 34))
	fmt.Println(assistantMessage)
	ui.Println(ui.Green, strings.Repeat("-", 50))
	fmt.Println()
}
