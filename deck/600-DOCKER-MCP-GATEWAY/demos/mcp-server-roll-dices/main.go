package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"math/rand"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {

	// Create MCP server
	s := server.NewMCPServer(
		"mcp-roll-dices",
		"0.0.0",
	)

	rollDices := mcp.NewTool("roll_dice",
		mcp.WithDescription("Roll dice to get a random result."),
		mcp.WithNumber("nb_dices",
			mcp.Required(),
			mcp.Description("The number of dice to roll"),
		),
		mcp.WithNumber("nb_sides",
			mcp.Required(),
			mcp.Description("The number of faces on the dice"),
		),
	)

	s.AddTool(rollDices, rollDicesHandler)

	// ---------------------------------------------------------
	// Start the HTTP server
	// ---------------------------------------------------------
	httpPort := "9090"
	fmt.Println("🌍 MCP HTTP Port:", httpPort)

	log.Println("MCP StreamableHTTP server is running on port", httpPort)

	// Create a custom mux to handle both MCP and health endpoints
	mux := http.NewServeMux()

	// Add healthcheck endpoint
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

func rollDicesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

	nbDices := request.GetInt("nb_dices", 1)
	sides := request.GetInt("nb_sides", 6)

	log.Printf("🎲 Rolling %d dice(s) with %d sides each...\n", nbDices, sides)

	roll := func(n, x int) int {
		if n <= 0 || x <= 0 {
			return 0
		}

		results := make([]int, n)
		sum := 0

		for i := range n {
			roll := rand.Intn(x) + 1 // +1 because rand.Intn(x) gives 0 to x-1
			results[i] = roll
			sum += roll
		}

		return sum
	}

	// Simulate rolling dice
	result := roll(nbDices, sides)

	return mcp.NewToolResultText("The result of rolling " + strconv.Itoa(nbDices) + " dice with " + strconv.Itoa(sides) + " faces is: " + strconv.Itoa(result)), nil

}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	response := map[string]any{
		"status": "healthy",
	}
	json.NewEncoder(w).Encode(response)
}
