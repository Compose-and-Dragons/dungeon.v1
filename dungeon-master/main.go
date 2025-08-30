package main

import (
	"context"
	"dungeon-master/agents"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/micro-agent/micro-agent-go/agent/helpers"
	"github.com/micro-agent/micro-agent-go/agent/msg"
	"github.com/micro-agent/micro-agent-go/agent/mu"
	"github.com/micro-agent/micro-agent-go/agent/tools"
	"github.com/micro-agent/micro-agent-go/agent/ui"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	"github.com/openai/openai-go/v2/shared"
)

var agentsTeam map[string]mu.Agent
var selectedAgent mu.Agent

func main() {

	ctx := context.Background()

	llmURL := helpers.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")
	mcpHost := helpers.GetEnvOrDefault("MCP_HOST", "http://localhost:9011/mcp")
	dungeonMasterModel := helpers.GetEnvOrDefault("DUNGEON_MASTER_MODEL", "hf.co/menlo/jan-nano-gguf:q4_k_m")

	fmt.Println("🌍 LLM URL:", llmURL)
	fmt.Println("🌍 MCP Host:", mcpHost)
	fmt.Println("🌍 Dungeon Master Model:", dungeonMasterModel)

	similaritySearchLimit := helpers.StringToFloat(helpers.GetEnvOrDefault("SIMILARITY_LIMIT", "0.5"))
	similaritySearchMaxResults := helpers.StringToInt(helpers.GetEnvOrDefault("SIMILARITY_MAX_RESULTS", "2"))

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

	// ---------------------------------------------------------
	// TOOLS: adding tools to the mcp tools index
	// ---------------------------------------------------------
	speakToAnAgentTool := openai.ChatCompletionFunctionTool(shared.FunctionDefinitionParam{
		Name:        "speak_to_somebody",
		Description: openai.String("Speak to somebody by name"),
		Parameters: shared.FunctionParameters{
			"type": "object",
			"properties": map[string]interface{}{
				"name": map[string]string{
					"type":        "string",
					"description": "The name of the person to speak to",
				},
			},
			"required": []string{"name"},
		},
	})

	toolsIndex = append(toolsIndex, speakToAnAgentTool)

	displayToolsIndex(toolsIndex)

	// ---------------------------------------------------------
	// AGENT: This is the Dungeon Master agent using tools
	// ---------------------------------------------------------
	dungeonMasterToolsAgentName := helpers.GetEnvOrDefault("DUNGEON_MASTER_NAME", "Sam")

	dungeonMasterToolsAgent, err := mu.NewAgent(ctx, dungeonMasterToolsAgentName,
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

	// SYSTEM MESSAGE:
	instructions := fmt.Sprintf(`Your name is "%s the Dungeon Master".`, dungeonMasterToolsAgentName) + "\n" + helpers.GetEnvOrDefault("DUNGEON_MASTER_SYSTEM_INSTRUCTIONS", dungeonMasterToolsAgentName)
	dungeonMasterSystemInstructions := openai.SystemMessage(instructions)

	// note used but could be useful later
	conversationalMemory := []openai.ChatCompletionMessageParamUnion{
		dungeonMasterSystemInstructions,
	}

	// ---------------------------------------------------------
	// AGENT: This is the Ghost agent
	// ---------------------------------------------------------
	// REMARK: Ghost agent is for testing only.
	// IMPORTANT: keep the name simple in only one word
	ghostAgentName := "Casper"
	ghostAgent := agents.NewGhostAgent(ghostAgentName)

	// ---------------------------------------------------------
	// AGENT: This is the Guard agent
	// ---------------------------------------------------------
	guardAgent := agents.GetGuardAgent(ctx, client)

	// ---------------------------------------------------------
	// AGENT: This is the Sorcerer agent
	// ---------------------------------------------------------
	sorcererAgent := agents.GetSorcererAgent(ctx, client)

	// ---------------------------------------------------------
	// AGENT: This is the Merchant agent
	// ---------------------------------------------------------
	merchantAgent := agents.GetMerchantAgent(ctx, client)

	// ---------------------------------------------------------
	// AGENT: This is the Healer agent
	// ---------------------------------------------------------
	healerAgent := agents.GetHealerAgent(ctx, client)

	// ---------------------------------------------------------
	// AGENTS: Creating the agents team of the dungeon
	// ---------------------------------------------------------
	idDungeonMasterToolsAgent := strings.ToLower(dungeonMasterToolsAgentName)
	idGhostAgent := strings.ToLower(ghostAgentName)
	idGuardAgent := strings.ToLower(guardAgent.GetName())
	idSorcererAgent := strings.ToLower(sorcererAgent.GetName())
	idMerchantAgent := strings.ToLower(merchantAgent.GetName())
	idHealerAgent := strings.ToLower(healerAgent.GetName())

	agentsTeam = map[string]mu.Agent{
		idDungeonMasterToolsAgent: dungeonMasterToolsAgent,
		idGhostAgent:              ghostAgent,
		idGuardAgent:              guardAgent,
		idSorcererAgent:           sorcererAgent,
		idMerchantAgent:           merchantAgent,
		idHealerAgent:             healerAgent,
	}
	selectedAgent = agentsTeam[idDungeonMasterToolsAgent]

	displayAgentsTeam()

	for {
		var promptText string
		if selectedAgent.GetName() == dungeonMasterToolsAgentName {
			// Dungeon Master prompt
			promptText = "🤖 (/bye to exit) [" + selectedAgent.GetName() + "]>"
		} else {
			// Non Player Character prompt
			promptText = "🙂 (/bye to exit /dm to go back to the DM) [" + selectedAgent.GetName() + "]>"
		}

		// PROMPT:
		content, _ := ui.SimplePrompt(promptText, "Type your command here...")

		// USER MESSAGE:
		userMessage := openai.UserMessage(content.Input)

		// ---------------------------------------------------------
		// Bye [COMMAND]
		// ---------------------------------------------------------
		if strings.HasPrefix(content.Input, "/bye") {
			fmt.Println("👋 Goodbye! Thanks for playing!")
			break
		}

		// ---------------------------------------------------------
		// Back to the Dungeon Master [COMMAND]
		// ---------------------------------------------------------
		if strings.HasPrefix(content.Input, "/back-to-dm") || strings.HasPrefix(content.Input, "/dm") || strings.HasPrefix(content.Input, "/dungeonmaster") && selectedAgent.GetName() != dungeonMasterToolsAgentName {
			selectedAgent = agentsTeam[idDungeonMasterToolsAgent]
			ui.Println(ui.Pink, "👋 You are back to the Dungeon Master:", selectedAgent.GetName())
			continue
			/*
				In Go, the continue keyword in a loop immediately jumps to the next iteration of the loop, skipping the rest
				of the code in the current iteration.

				Specifically:
				- In a for loop, continue returns to the beginning of the loop for the next iteration
				- Code after continue in the same iteration is not executed
				- The loop condition is evaluated normally
			*/
		}

		// ---------------------------------------------------------
		// For DEBUG: [COMMAND] to print messages history
		// ---------------------------------------------------------
		if strings.HasPrefix(content.Input, "/messages") {

			fmt.Println("📝 Messages history / Conversational memory:")
			for i, message := range conversationalMemory {
				printableMessage, err := msg.MessageToMap(message)
				if err != nil {
					fmt.Printf("Error converting message to map: %v\n", err)
					continue
				}
				fmt.Println("-", i, printableMessage)
			}
			continue
		}

		conversationalMemory = append(conversationalMemory, userMessage)

		// ---------------------------------------------------------
		// Get the AGENTS team list
		// ---------------------------------------------------------
		if strings.HasPrefix(content.Input, "/agents") {
			displayAgentsTeam()
			continue
		}

		if strings.HasPrefix(content.Input, "/tools") {
			displayToolsIndex(toolsIndex)
			continue
		}

		switch selectedAgent.GetName() {
		// ---------------------------------------------------------
		// TALK TO: AGENT: Dungeon master [COMPLETION] with [TOOLS]
		// ---------------------------------------------------------
		case dungeonMasterToolsAgentName:
			ui.Println(ui.Yellow, "<", selectedAgent.GetName(), "speaking...>")

			thinkingCtrl := ui.NewThinkingController()
			//thinkingCtrl.Start(ui.Blue, "Tools detection.....")
			thinkingCtrl.Start(ui.Cyan, "Tools detection.....")

			// Create executeFunction with MCP client option
			// Tool execution callback
			executeFn := executeFunction(mcpClient, thinkingCtrl)

			dungeonMasterMessages := []openai.ChatCompletionMessageParamUnion{
				dungeonMasterSystemInstructions,
				userMessage,
			}
			// QUESTION: should I keep the last message?

			// TOOLS DETECTION:
			_, toolCallsResults, assistantMessage, err := selectedAgent.DetectToolCalls(dungeonMasterMessages, executeFn)
			if err != nil {
				panic(err)
			}

			thinkingCtrl.Stop()

			if len(toolCallsResults) > 0 {
				// IMPORTANT: This is the answer from the [MCP] server
				displayFirstToolCallResult(toolCallsResults)
			}

			// ASSISTANT MESSAGE:
			// This is the final answer from the agent

			// TODO: improve the formatting of the assistant message
			formattedMessage := formatAssistantMessage(assistantMessage)
			ui.Println(ui.Green, formattedMessage)
			fmt.Println()

			// ui.PrintMarkdown(formattedMessage)
			// fmt.Println()

			// not used but could be useful later
			conversationalMemory = append(conversationalMemory, openai.AssistantMessage(assistantMessage))

		// ---------------------------------------------------------
		// TALK TO: AGENT:: **GHOST** for [TESTING] only
		// ---------------------------------------------------------
		case ghostAgentName:
			ui.Println(ui.Orange, "<", selectedAgent.GetName(), "speaking...>")

			ghostAgentMessages := []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage("You are a friendly and helpful ghost."),
				openai.UserMessage(content.Input),
			}
			_, err := selectedAgent.RunStream(ghostAgentMessages, func(content string) error {
				fmt.Print(content)
				return nil
			})

			if err != nil {
				ui.Println(ui.Red, "Error:", err)
			}

			fmt.Println()
			fmt.Println()

		// ---------------------------------------------------------
		// TALK TO: AGENT:: **GUARD** + [RAG]
		// ---------------------------------------------------------
		case guardAgent.GetName():
			ui.Println(ui.Brown, "<", selectedAgent.GetName(), "speaking...>")

			// ---------------------------------------------------------
			// [RAG] SIMILARITY SEARCH:
			// ---------------------------------------------------------
			guardAgentMessages, err := generatePromptMessagesWithSimilarities(ctx, &client, guardAgent.GetName(), content.Input, similaritySearchLimit, similaritySearchMaxResults)

			if err != nil {
				ui.Println(ui.Red, "Error:", err)
			}

			// NOTE: RunStreams adds the messages to the agent's memory
			_, err = selectedAgent.RunStream(guardAgentMessages, func(content string) error {
				fmt.Print(content)
				return nil
			})

			if err != nil {
				ui.Println(ui.Red, "Error:", err)
			}

			// DEBUG: display the messages history
			if strings.HasPrefix(content.Input, "/debug") {
				msg.DisplayHistory(selectedAgent)
			}

			fmt.Println()
			fmt.Println()

		// ---------------------------------------------------------
		// TALK TO: AGENT:: **SORCERER** + [RAG]
		// ---------------------------------------------------------
		case sorcererAgent.GetName():
			ui.Println(ui.Purple, "<", selectedAgent.GetName(), "speaking...>")

			// ---------------------------------------------------------
			// [RAG] SIMILARITY SEARCH:
			// ---------------------------------------------------------
			sorcererAgentMessages, err := generatePromptMessagesWithSimilarities(ctx, &client, sorcererAgent.GetName(), content.Input, similaritySearchLimit, similaritySearchMaxResults)

			if err != nil {
				ui.Println(ui.Red, "Error:", err)
			}

			// NOTE: RunStreams adds the messages to the agent's memory
			_, err = selectedAgent.RunStream(sorcererAgentMessages, func(content string) error {
				fmt.Print(content)
				return nil
			})

			if err != nil {
				ui.Println(ui.Red, "Error:", err)
			}

			// DEBUG: display the messages history
			if strings.HasPrefix(content.Input, "/debug") {
				msg.DisplayHistory(selectedAgent)
			}

			fmt.Println()
			fmt.Println()

		// ---------------------------------------------------------
		// TALK TO: AGENT:: **MERCHANT** + [RAG]
		// ---------------------------------------------------------
		case merchantAgent.GetName():
			ui.Println(ui.Cyan, "<", selectedAgent.GetName(), "speaking...>")

			// ---------------------------------------------------------
			// [RAG] SIMILARITY SEARCH:
			// ---------------------------------------------------------
			merchantAgentMessages, err := generatePromptMessagesWithSimilarities(ctx, &client, merchantAgent.GetName(), content.Input, similaritySearchLimit, similaritySearchMaxResults)

			if err != nil {
				ui.Println(ui.Red, "Error:", err)
			}

			// NOTE: RunStreams adds the messages to the agent's memory
			_, err = selectedAgent.RunStream(merchantAgentMessages, func(content string) error {
				fmt.Print(content)
				return nil
			})

			if err != nil {
				ui.Println(ui.Red, "Error:", err)
			}

			// DEBUG: display the messages history
			if strings.HasPrefix(content.Input, "/debug") {
				msg.DisplayHistory(selectedAgent)
			}

			fmt.Println()
			fmt.Println()

		// ---------------------------------------------------------
		// TALK TO: AGENT:: **HEALER** + [RAG]
		// ---------------------------------------------------------
		case healerAgent.GetName():
			ui.Println(ui.Magenta, "<", selectedAgent.GetName(), "speaking...>")

			// ---------------------------------------------------------
			// [RAG] SIMILARITY SEARCH:
			// ---------------------------------------------------------
			healerAgentMessages, err := generatePromptMessagesWithSimilarities(ctx, &client, healerAgent.GetName(), content.Input, similaritySearchLimit, similaritySearchMaxResults)

			if err != nil {
				ui.Println(ui.Red, "Error:", err)
			}

			// NOTE: RunStreams adds the messages to the agent's memory
			_, err = selectedAgent.RunStream(healerAgentMessages, func(content string) error {
				fmt.Print(content)
				return nil
			})

			if err != nil {
				ui.Println(ui.Red, "Error:", err)
			}

			// DEBUG: display the messages history
			if strings.HasPrefix(content.Input, "/debug") {
				msg.DisplayHistory(selectedAgent)
			}

			fmt.Println()
			fmt.Println()

		default:
			ui.Printf(ui.Cyan, "\n🤖 %s is thinking...\n", selectedAgent.GetName())
		}
	}
}


