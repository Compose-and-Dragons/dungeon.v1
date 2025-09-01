#!/bin/bash
echo "🐳 Starting the Compose and Dragons Dungeon MCP server and Gateway services..."
echo "docker compose up -d"
docker compose up -d
echo "docker compose logs -f mcp-dungeon mcp-gateway dungeon-end-of-level-boss"
docker compose logs -f mcp-dungeon mcp-gateway dungeon-end-of-level-boss