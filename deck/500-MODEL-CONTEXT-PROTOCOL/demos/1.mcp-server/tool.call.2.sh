#!/bin/bash
: <<'COMMENT'
# Use tool "add"
COMMENT

# STEP 1: Load the session ID from the environment file
source mcp.session.env
source mcp.server.env

MCP_SERVER=${MCP_SERVER:-"http://localhost:${MCP_HTTP_PORT}"}

read -r -d '' DATA <<- EOM
{
  "jsonrpc": "2.0",
  "id": "test",
  "method": "tools/call",
  "params": {
    "name": "detect_real_topic_in_user_message",
    "arguments": {
      "topic_name": "magic"
    }
  }
}
EOM

curl ${MCP_SERVER}/mcp \
  -H "Content-Type: application/json" \
  -H "Mcp-Session-Id: $SESSION_ID" \
  -d "${DATA}" | jq 