func executeFunction(mcpClient *tools.MCPClient, thinkingCtrl *ui.ThinkingController) func(string, string) (string, error) {

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
			// non MCP TOOL CALLS: implementations
			// ---------------------------------------------------------
			case "speak_to_somebody":

				argumentsStructured := struct {
					Name string `json:"name"`
				}{}
				err := json.Unmarshal([]byte(arguments), &argumentsStructured)
				if err != nil {
					return "", fmt.Errorf("failed to parse arguments: %v", err)
				}

				checkIfTheAgentExistInTheTeam := agentsTeam[strings.ToLower(argumentsStructured.Name)]

				if checkIfTheAgentExistInTheTeam == nil {
					return fmt.Sprintf(`{"result": "😕 There is no NPC named %s"}`, argumentsStructured.Name), nil
				} else {
					selectedAgent = agentsTeam[strings.ToLower(argumentsStructured.Name)]
				}
				// Use the /dm command to go back to the Dungeon Master
				return fmt.Sprintf(`{"result": "😃 You speak to %s. They greet you warmly and are eager to assist you on your quest."}`, arguments), nil

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

func generatePromptMessagesWithSimilarities(ctx context.Context, client *openai.Client, agentName, input string, similarityLimit float64, maxResults int) ([]openai.ChatCompletionMessageParamUnion, error) {
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

func displayFirstToolCallResult(results []string) {
	fmt.Println(strings.Repeat("-", 3)+"[MCP RESPONSE]"+ strings.Repeat("-", 33))
	fmt.Println(results[0])
	fmt.Println(strings.Repeat("-", 50))
}

func displayToolsIndex(toolsIndex []openai.ChatCompletionToolUnionParam) {
	for _, tool := range toolsIndex {
		ui.Printf(ui.Magenta, "Tool: %s - %s\n", tool.GetFunction().Name, tool.GetFunction().Description)
	}
	fmt.Println()
}

func displayAgentsTeam() {
	for agentId, agent := range agentsTeam {
		ui.Printf(ui.Cyan, "Agent ID: %s agent name: %s model: %s\n", agentId, agent.GetName(), agent.GetModel())
	}
	fmt.Println()
}

func formatAssistantMessage(message string) string {
	// Supprimer les espaces en début et fin
	message = strings.TrimSpace(message)
	
	// Remplacer les multiples sauts de ligne par un double saut de ligne
	re := regexp.MustCompile(`\n\s*\n\s*\n+`)
	message = re.ReplaceAllString(message, "\n\n")
	
	// Nettoyer les espaces multiples en fin de ligne
	re = regexp.MustCompile(`[ \t]+\n`)
	message = re.ReplaceAllString(message, "\n")
	
	// S'assurer qu'il n'y a pas plus de 2 sauts de ligne consécutifs
	re = regexp.MustCompile(`\n{3,}`)
	message = re.ReplaceAllString(message, "\n\n")
	
	return message
}