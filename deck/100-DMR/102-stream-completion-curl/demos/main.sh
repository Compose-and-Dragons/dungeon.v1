#!/bin/bash
BASE_URL=${MODEL_RUNNER_BASE_URL:-http://localhost:12434/engines/llama.cpp/v1}
MODEL=${MODEL_QWEN2_5_LARGE:-"ai/qwen2.5:latest"}

read -r -d '' SYSTEM_INSTRUCTION <<- EOM
You are an expert of medieval role playing games. 
EOM

read -r -d '' USER_CONTENT <<- EOM
[Brief] What is a dungeon crawler game?
EOM

read -r -d '' DATA <<- EOM
{
  "model":"${MODEL}",
  "options": {
    "temperature": 0.5,
    "repeat_last_n": 2
  },
  "messages": [
    {"role":"system", "content": "${SYSTEM_INSTRUCTION}"},
    {"role":"user", "content": "${USER_CONTENT}"}
  ],
  "stream": true
}
EOM

# Remove newlines from DATA 
DATA=$(echo ${DATA} | tr -d '\n')

echo "Using DATA: ${DATA}"

# JSON_RESULT=$(curl --silent ${BASE_URL}/chat/completions \
#     -H "Content-Type: application/json" \
#     -d "${DATA}"
# )
unescape_quotes() {
    local str="$1"
    str="${str//\\\"/\"}"  # Replace \" by "
    echo "$str"
}

remove_quotes() {
    local str="$1"
    str="${str%\"}"   # remove " at the end
    str="${str#\"}"   # remove " at start
    echo "$str"
}

callback() {
  echo -ne "$1" 
}


curl --no-buffer --silent ${BASE_URL}/chat/completions \
    -H "Content-Type: application/json" \
    -d "${DATA}" \
    | while IFS= read -r line; do
        #echo "📝 $line"
        if [[ $line == data:* ]]; then
        
            json_data="${line#data: }"
            if [[ $json_data != "[DONE]" ]]; then
                # Extract content if it exists, else return "null"
                content_chunk=$(echo "$json_data" | jq '.choices[0].delta.content // "null"' 2>/dev/null)

                if [[ "$content_chunk" != "\"null\"" ]]; then
                    result=$(remove_quotes "$content_chunk")
                    clean_result=$(unescape_quotes "$result")
                    callback "$clean_result"
                fi

            fi
        fi
    done        


echo -e "\n"


