---
marp: true
theme: default
paginate: true
---
# Chat Stream Completion with Docker Model Runner & **`Curl`**

---
## **`"stream": true`**

```bash
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
  "stream": true # ⬅️ 👋👋👋
}
EOM
```

---
## **`callback()`** function

```bash
callback() {
  echo -ne "$1" 
}
```

---
## **`curl`** request with streaming

```bash
curl --no-buffer --silent ${BASE_URL}/chat/completions \
    -H "Content-Type: application/json" \
    -d "${DATA}" \
    | while IFS= read -r line; do
        if [[ $line == data:* ]]; then
          # Extract JSON data after "data: ", 
          content_chunk=$(echo "$json_data" | jq '.choices[0].delta.content // "null"' 2>/dev/null)
          # Clean content_chunk and call the callback
          # ...
          callback "$clean_result"
        fi
    done      
```

[← Previous](../101-simple-completion-go/101-dmr-chat-completion-openai-go.md) | [Next →](../103-stream-completion-go/000-dmr-chat-stream-completion-openai-go.md)
