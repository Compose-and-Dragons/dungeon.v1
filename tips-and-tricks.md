# Tips and Tricks

## Player-NPC Interaction Validation

The verification of whether a player can talk to an NPC (if they are in the same room) occurs in the code at the following location:

**Main validation**: `dungeon-master/main.go:37` in the `speak_to_somebody` function. The system uses the MCP tool `is_player_in_same_room_as_npc` to validate that the player and NPC are in the same room before allowing the interaction.

**Complete validation logic**: The full validation implementation can be found in `dungeon-crawler-mcp-server/tools/is-player-in-same-room-as-npc.go:30-98`.