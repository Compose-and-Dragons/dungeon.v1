package main

import (
	"context"
	"dungeon-master/agents"
	"encoding/json"
	"fmt"
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
	// [MCP Client] to connect to the [MCP Dungeon Server]
	mcpClient, err := tools.NewStreamableHttpMCPClient(ctx, mcpHost)
	if err != nil {
		panic(fmt.Errorf("failed to create MCP client: %v", err))
	}

	ui.Println(ui.Orange, "MCP Client initialized successfully")

	// ---------------------------------------------------------
	// TOOLS CATALOG: get the [list of tools] from the [MCP Client]
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

	DisplayToolsIndex(toolsIndex)

	// ---------------------------------------------------------
	// AGENT: This is the Dungeon Master Agent using tools
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
			Tools:             toolsIndex, // <- important 
			ParallelToolCalls: openai.Opt(false),
		}),
	)
	if err != nil {
		panic(err)
	}

	// SYSTEM MESSAGE:
	instructions := fmt.Sprintf(`Your name is "%s the Dungeon Master".`, dungeonMasterToolsAgentName) + "\n" + helpers.GetEnvOrDefault("DUNGEON_MASTER_SYSTEM_INSTRUCTIONS", dungeonMasterToolsAgentName)
	dungeonMasterSystemInstructions := openai.SystemMessage(instructions)

	// ---------------------------------------------------------
	// [FAKE] AGENT: This is the Ghost agent
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
	// [REMOTE] AGENT: This is the Boss agent
	// ---------------------------------------------------------
	bossAgent := agents.NewBossAgent(
		helpers.GetEnvOrDefault("BOSS_NAME", "Boss"),
		helpers.GetEnvOrDefault("BOSS_REMOTE_AGENT_URL", "http://localhost:8080/agent/boss"),
	)

	// ---------------------------------------------------------
	// AGENTS: Creating the Agents Team of the Dungeon
	// ---------------------------------------------------------
	idDungeonMasterToolsAgent := strings.ToLower(dungeonMasterToolsAgentName)
	idGhostAgent := strings.ToLower(ghostAgentName)
	idGuardAgent := strings.ToLower(guardAgent.GetName())
	idSorcererAgent := strings.ToLower(sorcererAgent.GetName())
	idMerchantAgent := strings.ToLower(merchantAgent.GetName())
	idHealerAgent := strings.ToLower(healerAgent.GetName())
	idBossAgent := strings.ToLower(bossAgent.GetName())

	agentsTeam = map[string]mu.Agent{
		idDungeonMasterToolsAgent: dungeonMasterToolsAgent,
		idGhostAgent:              ghostAgent,
		idGuardAgent:              guardAgent,
		idSorcererAgent:           sorcererAgent,
		idMerchantAgent:           merchantAgent,
		idHealerAgent:             healerAgent,
		idBossAgent:               bossAgent,
	}
	selectedAgent = agentsTeam[idDungeonMasterToolsAgent]

	DisplayAgentsTeam()
	// Loop to interact with the agents
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

		// DEBUG:
		if strings.HasPrefix(content.Input, "/memory") {
			msg.DisplayHistory(dungeonMasterToolsAgent)
			continue
		}

		// ---------------------------------------------------------
		// Get the AGENTS team list
		// ---------------------------------------------------------
		if strings.HasPrefix(content.Input, "/agents") {
			DisplayAgentsTeam()
			continue
		}

		if strings.HasPrefix(content.Input, "/tools") {
			DisplayToolsIndex(toolsIndex)
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
			// Tool execution callback if detected by the agent
			executeFn := ExecuteFunction(mcpClient, thinkingCtrl)

			dungeonMasterMessages := []openai.ChatCompletionMessageParamUnion{
				dungeonMasterSystemInstructions,
				userMessage,
			}
			
			

			// TOOLS DETECTION:
			_, toolCallsResults, assistantMessage, err := selectedAgent.DetectToolCalls(dungeonMasterMessages, executeFn)
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
			guardAgentMessages, err := GeneratePromptMessagesWithSimilarities(ctx, &client, guardAgent.GetName(), content.Input, similaritySearchLimit, similaritySearchMaxResults)

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
			sorcererAgentMessages, err := GeneratePromptMessagesWithSimilarities(ctx, &client, sorcererAgent.GetName(), content.Input, similaritySearchLimit, similaritySearchMaxResults)

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
			merchantAgentMessages, err := GeneratePromptMessagesWithSimilarities(ctx, &client, merchantAgent.GetName(), content.Input, similaritySearchLimit, similaritySearchMaxResults)

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
			healerAgentMessages, err := GeneratePromptMessagesWithSimilarities(ctx, &client, healerAgent.GetName(), content.Input, similaritySearchLimit, similaritySearchMaxResults)

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

		// ---------------------------------------------------------
		// TALK TO: AGENT:: **BOSS**
		// ---------------------------------------------------------
		case bossAgent.GetName():
			ui.Println(ui.Red, "<", selectedAgent.GetName(), "speaking...>")

			bossAgentMessages := []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(content.Input),
			}

			answer, err := selectedAgent.RunStream(bossAgentMessages, func(content string) error {
				fmt.Print(content)
				return nil
			})

			// IMPORTANT: Check if the player has defeated the boss
			// 👀 Look at /dungeon-end-of-level-boss/data/boss_system_instructions.md
			// You lose 😢
			if strings.Contains(strings.ToLower(answer), "you are trapped") {
				ui.Println(ui.Red, "\n💀 You have been defeated by the Boss! Game Over! 💀")
				ui.Println(ui.Red, "👹 The Boss reigns supreme in the dungeon! 👹")
				ui.Println(ui.Red, "🎲 Better luck next time! 🎲")

				// [DIRECT CALL TO MCP]
				result, err := mcpClient.CallTool(ctx, "get_player_info", "{}")
				if err == nil && len(result.Content) > 0 {
					playerInfo := result.Content[0].(mcp.TextContent).Text
					ui.Println(ui.Red, "📝 Your player information:\n", playerInfo)
				}

				continue
				//break
			}
			// You win 🎉
			if strings.Contains(strings.ToLower(answer), "you are free") {
				ui.Println(ui.Green, "\n💀 You have defeated the Boss! Congratulations, brave adventurer! 💀")
				ui.Println(ui.Green, "👑 You are now the new ruler of the dungeon! 👑")
				ui.Println(ui.Green, "🎉 Thanks for playing! 🎉")

				// [DIRECT CALL TO MCP]
				result, err := mcpClient.CallTool(ctx, "get_player_info", "{}")
				if err == nil && len(result.Content) > 0 {
					playerInfo := result.Content[0].(mcp.TextContent).Text
					ui.Println(ui.Green, "📝 Your player information:\n", playerInfo)
				}

				continue
				//break
			}

			if err != nil {
				ui.Println(ui.Red, "Error:", err)
			}

			fmt.Println()
			fmt.Println()

		default:
			ui.Printf(ui.Cyan, "\n🤖 %s is thinking...\n", selectedAgent.GetName())
		}
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

				invokedAgent := agentsTeam[strings.ToLower(argumentsStructured.Name)]
				if invokedAgent == nil {
					return fmt.Sprintf(`{"result": "😕 There is no NPC named %s"}`, argumentsStructured.Name), nil
				}
				// NOTE: SPECIAL CASE: Ghost agent is always available
				if strings.EqualFold(invokedAgent.GetName(), "Casper") {
					selectedAgent = agentsTeam[strings.ToLower(argumentsStructured.Name)]
					return fmt.Sprintf(`{"result": "😃 You speak to %s."}`, arguments), nil
				}

				// THIS IS FOR TEST: IMPORTANT: TODO: TO BE REMOVED NOTE:
				// SPECIAL CASE: Shesepankh agent is always available
				if strings.EqualFold(invokedAgent.GetName(), "Shesepankh") {
					selectedAgent = agentsTeam[strings.ToLower(argumentsStructured.Name)]
					return fmt.Sprintf(`{"result": "😃 You speak to %s."}`, arguments), nil
				}

				// ===================================================================================
				// IMPORTANT: check the position of the agent in the dungeon and the player position
				// NOTE: make a **direct call** to the MCP server to invoke a tool
				// ===================================================================================
				if mcpClient != nil {
					ctx := context.Background()

					// [DIRECT CALL TO MCP]
					result, err := mcpClient.CallTool(ctx, "is_player_in_same_room_as_npc", arguments)

					if err != nil {
						fmt.Println("🔴 Error calling tool is_player_in_same_room_as_npc:", err)
						return fmt.Sprintf(`{"result": "😕 You cannot speak to %s. (%s)"}`, argumentsStructured.Name, err.Error()), nil

					}
					if len(result.Content) > 0 {
						var toolResponse struct {
							InSameRoom bool   `json:"in_same_room"`
							PlayerRoom string `json:"player_room_id"`
							NPCRoom    string `json:"npc_room_id,omitempty"`
							Message    string `json:"message"`
						}
						err = json.Unmarshal([]byte(result.Content[0].(mcp.TextContent).Text), &toolResponse)
						if err != nil {
							fmt.Println("🔴 Error unmarshaling tool response:", err)
							return fmt.Sprintf(`{"result": "😕 You cannot speak to %s. (%s)"}`, argumentsStructured.Name, err.Error()), nil
						}
						if !toolResponse.InSameRoom {
							return fmt.Sprintf(`{"result": "😕 You cannot speak to %s because you are not in the same room."}`, argumentsStructured.Name), nil
						}

						// NOTE: the player is in the same room as the NPC
						// Check if the agent exist in the team
						// invokedAgent := agentsTeam[strings.ToLower(argumentsStructured.Name)]

						// if invokedAgent == nil {
						// 	return fmt.Sprintf(`{"result": "😕 There is no NPC named %s"}`, argumentsStructured.Name), nil
						// } else {
						// 	selectedAgent = agentsTeam[strings.ToLower(argumentsStructured.Name)]
						// }
						selectedAgent = agentsTeam[strings.ToLower(argumentsStructured.Name)]
						// Use the /dm command to go back to the Dungeon Master
						return fmt.Sprintf(`{"result": "😃 You speak to %s. They greet you warmly and are eager to assist you on your quest."}`, arguments), nil

					}

				}
				return fmt.Sprintf(`{"result": "😕 You cannot speak to %s. (%s)"}`, argumentsStructured.Name, "unable to connect to the MCP Server"), nil

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

