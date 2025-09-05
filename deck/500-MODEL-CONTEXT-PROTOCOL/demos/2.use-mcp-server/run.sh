#!/bin/bash
export MODEL_RUNNER_BASE_URL="http://localhost:12434/engines/llama.cpp/v1"
export MCP_HTTP_SERVER_URL="http://localhost:9090/mcp"

go run main.go