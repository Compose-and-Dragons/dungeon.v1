package tools

import (
	"context"
	"dungeon-mcp-server/helpers"
	"dungeon-mcp-server/types"
	"encoding/json"
	"fmt"
	"math/rand"

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

func GetMovePlayerTool() mcp.Tool {

	movePlayer := mcp.NewTool("move_player",
		mcp.WithDescription(`Move the player in the dungeon by specifying a cardinal direction. This is the primary navigation tool for exploring rooms. Usage: "move player north" or "go east".`),
		mcp.WithString("direction",
			mcp.Required(),
			mcp.Description("Cardinal direction to move the player. MUST be exactly one of these values: 'north', 'south', 'east', 'west' (lowercase only)"),
		),
	)
	return movePlayer

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

		// Update player's current room ID => useful to get current room info
		player.RoomID = fmt.Sprintf("room_%d_%d", newX, newY)

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
			dungeonAgentRoomSystemInstruction := helpers.GetEnvOrDefault("DUNGEON_AGENT_ROOM_SYSTEM_INSTRUCTION", "You are a Dungeon Master. You create rooms in a dungeon. Each room has a name and a short description.")

			response, err := dungeonAgent.Run([]openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(dungeonAgentRoomSystemInstruction),
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

			// Add NPCs, monsters, and items based on probabilities of appearance

			// ---------------------------------------------------------
			// BEGIN: Create NPC 🧙‍♂️
			// ---------------------------------------------------------
			// QUESTION: comment un tool peut déclencher un autre tool ?
			// NOTE: sinon cela dit qu'il y a un NPC dans la pièce
			// et le joueur peut utiliser le tool "talk_to_npc" pour interagir avec lui

			var hasNonPlayerCharacter bool
			var nonPlayerCharacter types.NonPlayerCharacter

			merchantRoom := helpers.GetEnvOrDefault("MERCHANT_ROOM", "room_1_1")
			guardRoom := helpers.GetEnvOrDefault("GUARD_ROOM", "room_0_2")
			sorcererRoom := helpers.GetEnvOrDefault("SORCERER_ROOM", "room_2_0")
			healerRoom := helpers.GetEnvOrDefault("HEALER_ROOM", "room_2_2")

			switch roomID {
			case merchantRoom:
				hasNonPlayerCharacter = true
				nonPlayerCharacter = types.NonPlayerCharacter{
					Type:     types.Merchant,
					Name:     helpers.GetEnvOrDefault("MERCHANT_NAME", "[default]Gorim the Merchant"),
					Race:     helpers.GetEnvOrDefault("MERCHANT_RACE", "Dwarf"),
					Position: types.Coordinates{X: newX, Y: newY},
					RoomID:   roomID,
				}

			case guardRoom:
				hasNonPlayerCharacter = true
				nonPlayerCharacter = types.NonPlayerCharacter{
					Type:     types.Guard,
					Name:     helpers.GetEnvOrDefault("GUARD_NAME", "[default]Lyria the Guard"),
					Race:     helpers.GetEnvOrDefault("GUARD_RACE", "Elf"),
					Position: types.Coordinates{X: newX, Y: newY},
					RoomID:   roomID,
				}

			case sorcererRoom:
				hasNonPlayerCharacter = true
				nonPlayerCharacter = types.NonPlayerCharacter{
					Type:     types.Sorcerer,
					Name:     helpers.GetEnvOrDefault("SORCERER_NAME", "[default]Eldrin the Sorcerer"),
					Race:     helpers.GetEnvOrDefault("SORCERER_RACE", "Human"),					
					Position: types.Coordinates{X: newX, Y: newY},
					RoomID:   roomID,
				}

			case healerRoom:
				hasNonPlayerCharacter = true
				nonPlayerCharacter = types.NonPlayerCharacter{
					Type:     types.Healer,
					Name:     helpers.GetEnvOrDefault("HEALER_NAME", "[default]Mira the Healer"),
					Race:     helpers.GetEnvOrDefault("HEALER_RACE", "Half-Elf"),					
					Position: types.Coordinates{X: newX, Y: newY},
					RoomID:   roomID,
				}

			default:
				hasNonPlayerCharacter = false
				nonPlayerCharacter = types.NonPlayerCharacter{}
			}

			// ---------------------------------------------------------
			// END: Create NPC
			// ---------------------------------------------------------

			// ---------------------------------------------------------
			// BEGIN: Create Monster 👹
			// ---------------------------------------------------------
			monsterProbability := helpers.StringToFloat(helpers.GetEnvOrDefault("MONSTER_PROBABILITY", "0.25"))

			// NOTE: 

			// 100 x monsterProbability % of chance to have a monster in the room
			// except if there is already a NPC in the room
			if rand.Float64() < monsterProbability && !hasNonPlayerCharacter {
				// TODO: create a monster with a model
			}

			// ---------------------------------------------------------
			// END: Create Monster
			// ---------------------------------------------------------

			// ---------------------------------------------------------
			// BEGIN: Create Gold coins, potions, and items ⭐️
			// ---------------------------------------------------------
			itemProbability := helpers.StringToFloat(helpers.GetEnvOrDefault("ITEM_PROBABILITY", "0.20"))

			// 100 x itemProbability % of chance to have an item in the room
			var hasTreasure, hasMagicPotion bool
			var regenerationHealth, goldCoins int
			if rand.Float64() < itemProbability {
				if rand.Float64() < 0.5 {
					hasTreasure = true
					goldCoins = rand.Intn(50) + 10 // between 10 and 59 gold coins
				} else {
					hasMagicPotion = true
					regenerationHealth = rand.Intn(20) + 5 // between 5 and 24 health points
				}
			}

			// ---------------------------------------------------------
			// END: Create Gold coins, potions, and items
			// ---------------------------------------------------------
			newRoom := types.Room{
				ID:                    roomID,
				Name:                  roomResponse.Name,
				Description:           roomResponse.Description,
				Coordinates:           types.Coordinates{X: newX, Y: newY},
				Visited:               true,
				IsEntrance:            newX == dungeon.EntranceCoords.X && newY == dungeon.EntranceCoords.Y,
				IsExit:                newX == dungeon.ExitCoords.X && newY == dungeon.ExitCoords.Y,
				HasTreasure:           hasTreasure,
				GoldCoins:             goldCoins,
				HasMagicPotion:        hasMagicPotion,
				RegenerationHealth:    regenerationHealth,
				HasNonPlayerCharacter: hasNonPlayerCharacter,
				NonPlayerCharacter:    &nonPlayerCharacter,
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
