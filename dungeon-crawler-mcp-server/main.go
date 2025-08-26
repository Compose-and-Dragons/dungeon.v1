package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"dungeon-mcp-server/helpers"
	"dungeon-mcp-server/tools"
	"dungeon-mcp-server/types"

	"github.com/mark3labs/mcp-go/server"
	"github.com/micro-agent/micro-agent-go/agent/mu"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	// imported as openai
)


func main() {

	// ---------------------------------------------------------
	// Create MCP server
	// ---------------------------------------------------------
	s := server.NewMCPServer(
		"dungeon-mcp-server",
		"0.0.0",
	)

	// ---------------------------------------------------------
	// Create a "micro" agent
	// ---------------------------------------------------------
	ctx := context.Background()

	fmt.Println("🤖 Initializing Dungeon Agent...", os.Getenv("MODEL_RUNNER_BASE_URL"), "+++")

	baseURL := helpers.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")
	fmt.Println("🔵🌍 Model Runner Base URL:", baseURL)
	dungeonModel := helpers.GetEnvOrDefault("DUNGEON_MODEL", "ai/qwen2.5:1.5B-F16")
	fmt.Println("🌍 Dungeon Model:", dungeonModel)
	// Initialize OpenAI client
	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey(""),
	)

	temperature := helpers.StringToFloat(helpers.GetEnvOrDefault("DUNGEON_MODEL_TEMPERATURE", "0.7"))

	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": map[string]any{
				"type": "string",
			},
			"description": map[string]any{
				"type": "string",
			},
		},
		"required": []string{"name", "description"},
	}

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "room_info",
		Description: openai.String("name and description of the room"),
		Schema:      schema,
		Strict:      openai.Bool(true),
	}

	dungeonAgent, err := mu.NewAgent(ctx, "dungeon-agent",
		mu.WithClient(client),
		mu.WithParams(openai.ChatCompletionNewParams{
			Model:       dungeonModel,
			Temperature: openai.Opt(temperature),
			ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
				OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
					JSONSchema: schemaParam,
				},
			},
		}),
	)
	if err != nil {
		panic(err)
	}

	// ---------------------------------------------------------
	// Game initialisation
	// ---------------------------------------------------------
	currentPlayer := types.Player{
		Name: "Unknown",
	}

	width := helpers.StringToInt(helpers.GetEnvOrDefault("DUNGEON_WIDTH", "3"))
	height := helpers.StringToInt(helpers.GetEnvOrDefault("DUNGEON_HEIGHT", "3"))
	entranceX := helpers.StringToInt(helpers.GetEnvOrDefault("DUNGEON_ENTRANCE_X", "0"))
	entranceY := helpers.StringToInt(helpers.GetEnvOrDefault("DUNGEON_ENTRANCE_Y", "0"))
	exitX := helpers.StringToInt(helpers.GetEnvOrDefault("DUNGEON_EXIT_X", "2"))
	exitY := helpers.StringToInt(helpers.GetEnvOrDefault("DUNGEON_EXIT_Y", "2"))

	dungeonName := helpers.GetEnvOrDefault("DUNGEON_NAME", "The Dark Labyrinth")
	dungeonDescription := helpers.GetEnvOrDefault("DUNGEON_DESCRIPTION", "A sprawling underground maze filled with monsters, traps, and treasure.")

	fmt.Println("🧙 Dungeon Name:", dungeonName)
	fmt.Println("📝 Dungeon Description:", dungeonDescription)

	fmt.Println("🏰 Dungeon Size:", width, "x", height)

	dungeon := types.Dungeon{
		Name:        dungeonName,
		Description: dungeonDescription,
		Width:       width,
		Height:      height,
		Rooms:       []types.Room{},
		EntranceCoords: types.Coordinates{
			X: entranceX,
			Y: entranceY,
		},
		ExitCoords: types.Coordinates{
			X: exitX,
			Y: exitY,
		},
	}
	// TODO:
	// make the dungeon settings configurable via env vars or a config file
	fmt.Println("🚪 Dungeon Entrance Coords:", dungeon.EntranceCoords)
	fmt.Println("🚪 Dungeon Exit Coords:", dungeon.ExitCoords)

	// Create the entrance room of the dungeon
	// TODO: generate room name and description with a model

	// ---------------------------------------------------------
	// BEGIN: Generate the entrance room with the dungeon agent
	// ---------------------------------------------------------
	dungeonAgentSystemInstruction := helpers.GetEnvOrDefault("DUNGEON_AGENT_SYSTEM_INSTRUCTION", "You are a Dungeon Master. You create rooms in a dungeon. Each room has a name and a short description.")

	response, err := dungeonAgent.Run([]openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(dungeonAgentSystemInstruction),
		openai.UserMessage("Create an dungeon entrance room with a name and a short description."),
	})

	if err != nil {
		fmt.Println("🔴 Error generating room:", err)
		return
	}

	// following the schema response is a JSON String
	// {
	//   "room_info": {
	//     "name": "Dungeon Entrance",
	//     "description": "The dark and foreboding entrance to the dungeon."
	//   }

	fmt.Println("📝 Dungeon Entrance Room Response:", response)

	var roomResponse struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	err = json.Unmarshal([]byte(response), &roomResponse)
	if err != nil {
		fmt.Println("Error unmarshaling room response:", err)
		return
	}

	fmt.Println("👋🏰 Entrance Room:", roomResponse)
	// ---------------------------------------------------------
	// END: of Generate the entrance room with the dungeon agent
	// ---------------------------------------------------------

	entranceRoom := types.Room{
		ID:          "room_0_0",
		Name:        roomResponse.Name,
		Description: roomResponse.Description,
		IsEntrance:  true,
		IsExit:      false,
		Coordinates: types.Coordinates{
			X: entranceX,
			Y: entranceY,
		},
		Visited:     false,
		HasMonster:  false,
		HasNPC:      false,
		HasTreasure: false,
	}
	dungeon.Rooms = append(dungeon.Rooms, entranceRoom)

	// ---------------------------------------------------------
	// TOOLS:
	// ---------------------------------------------------------
	// TODO:
	// - add dungeon map tool
	// - add look around tool or get current room info tool
	// ---------------------------------------------------------
	// Register tools and their handlers
	// ---------------------------------------------------------
	// Create Player
	createPlayerToolInstance := tools.CreatePlayerTool()
	s.AddTool(createPlayerToolInstance, tools.CreatePlayerToolHandler(&currentPlayer, &dungeon))

	// Get Player Info
	getPlayerInfoToolInstance := tools.GetPlayerInformationTool()
	s.AddTool(getPlayerInfoToolInstance, tools.GetPlayerInformationToolHandler(&currentPlayer, &dungeon))

	// Get Dungeon Info
	getDungeonInfoToolInstance := tools.GetDungeonInformationTool()
	s.AddTool(getDungeonInfoToolInstance, tools.GetDungeonInformationToolHandler(&currentPlayer, &dungeon))

	// Move in the dungeon
	moveIntoTheDungeonToolInstance := tools.GetMoveIntoTheDungeonTool()
	s.AddTool(moveIntoTheDungeonToolInstance, tools.MoveByDirectionToolHandler(&currentPlayer, &dungeon, dungeonAgent))

	// Get Current Room Info
	getCurrentRoomInfoToolInstance := tools.GetCurrentRoomInformationTool()
	s.AddTool(getCurrentRoomInfoToolInstance, tools.GetCurrentRoomInformationToolHandler(&currentPlayer, &dungeon))

	// ---------------------------------------------------------
	// Start the HTTP server
	// ---------------------------------------------------------
	httpPort := helpers.GetEnvOrDefault("MCP_HTTP_PORT", "9090")
	fmt.Println("🌍 MCP HTTP Port:", httpPort)

	log.Println("[Dungeon]MCP StreamableHTTP server is running on port", httpPort)

	// Create a custom mux to handle both MCP and health endpoints
	mux := http.NewServeMux()
	// Add healthcheck endpoint (for Docker MCP Gateway with Docker Compose)
	mux.HandleFunc("/health", healthCheckHandler)
	// Add MCP endpoint
	httpServer := server.NewStreamableHTTPServer(s,
		server.WithEndpointPath("/mcp"),
	)
	// Register MCP handler with the mux
	mux.Handle("/mcp", httpServer)
	// Start the HTTP server with custom mux
	log.Fatal(http.ListenAndServe(":"+httpPort, mux))
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]any{
		"status": "healthy",
	}
	json.NewEncoder(w).Encode(response)
}
