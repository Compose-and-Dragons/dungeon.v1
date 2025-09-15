# Customize `models` env variables passed to service containers

To run the demo
```shell
$ docker compose up --build -d
$ docker compose attach sorcerer-agent
🔶 Loading vector store from: /app/data/elara_vector_store.json
✅ Vector store loaded successfully with 19 records
┃ 🤖 (/bye to exit) [Elara]>                                                                                                                                                                                                                                                  
┃ Type your command here...                                                                                                                                                                                                                                                   
┃                                                                                                                                                                                                                                                                             
┃                                                                                                                                                                                                                                                                             
┃                                                                                                                                                                                                                                                                             
┃                                                                                                                                                                                                                                                                             
┃                                                                                                                                                                                                                                                                             
                                                                                                                                                                                                                                                                              
alt+enter / ctrl+j new line • enter submit                                                                                                                                                                                                                                    
```

## Keypoints to check
### Compose file 
[Compose file](compose.yaml):
- line 19 to 24: customization of models' variables that will be sent to service containers

### Source code

- [main.go](main.go) line 20 - `llmURL := helpers.GetEnvOrDefault("MODEL_RUNNER_BASE_URL", "http://localhost:12434/engines/llama.cpp/v1")`
- [sorcerer.agent.go](agents/sorcerer.agent.go) line 31 - `model := helpers.GetEnvOrDefault("ELARA_MODEL", "ai/qwen2.5:1.5B-F16")`
