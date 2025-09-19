---
marp: true
theme: default
paginate: true
---
# 🏰 Dungeon crawler **MCP server**
> Streamable HTTP server
---
## Initialize: `func main()`

- [Create `MCP server`: /dungeon-crawler-mcp-server/main.go#L26](/dungeon-crawler-mcp-server/main.go#L26)
- [Create `OpenAI client`: /dungeon-crawler-mcp-server/main.go#L43](/dungeon-crawler-mcp-server/main.go#L43)
- [`Agent` Creation: /dungeon-crawler-mcp-server/main.go#L51](/dungeon-crawler-mcp-server/main.go#L51)
  - Role: <!-- TODO: -->
  - Use of JSON output format (the LLM will "answer in JSON")
    `data.GetRoomSchema()` [Room schema: /dungeon-crawler-mcp-server/data/response.schema.go#L5](/dungeon-crawler-mcp-server/data/response.schema.go#L5)

---
## Initialize: `func main()`

- [Initialize the `current player`: /dungeon-crawler-mcp-server/main.go#L72](/dungeon-crawler-mcp-server/main.go#L72)
  - Empty player struct: `types.Player{}`
  - CurrentRoomID: will be set to the entrance room ID when the entrance room is created
- [Initialize the `Dungeon struct`: /dungeon-crawler-mcp-server/main.go#L92](/dungeon-crawler-mcp-server/main.go#L92)
  - Name, Description, Size, Entrance/Exit Coordinates
  - Rooms: empty array `[]types.Room{}`
---
## Initialize: `func main()`

- [Generate the `entrance room` with the `dungeon agent`: /dungeon-crawler-mcp-server/main.go#L116](/dungeon-crawler-mcp-server/main.go#L116)
  - create `types.Room{ ID: "room_0_0"}`
  - Add the room to the Dungeon: `dungeon.Rooms = append(dungeon.Rooms, entranceRoom)`

---
## Initialize: `func main()`

### Tools Section
- [Define the `MCP tools`: /dungeon-crawler-mcp-server/main.go#L166](/dungeon-crawler-mcp-server/main.go#L166)
  - 🤚 **These tools will be used by the dungeon-master program**
  - Example: `createPlayerToolInstance := tools.CreatePlayerTool()`
    - Define the **tool**: name, description, parameters (if any)
    - Define the **handler function**: what the tool does when called
    - **Register** the tools with the MCP server: `s.AddTool(createPlayerToolInstance, tools.CreatePlayerToolHandler(&currentPlayer, &dungeon))`

---
## Initialize: `func main()`

### Tools List 1/2
- `create_player`: `tools.CreatePlayerTool()`
- `get_player_info`: `tools.GetPlayerInformationTool()`
- `get_dungeon_info`: `tools.GetDungeonInformationTool()`
- `move_by_direction`: `tools.GetMoveIntoTheDungeonTool()` -> **`MoveByDirectionToolHandler`**
- `move_player`: `tools.GetMovePlayerTool()` -> **`MoveByDirectionToolHandler`**

---
## Initialize: `func main()`

### Tools List 2/2
- `get_current_room_info`: `tools.tools.GetCurrentRoomInformationTool()`
- `get_dungeon_map`: `tools.GetDungeonMapTool()` -> **`GetDungeonMapToolHandler`**
- `collect_gold`: `tools.CollectGoldTool()`
- `collect_magic_potion`: `tools.CollectMagicPotionTool()`
- `fight_monster`: `tools.FightMonsterTool()` -> **`FightMonsterToolHandler`**
- `is_player_in_same_room_as_npc`: `tools.IsPlayerInSameRoomAsNPCTool()`


---
## Start: `func main()`

- [Start the `MCP server`: /dungeon-crawler-mcp-server/main.go#L216](/dungeon-crawler-mcp-server/main.go#L216)
  - And add a health check endpoint `/health` (for the usage with Docker MCP Gateway and Docker Compose)

---
## Handlers: `MoveByDirectionToolHandler`



---
## Handlers: `GetDungeonMapToolHandler`


---
## Handlers: `FightMonsterToolHandler`

