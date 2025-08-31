package tools

import (
	"context"
	"dungeon-mcp-server/types"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

func CollectGoldTool() mcp.Tool {
	return mcp.NewTool("collect_gold",
		mcp.WithDescription(`Collect gold coins from the current room if available. Try: "Collect the gold coins"`),
	)
}

func CollectGoldToolHandler(player *types.Player, dungeon *types.Dungeon) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		
		if player.Name == "Unknown" {
			message := "✋ No player exists. Please create a player first."
			fmt.Println(message)
			return mcp.NewToolResultText(message), fmt.Errorf("no player exists")
		}

		var currentRoom *types.Room
		for i := range dungeon.Rooms {
			if dungeon.Rooms[i].ID == player.RoomID {
				currentRoom = &dungeon.Rooms[i]
				break
			}
		}

		if currentRoom == nil {
			message := "❌ Player is not in any room."
			fmt.Println(message)
			return mcp.NewToolResultText(message), fmt.Errorf("player not in any room")
		}

		if currentRoom.GoldCoins <= 0 {
			message := fmt.Sprintf("💰 There are no gold coins to collect in %s.", currentRoom.Name)
			fmt.Println(message)
			return mcp.NewToolResultText(message), nil
		}

		collectedGold := currentRoom.GoldCoins
		player.GoldCoins += collectedGold
		currentRoom.GoldCoins = 0

		message := fmt.Sprintf("💰 You collected %d gold coins from %s! Your total gold coins: %d", 
			collectedGold, currentRoom.Name, player.GoldCoins)
		fmt.Println(message)
		return mcp.NewToolResultText(message), nil
	}
}