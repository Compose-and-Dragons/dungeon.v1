package tools

import (
	"context"
	"dungeon-mcp-server/types"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

func CollectMagicPotionTool() mcp.Tool {
	return mcp.NewTool("collect_magic_potion",
		mcp.WithDescription(`Collect magic potions from the current room if available. Try: "Collect the magic potions"`),
	)
}

func CollectMagicPotionToolHandler(player *types.Player, dungeon *types.Dungeon) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

		if !currentRoom.HasMagicPotion || currentRoom.RegenerationHealth <= 0 {
			message := fmt.Sprintf("🧪 There are no magic potions to collect in %s.", currentRoom.Name)
			fmt.Println(message)
			return mcp.NewToolResultText(message), nil
		}

		collectedPotion := currentRoom.RegenerationHealth
		player.Health += collectedPotion
		currentRoom.HasMagicPotion = false
		currentRoom.RegenerationHealth = 0

		message := fmt.Sprintf("🧪 You collected a magic potion from %s! You gained %d health points. Your current health: %d", 
			currentRoom.Name, collectedPotion, player.Health)
		fmt.Println(message)
		return mcp.NewToolResultText(message), nil
	}
}