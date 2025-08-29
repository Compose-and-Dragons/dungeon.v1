#!/bin/bash
echo "🐳 Starting the Compose and Dragons Dungeon Master..."
echo "docker attach $(docker compose ps -q dungeon-master)"
docker attach $(docker compose ps -q dungeon-master)