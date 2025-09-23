# Run Zephyr Agent and give it access to MCP tools via MCP Gateway

To run the demo
```shell
$ docker compose up --build -d
$ docker compose attach zephyr
MCP Client initialized successfully
Tool: choose_character_by_species - select a species from among these: [Human, Orc, Elf, Dwarf] by saying: I want to talk to a <species_name>.
Tool: detect_real_topic_in_user_message - select a topic from among these: [justice, war, combat, magic, poetry, craftsmanship, forge] by saying: I have a question about <topic_name>.                                                                                                                                                                                   
Tool: roll_dice - Roll dice to get a random result.

┃ 🤖 (/bye to exit) [Zephyr]>                                                                                                                                                        
┃ Type your command here...                                                                                                                                                          
┃                                                                                                                                                                                    
┃                                                                                                                                                                                    
┃                                                                                                                                                                                    
┃                                                                                                                                                                                    
┃                                                                                                                                                                                    
                                                                                                                                                                                     
alt+enter / ctrl+j new line • enter submit                                                                                                                                                                                                                                      
```