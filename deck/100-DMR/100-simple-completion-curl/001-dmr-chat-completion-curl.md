---
marp: true
theme: default
paginate: true
---
# Chat Completion with Docker Model Runner & **`Curl`**

---
## System Instruction & User Content

```bash
read -r -d '' SYSTEM_INSTRUCTION <<- EOM
You are an expert of medieval role playing games. 
EOM

read -r -d '' USER_CONTENT <<- EOM
[Brief] What is a dungeon crawler game?
EOM
```
---
## Request Data

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
  "stream": false
}
EOM
```
---
## Curl Command

```bash
# Remove newlines from DATA 
DATA=$(echo ${DATA} | tr -d '\n')

JSON_RESULT=$(curl --silent ${BASE_URL}/chat/completions \
    -H "Content-Type: application/json" \
    -d "${DATA}"
)

CONTENT=$(echo "${JSON_RESULT}" | jq -r '.choices[0].message.content')
echo "${CONTENT}"
```
