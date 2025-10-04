package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {

	// Create MCP server
	s := server.NewMCPServer(
		"mcp-dd",
		"0.0.0",
	)
	// Add a TOOL:
	chooseCharacterBySpecies := mcp.NewTool("choose_character_by_species",
		mcp.WithDescription(`select a species from among these: [Human, Orc, Elf, Dwarf] by saying: I want to talk to a <species_name>.`),
		mcp.WithString("species_name", // PARAMETER:
			mcp.Required(),
			mcp.Description("The species to detect in the user message. The species can be one of the following: [Human, Orc, Elf, Dwarf]."),
		),
	)
	s.AddTool(chooseCharacterBySpecies, chooseCharacterBySpeciesHandler)

	// Add another TOOL:
	detectTheRealTopicInUserMessage := mcp.NewTool("detect_real_topic_in_user_message",
		mcp.WithDescription(`select a topic from among these: [justice, war, combat, magic, poetry, craftsmanship, forge] by saying: I have a question about <topic_name>.`),
		mcp.WithString("topic_name", // PARAMETER:
			mcp.Required(),
			mcp.Description("The topic to detect in the user message. The topic can be one of the following: [justice, war, combat, magic, poetry, craftsmanship, forge]."),
		),
	)
	s.AddTool(detectTheRealTopicInUserMessage, detectTheRealTopicInUserMessageHandler)

	// ---------------------------------------------------------
	// Start the HTTP server
	// ---------------------------------------------------------
	httpPort := "9090"
	fmt.Println("🌍 MCP HTTP Port:", httpPort)

	log.Println("MCP StreamableHTTP server is running on port", httpPort)

	// Create a custom mux to handle both MCP and health endpoints
	mux := http.NewServeMux()
	// Add MCP endpoint
	httpServer := server.NewStreamableHTTPServer(s,
		server.WithEndpointPath("/mcp"),
	)
	// Register MCP handler with the mux
	mux.Handle("/mcp", httpServer)
	// Start the HTTP server with custom mux
	log.Fatal(http.ListenAndServe(":"+httpPort, mux))

}

func chooseCharacterBySpeciesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	fmt.Println("🛠️  chooseCharacterBySpeciesHandler called with args:", args)
	// Check if the species_name argument is provided
	if len(args) == 0 {
		return mcp.NewToolResultText("zephyr"), nil
	}
	var content = "zephyr" // default character
	if speciesName, ok := args["species_name"]; ok {

		switch strings.ToLower(speciesName.(string)) {
		case "dwarf":
			content = "galdor"
		case "elf":
			content = "thrain"
		case "human":
			content = "elara"
		case "half-elf":
			content = "liora"
		default:
			content = "zephyr"
		}
	}
	return mcp.NewToolResultText(content), nil

}

func detectTheRealTopicInUserMessageHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	fmt.Println("🛠️  detectTheRealTopicInUserMessageHandler called with args:", args)
	// Check if the topic_name argument is provided
	if len(args) == 0 {
		return mcp.NewToolResultText("zephyr"), nil
	}
	var content = "zephyr" //default character
	if topicName, ok := args["topic_name"]; ok {
		switch strings.ToLower(topicName.(string)) {
		case "food", "gold":
			content = "galdor"
		case "war", "weapons", "fight":
			content = "thrain"
		case "magic", "poetry", "art":
			content = "elara"
		case "medicine", "health", "potion":
			content = "liora"
		default:
			content = "zephyr"
		}
	}
	return mcp.NewToolResultText(content), nil
}
