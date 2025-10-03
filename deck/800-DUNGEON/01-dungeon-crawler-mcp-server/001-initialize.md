---
marp: true
theme: default
paginate: true
---
# 🏰 Dungeon crawler **MCP server**
> Streamable HTTP server
---
### Tools List 1/2
- `create_player`: -> **`CreatePlayerToolHandler`**
- `get_player_info`
- `get_dungeon_info`
- `move_by_direction`: -> **`MoveByDirectionToolHandler`**
- `move_player`: -> **`MoveByDirectionToolHandler`**

---
### Tools List 2/2
- `get_current_room_info`
- `get_dungeon_map`: -> **`GetDungeonMapToolHandler`**
- `collect_gold`
- `collect_magic_potion`
- `fight_monster`: -> **`FightMonsterToolHandler`**
- `is_player_in_same_room_as_npc`

---
## The most important Handlers: 

- `CreatePlayerToolHandler`
- `MoveByDirectionToolHandler`
- `GetDungeonMapToolHandler`
- `FightMonsterToolHandler`

[← Previous](../001-initialize.md) | [Next →](../02-dungeon-master/001-initialize.md)
