package tools

import (
	"context"
	"dungeon-mcp-server/helpers"
	"dungeon-mcp-server/types"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/micro-agent/micro-agent-go/agent/mu"
	"github.com/openai/openai-go/v2"
)

func GetMoveIntoTheDungeonTool() mcp.Tool {

	moveByDirection := mcp.NewTool("move_by_direction",
		mcp.WithDescription(`Move the player in a specified direction (north, south, east, west). Try "move by north".`),
		mcp.WithString("direction",
			mcp.Required(),
			mcp.Description("The direction to move in. Must be one of: north, south, east, west"),
		),
	)
	return moveByDirection

}

// TODO:
/*
- generate room name and description with a model
- add monsters, items, and traps to rooms
- handle special rooms (entrance, exit)

*/

func MoveByDirectionToolHandler(player *types.Player, dungeon *types.Dungeon, dungeonAgent *mu.Agent) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if player.Name == "Unknown" {
			message := "✋ No player exists. Please create a player first."
			fmt.Println(message)
			return mcp.NewToolResultText(message), fmt.Errorf("no player exists")
		}

		args := request.GetArguments()
		direction := args["direction"].(string)

		fmt.Println("➡️ Move by direction:", direction)

		newX := player.Position.X
		newY := player.Position.Y

		switch direction {
		case "north":
			newY++
		case "south":
			newY--
		case "east":
			newX++
		case "west":
			newX--
		default:
			message := fmt.Sprintf("❌ Invalid direction: %s. Must be one of: north, south, east, west", direction)
			fmt.Println(message)
			return mcp.NewToolResultText(message), fmt.Errorf("invalid direction")
		}

		if newX < 0 || newX >= dungeon.Width || newY < 0 || newY >= dungeon.Height {
			message := fmt.Sprintf("❌ Cannot move %s from (%d, %d). Position (%d, %d) is outside the dungeon boundaries (0--%d, 0--%d).",
				direction, player.Position.X, player.Position.Y, newX, newY, dungeon.Width-1, dungeon.Height-1)
			fmt.Println(message)
			return mcp.NewToolResultText(message), fmt.Errorf("position outside dungeon boundaries")
		}

		player.Position.X = newX
		player.Position.Y = newY

		roomID := fmt.Sprintf("room_%d_%d", newX, newY)
		var currentRoom *types.Room

		for i := range dungeon.Rooms {
			if dungeon.Rooms[i].ID == roomID {
				currentRoom = &dungeon.Rooms[i]
				break
			}
		}

		if currentRoom == nil {
			// NOTE: it's a new room, create it => generate room name and description with a model

			// ---------------------------------------------------------
			// BEGIN: Generate the room with the dungeon agent
			// ---------------------------------------------------------
			dungeonAgentSystemInstruction := helpers.GetEnvOrDefault("DUNGEON_AGENT_SYSTEM_INSTRUCTION", "You are a Dungeon Master. You create rooms in a dungeon. Each room has a name and a short description.")

			response, err := dungeonAgent.Run([]openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(dungeonAgentSystemInstruction),
				openai.UserMessage("Create a new dungeon room with a name and a short description."),
			})

			if err != nil {
				fmt.Println("🔴 Error generating room:", err)
				// TODO: handle error better => generate a default room
				return mcp.NewToolResultText(""), err

			}
			fmt.Println("📝 Dungeon Room Response:", response)

			var roomResponse struct {
				Name        string `json:"name"`
				Description string `json:"description"`
			}
			err = json.Unmarshal([]byte(response), &roomResponse)
			if err != nil {
				fmt.Println("🔴 Error unmarshaling room response:", err)
				// TODO: handle error better => generate a default room
				return mcp.NewToolResultText(""), err
			}
			fmt.Println("👋🏰 Room:", roomResponse)

			// ---------------------------------------------------------
			// END: of Generate the room with the dungeon agent
			// ---------------------------------------------------------

			// newRoom := types.Room{
			// 	ID:          roomID,
			// 	Name:        fmt.Sprintf("Room at (%d, %d)", newX, newY),
			// 	Description: fmt.Sprintf("You are in a room at coordinates (%d, %d).", newX, newY),
			// 	Coordinates: types.Coordinates{X: newX, Y: newY},
			// 	Visited:     true,
			// 	IsEntrance:  newX == dungeon.EntranceCoords.X && newY == dungeon.EntranceCoords.Y,
			// 	IsExit:      newX == dungeon.ExitCoords.X && newY == dungeon.ExitCoords.Y,
			// }

			newRoom := types.Room{
				ID:          roomID,
				Name:        roomResponse.Name,
				Description: roomResponse.Description,
				Coordinates: types.Coordinates{X: newX, Y: newY},
				Visited:     true,
				IsEntrance:  newX == dungeon.EntranceCoords.X && newY == dungeon.EntranceCoords.Y,
				IsExit:      newX == dungeon.ExitCoords.X && newY == dungeon.ExitCoords.Y,
			}

			dungeon.Rooms = append(dungeon.Rooms, newRoom)
			currentRoom = &dungeon.Rooms[len(dungeon.Rooms)-1]
		} else {
			currentRoom.Visited = true
		}

		resultMessage := fmt.Sprintf("✅ Moved %s to position (%d, %d)\n\n🏠 %s\n%s",
			direction, newX, newY, currentRoom.Name, currentRoom.Description)

		fmt.Println(resultMessage)
		return mcp.NewToolResultText(resultMessage), nil
	}
}
