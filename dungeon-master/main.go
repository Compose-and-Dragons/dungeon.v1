package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/micro-agent/micro-agent-go/agent/mu"
	"github.com/micro-agent/micro-agent-go/agent/tools"
	"github.com/micro-agent/micro-agent-go/agent/ui"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

func main() {

	ctx := context.Background()

	llmURL := getEnvOrDefault("LLM_URL", "http://localhost:12434/engines/llama.cpp/v1")
	mcpHost := getEnvOrDefault("MCP_HOST", "http://localhost:9011/mcp")
	dungeonMasterModel := getEnvOrDefault("DUNGEON_MASTER_MODEL", "hf.co/menlo/jan-nano-gguf:q4_k_m")

	fmt.Println("🌍 LLM URL:", llmURL)
	fmt.Println("🌍 MCP Host:", mcpHost)
	fmt.Println("🌍 Dungeon Master Model:", dungeonMasterModel)

	client := openai.NewClient(
		option.WithBaseURL(llmURL),
		option.WithAPIKey(""),
	)

	mcpClient, err := tools.NewStreamableHttpMCPClient(ctx, mcpHost)
	if err != nil {
		panic(fmt.Errorf("failed to create MCP client: %v", err))
	}

	ui.Println(ui.Purple, "MCP Client initialized successfully")
	toolsIndex := mcpClient.OpenAITools()
	for _, tool := range toolsIndex {
		ui.Printf(ui.Magenta, "Tool: %s - %s\n", tool.GetFunction().Name, tool.GetFunction().Description)
	}

	toolAgent, err := mu.NewAgent(ctx, "Bob",
		mu.WithClient(client),
		mu.WithParams(openai.ChatCompletionNewParams{
			Model:       dungeonMasterModel,
			Temperature: openai.Opt(0.0),
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

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(`
			Your name is "Sam the Dungeon Master".
			You are a friendly and helpful Dungeon Master for a Dungeons & Dragons game.
			You will guide the player through a fantasy adventure, describing scenes, challenges, and characters.
			You will use tools to manage the game state, such as creating a player, starting a quest, and rolling dice.
			You will always describe the game world in vivid detail and engage the player with interesting scenarios.
			You will keep track of the player's stats, inventory, and progress through the adventure.
			You will respond to the player's actions and decisions, adapting the story accordingly.
			You will use a conversational and immersive style, making the player feel like they are part of the adventure.
			You will avoid breaking character and stay in the role of the Dungeon Master at all times.
			You will ensure the game is fun and exciting for the player.
			You will end each response with a question or prompt to encourage the player to take action.
			Always refer to the player by their name.
		`),
	}

	for {
		content, _ := ui.SimplePrompt("🤖 (/bye to exit)>", "Type your command here...")

		if content.Input == "/bye" {
			fmt.Println("👋 Goodbye! Thanks for playing!")
			break
		}

		userMessage := openai.UserMessage(content.Input)
		messages = append(messages, userMessage)

		// Stream callback for real-time content display
		// streamCallback := func(thinkingCtrl *ui.ThinkingController) func(string) error {
		// 	return func(content string) error {
		// 		if thinkingCtrl.IsStarted() {
		// 			thinkingCtrl.Stop()
		// 		}
		// 		ui.Print(ui.Green, content)
		// 		return nil
		// 	}
		// }

		thinkingCtrl := ui.NewThinkingController()
		//thinkingCtrl.Start(ui.Blue, "Tools detection.....")
		thinkingCtrl.Start(ui.Cyan, "Tools detection.....")

		// Create executeFunction with MCP client option
		// Tool execution callback
		executeFn := executeFunction(mcpClient, thinkingCtrl)

		_, toolCallsResults, assistantMessage, err := toolAgent.DetectToolCalls(messages, executeFn)
		if err != nil {
			panic(err)
		}

		thinkingCtrl.Stop()

		//prettyPrintFirstToolCallResult(toolCallsResults)
		displayFirstToolCallResult(toolCallsResults)

		ui.Println(ui.Green, assistantMessage)
		//fmt.Println(assistantMessage)
		fmt.Println()
	}

}

func executeFunction(mcpClient *tools.MCPClient, thinkingCtrl *ui.ThinkingController) func(string, string) (string, error) {

	return func(functionName string, arguments string) (string, error) {

		fmt.Printf("🟢 %s with arguments: %s\n", functionName, arguments)

		thinkingCtrl.Pause()
		//choice := ui.GetConfirmation(ui.Gray, "Do you want to execute this function?", true)
		//choice := ui.GetChoice(ui.Gray, "Do you want to execute this function? (y)es (n)o (a)bort", []string{"y", "n", "a"}, "y")
		choice := ui.GetChoice(ui.Yellow, "Do you want to execute this function? (y)es (n)o (a)bort", []string{"y", "n", "a"}, "y")

		thinkingCtrl.Resume()

		switch choice {
		case "n":
			return `{"result": "Function not executed"}`, nil
		case "a": // abort
			return `{"result": "Function not executed"}`,
				&mu.ExitToolCallsLoopError{Message: "Tool execution aborted by user"}

		default:

			// If MCP client is available, use it to execute the tool
			if mcpClient != nil {
				ctx := context.Background()
				result, err := mcpClient.CallTool(ctx, functionName, arguments)
				if err != nil {
					return "", fmt.Errorf("MCP tool execution failed: %v", err)
				}
				// resultContent = toolResponse.Content[0].(mcp.TextContent).Text
				// Convert MCP result to JSON string
				if len(result.Content) > 0 {
					// Take the first content item and return its text
					resultContent := result.Content[0].(mcp.TextContent).Text
					//fmt.Printf("✅ Tool executed successfully, result: %s\n", resultContent)
					fmt.Println("✅ Tool executed successfully")
					//return fmt.Sprintf(`{"result": "%s"}`, resultContent), nil
					// ✋ could be JSON or not
					return resultContent, nil

				}
				return `{"result": "Tool executed successfully but returned no content"}`, nil
			}
			return `{"result": "Function not executed"}`, nil
		}
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func prettyPrintFirstToolCallResult(results []string) {
	fmt.Println(strings.Repeat("-", 50))
	// results[0] is like {"result": "text"}
	// we want to extract the text value using JSON unmarshalling
	var resultMap map[string]string

	// Debug: print raw result
	fmt.Printf("Raw result: %q\n", results[0])

	cleanedResult := strings.ReplaceAll(results[0], "\n", "\\n")

	// Debug: print cleaned result
	fmt.Printf("Cleaned result: %q\n", cleanedResult)

	err := json.Unmarshal([]byte(cleanedResult), &resultMap)
	if err != nil {
		fmt.Println("Error unmarshalling result:", err)
	} else {
		if result, ok := resultMap["result"]; ok {
			fmt.Println(result)
		} else {
			fmt.Println("Error: result field not found in tool result")
		}
	}
	fmt.Println(strings.Repeat("-", 50))
}

func displayFirstToolCallResult(results []string) {
	fmt.Println(strings.Repeat("-", 50))
	fmt.Println(results[0])
	fmt.Println(strings.Repeat("-", 50))
}
