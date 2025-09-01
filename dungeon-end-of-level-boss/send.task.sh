#!/bin/bash
: <<'COMMENT'
# Send task to the agent server
COMMENT

HTTP_PORT=8888
#AGENT_BASE_URL=http://0.0.0.0:${HTTP_PORT}
AGENT_BASE_URL=http://localhost:${HTTP_PORT}

# host.docker.internal

read -r -d '' DATA <<- EOM
{
    "jsonrpc": "2.0",
    "id": "1111",
    "method": "message/send",
    "params": {
      "message": {
        "role": "user",
        "parts": [
          {
            "text": "Tell me something about your family?"
          }
        ]
      },
      "metadata": {
        "skill": "ask_for_something"
      }
    }
}
EOM

function on_chunk() {
    chunk=$1
    
    # Vérifie si la ligne commence par "data: "
    if [[ ${chunk} == data:* ]]; then
        # Extrait le JSON après "data: "
        json_data=${chunk#data: }
        
        # Vérifie si c'est un chunk de streaming
        if echo "${json_data}" | jq -e '.type == "chunk"' >/dev/null 2>&1; then
            content=$(echo "${json_data}" | jq -r '.content')
            echo -n "${content}"
        elif echo "${json_data}" | jq -e '.result.history[0].parts[0].text' >/dev/null 2>&1; then
            # Réponse finale complète
            echo ""
            echo "✅ Final response received"
        fi
    fi
}

function send_message() {
    AGENT_BASE_URL="${1}/stream"
    DATA="${2}"
    CALL_BACK=${3}

    echo "🌊 Starting streaming request..."
    curl --no-buffer --silent ${AGENT_BASE_URL} \
        -H "Content-Type: application/json" \
        -H "Accept: text/event-stream" \
        -d "${DATA}" | while read linestream
        do
            ${CALL_BACK} "${linestream}"
        done 
    echo ""
    echo "🏁 Stream completed"
}

send_message "${AGENT_BASE_URL}" "${DATA}" on_chunk



