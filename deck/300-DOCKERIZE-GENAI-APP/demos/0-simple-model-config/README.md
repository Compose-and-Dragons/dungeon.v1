# Simplest `models` configuration within Compose

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
- line 19 to 21: declaration of the models that will be attached to service containers
- line 23 to 27: definition of the models to use in this Compose application

### Source code

- [main.go](main.go) line 20 - `llmURL := helpers.GetEnvOrDefault("SORCERER_URL", "http://localhost:12434/engines/llama.cpp/v1")`
- [store.go](agents/stores.go) line 53 - `Model: helpers.GetEnvOrDefault("EMBEDDING_MODEL", "ai/mxbai-embed-large:latest"),`
- [sorcerer.agent.go](agents/sorcerer.agent.go) line 31 - `model := helpers.GetEnvOrDefault("SORCERER_MODEL", "ai/qwen2.5:1.5B-F16")`

### Notes
> The configuration which was done in the shell scripts previously is now directly done in the Compose file