func GeneratePromptMessagesWithSimilarities(ctx context.Context, client *openai.Client, agentName, input string, similarityLimit float64, maxResults int) ([]openai.ChatCompletionMessageParamUnion, error) {
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

func DisplayMCPToolCallResult(results []string) {
	fmt.Println(strings.Repeat("-", 3) + "[MCP RESPONSE]" + strings.Repeat("-", 33))
	fmt.Println(results[0])
	fmt.Println(strings.Repeat("-", 50))
}

func DisplayToolsIndex(toolsIndex []openai.ChatCompletionToolUnionParam) {
	for _, tool := range toolsIndex {
		ui.Printf(ui.Magenta, "Tool: %s - %s\n", tool.GetFunction().Name, tool.GetFunction().Description)
	}
	fmt.Println()
}

func DisplayAgentsTeam() {
	for agentId, agent := range agentsTeam {
		ui.Printf(ui.Cyan, "Agent ID: %s agent name: %s model: %s\n", agentId, agent.GetName(), agent.GetModel())
	}
	fmt.Println()
}

func DisplayDMResponse(assistantMessage string) {
	ui.Println(ui.Green, strings.Repeat("-", 3)+"[DM RESPONSE]"+strings.Repeat("-", 34))
	fmt.Println(assistantMessage)
	ui.Println(ui.Green, strings.Repeat("-", 50))
	fmt.Println()
}
