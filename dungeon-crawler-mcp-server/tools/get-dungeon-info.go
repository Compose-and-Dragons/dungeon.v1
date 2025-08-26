package tools

import (
	"context"
	"dungeon-mcp-server/types"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
)

func GetDungeonInformationTool() mcp.Tool {
	return mcp.NewTool("get_dungeon_info",
		mcp.WithDescription(`Get the current dungeon's information including its layout, rooms, entrance and exit coordinates.`),
	)
}

func GetDungeonInformationToolHandler(player *types.Player, dungeon *types.Dungeon) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		dungeonJSON, err := json.MarshalIndent(*dungeon, "", "  ")
		if err != nil {
			return nil, err
		}

		return mcp.NewToolResultText(string(dungeonJSON)), nil
	}
}
