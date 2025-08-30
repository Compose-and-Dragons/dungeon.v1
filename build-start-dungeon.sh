#!/bin/bash
echo "🐳 Starting the Compose and Dragons Dungeon MCP server and Gateway services..."
echo "docker compose up --build -d"
docker compose up --build -d
echo "docker compose logs -f mcp-dungeon mcp-gateway"
docker compose logs -f mcp-dungeon mcp-gateway