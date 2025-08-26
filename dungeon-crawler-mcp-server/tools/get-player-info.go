package tools

import (
	"context"
	"dungeon-mcp-server/types"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

func GetPlayerInformationTool() mcp.Tool {
	return mcp.NewTool("get_player_info",
		mcp.WithDescription(`Get the current player's information. Try: "Who am I?"`),
	)
}

func GetPlayerInformationToolHandler(player *types.Player, dungeon *types.Dungeon) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if player.Name == "Unknown" {
			message := "✋ No player exists. Please create a player first."
			fmt.Println(message)
			return mcp.NewToolResultText(message), fmt.Errorf("no player exists")
		}

		playerJSON, err := json.MarshalIndent(*player, "", "  ")
		if err != nil {
			return nil, err
		}

		return mcp.NewToolResultText(string(playerJSON)), nil
	}
}
