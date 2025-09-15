# Customize `models` runtime

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
- line 29 to 34: customization of model runtime

### Source code
- [sorcerer.agent.go](agents/sorcerer.agent.go) we removed the temperature definition `	temperature := helpers.StringToFloat(helpers.GetEnvOrDefault("SORCERER_MODEL_TEMPERATURE", "0.0"))`
