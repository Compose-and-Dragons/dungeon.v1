# Compose and Dragons

## VSCode Setup

The `.vscode` folder is **mandatory**. Certain settings and extensions are used to improve document and source code display and facilitate presentation.

## Launching the Docker Compose Project

At the project root, run the following command to start the Docker containers:
```bash
docker compose up --build -d
docker compose logs -f mcp-dungeon mcp-gateway dungeon-end-of-level-boss
```
> Or use `build-start-dungeon.sh`

### What is being launched?

```bash
docker compose ps --services
```

- `mcp-dungeon` (project: `dungeon-crawler-mcp-server`): the custom MCP server for the dungeon game. (transport used: Streamable HTTP)
  - The server is started automatically
- `mcp-gateway`: the MCP Gateway server that allows connection to the `mcp-dungeon` MCP server and potentially other MCP servers.
  - The gateway is started automatically
- `dungeon-end-of-level-boss` (project: `dungeon-end-of-level-boss`): a **"functional"** "independent" AI agent that uses 2 technical agents (1 for chat completion, the other for RAG).
  - This is an "End-of-level Boss" NPC, named **"Shesepankh"** who will be used at the dungeon exit.
  - The agent configuration data is in the `dungeon-end-of-level-boss/data` folder
  - The agent is started automatically: on first startup, the agent will build the vector store in a JSON file from the configuration data. This may take a moment. It will then be instantaneous on subsequent startups.
- `dungeon-master` (project: `dungeon-master`): it is composed of **multiple** **"functional"** AI agents:
  - **"Zephyr"**: can be considered the "dungeon master" who manages the game state and interacts with the player.
    - The container is started but the `dungeon-master` program is **not** launched automatically.
    - It will detect player commands and call other **"functional"** agents if necessary.
    - It will use MCP Tools from the `mcp-dungeon` MCP server to respond to player actions.
    - It uses 1 technical agent (1 for chat completion & tool call detection).
  - The other **"functional"** agents (NPCs) are invoked by "Zephyr" **on demand** from the player.
    - **"Galdor"**: an NPC Merchant
    - **"Elara"**: an NPC Sorceress
    - **"Thrain"**: an NPC Guardian
    - **"Liora"**: an NPC Healer
    - Each of them uses 2 technical agents (1 for chat completion, the other for RAG).
    - ✋👻 *there is also a "Ghost agent" which is a fake AI agent for testing purposes - no utility for gameplay*.

### Launching the "Zephyr" agent (the Dungeon Master with user interface)

To launch the user interface, execute the following command in a new terminal:
```bash
docker attach $(docker compose ps -q dungeon-master)
```
On first startup, the 4 NPC agents will each build a vector store in a JSON file from configuration data (present in the `dungeon-master/data` folder). This may take a moment. It will then be instantaneous on subsequent startups.

## Game Operation Principle

- The dungeon is extremely simple: a 4x4 square. (configurable in the `compose.yml` file)
- The 4 NPCs as well as the end-level Boss are placed in advance in specific rooms.
  - This position configuration is done in the `compose.yml` file
- We use the "Dungeon Crawler" principle (dungeon exploration) with dynamically generated rooms:
  - Once a room is generated, it is stored in the `mcp-dungeon` MCP server (tool: `store_room`) and will no longer be modified (the player can return to the room).
  - When generating a room, we can also generate:
    - objects (potions, gold coins) which are stored in the `mcp-dungeon` MCP server.
    - enemies (monsters) which are also stored in the `mcp-dungeon` MCP server.
    - The player can interact with objects and enemies (combat, collect objects, etc.)
- The player can move in 4 directions (north, south, east, west) and interact with NPCs.
- The goal of the game is to find and defeat the end-level Boss (**Shesepankh**) who is located in a specific room.
- **[🚧 not yet implemented (1)]** The player will collect information from NPCs that will be useful during their discussion with **Shesepankh**. The player will therefore need to meet all NPCs and ask them the right questions.
- **[🚧 not yet implemented (2)]** The player will need to give the collected information to **Shesepankh** for her to agree to let them exit.

> - (1): this should be possible with prompts only.
> - (2): this should be possible with prompts + MCP tool(s).

## Game Flow

1. Launch the Docker Compose project (see above)
2. ⏳ Wait 
3. Launch the "Zephyr" agent (the Dungeon Master with user interface) (see above)
4. The game begins, the player is in the starting room (0,0)

### At game launch (Dungeon Master)

```
docker attach $(docker compose ps -q dungeon-master)
/app # ./dungeon-master
```

**If everything goes well**, you should see a display similar to this:
```
🌍 LLM URL: http://model-runner.docker.internal/engines/v1/
🌍 MCP Host: http://mcp-gateway:9011/mcp
🌍 Dungeon Master Model: hf.co/menlo/jan-nano-gguf:q4_k_m
MCP Client initialized successfully
Tool: collect_gold - Collect gold coins from the current room if available. Try: "Collect the gold coins"
Tool: collect_magic_potion - Collect magic potions from the current room if available. Try: "Collect the magic potions"
Tool: create_player - Create a new player. Try: "I'm Bob, the Dwarf Warrior."
Tool: fight_monster - Fight a monster in your current room using turn-based combat. Each call represents one combat turn with dice rolls for both player and monster.
Tool: get_current_room_info - Get information about the current room where the player is located. Try: "Where am I?" or "Look around"
Tool: get_dungeon_info - Get the current dungeon's information including its layout, rooms, entrance and exit coordinates.
Tool: get_dungeon_map - Generate an ASCII map of the discovered dungeon rooms showing the player position, NPCs, and monsters with a legend.
Tool: get_player_info - Get the current player's information. Try: "Who am I?"
Tool: move_by_direction - Move the player in a specified direction (north, south, east, west). Try "move by north".
Tool: move_player - Move the player in the dungeon by specifying a cardinal direction. This is the primary navigation tool for exploring rooms. Usage: "move player north" or "go east".
Tool: speak_to_somebody - Speak to somebody by name

🔶 Loading vector store from: ./data/thrain_vector_store.json
✅ Vector store loaded successfully with 11 records
🔶 Loading vector store from: ./data/elara_vector_store.json
✅ Vector store loaded successfully with 11 records
🔶 Loading vector store from: ./data/galdor_vector_store.json
✅ Vector store loaded successfully with 11 records
🔶 Loading vector store from: ./data/liora_vector_store.json
✅ Vector store loaded successfully with 11 records
🔍 Pinging agent...
✅ Connected to agent: Shesepankh
📝 Description: An ancient and wise Sphinx who guards the exit of the dungeon. Known for posing riddles to those who seek passage.
🔧 Available skills: 1

Agent ID: zephyr agent name: Zephyr model: hf.co/menlo/jan-nano-gguf:q4_k_m
Agent ID: casper agent name: Casper model: ghost-model
Agent ID: thrain agent name: Thrain model: ai/qwen2.5:1.5B-F16
Agent ID: elara agent name: Elara model: hf.co/menlo/lucy-gguf:q8_0
Agent ID: galdor agent name: Galdor model: ai/qwen2.5:0.5B-F16
Agent ID: liora agent name: Liora model: ai/qwen2.5:0.5B-F16
Agent ID: shesepankh agent name: Shesepankh model: Remote Model

┃ 🤖 (/bye to exit) [Zephyr]>                                                                                                    
┃ Type your command here...                                                                                                      
┃                                                                                                                                
┃                                                                                                                                
┃                                                                                                                                
┃                                                                                                                                
┃                                                                                                                                                                                                                                                            
alt+enter / ctrl+j new line • enter submit
```
> - By default the selected agent to converse with you is "Zephyr"
> - 👋 For now it logs a lot of debug messages (DEBUG level) - to be cleaned up... Or not.

### Character Creation

You must start by creating a character by entering their name. For example:

```raw
┃ 🤖 (/bye to exit) [Zephyr]>                                                                                                    
┃ Hello, I'm Bob  the Dwarf Warrior                                                                                                         
```

```raw
< Zephyr speaking...>
⠹ Tools detection.....
⠇ Tools detection.....🟢 create_player with arguments: {"class":"warrior","name":"Bob","race":"dwarf"}
Do you want to execute this function? (y)es (n)o (a)bort (y/n/a) [y]: y  
```

```raw
⠏ Tools detection.....✅ Tool executed successfully
---[MCP RESPONSE]---------------------------------
{
  "name": "Bob",
  "level": 1,
  "class": "warrior",
  "race": "dwarf",
  "position": {
    "x": 0,
    "y": 0
  },
  "room_id": "room_0_0",
  "health": 100,
  "strength": 10,
  "experience": 0,
  "gold_coins": 0,
  "is_dead": false
}
--------------------------------------------------
---[DM RESPONSE]----------------------------------
Welcome, Bob the Dwarf Warrior! You have been successfully created as a warrior of the dwarven race. You are currently in room_0_0 at coordinates (0, 0). You have 100 health, 10 strength, and 0 experience. You are equipped with no gold coins and are not dead. 

Where would you like to go? Would you like to look around the room or move in a particular direction?
--------------------------------------------------

┃ 🤖 (/bye to exit) [Zephyr]>                                                                                                                        
┃ Type your command here...                                                                                                                          
┃
```

### List of Available MCP Tools

🎉 You can now move around the dungeon! And you have access to several MCP tools to interact with the dungeon:

- `collect_gold` - Collect gold coins from the current room if available. Try: "Collect the gold coins"
- `collect_magic_potion` - Collect magic potions from the current room if available. Try: "Collect the magic potions"
- `create_player` - Create a new player. Try: "I'm Bob, the Dwarf Warrior."
- `fight_monster` - Fight a monster in your current room using turn-based combat. Each call represents one combat turn with dice rolls for both player and monster.
- `get_current_room_info` - Get information about the current room where the player is located. Try: "Where am I?" or "Look around"
- `get_dungeon_info` - Get the current dungeon's information including its layout, rooms, entrance and exit coordinates.
- `get_dungeon_map` - Generate an ASCII map of the discovered dungeon rooms showing the player position, NPCs, and monsters with a legend.
- `get_player_info` - Get the current player's information. Try: "Who am I?"
- `move_by_direction` - Move the player in a specified direction (north, south, east, west). Try "move by north".
- `move_player` - Move the player in the dungeon by specifying a cardinal direction. 
- `speak_to_somebody` - Speak to somebody by name

> You can type the `/tools` command to see the list of available tools.

#### `get_current_room_info`

**You are at the dungeon entrance, you can start by looking around:**
```raw
┃ 🤖 (/bye to exit) [Zephyr]>                                                                                                  
┃ give me information about the room      
```

**Zephyr's response:**
```raw
< Zephyr speaking...>
⠧ Tools detection.....🟢 get_current_room_info with arguments: {}
Do you want to execute this function? (y)es (n)o (a)bort (y/n/a) [y]:   
```

**Then:**
```raw
⠇ Tools detection.....✅ Tool executed successfully
---[MCP RESPONSE]---------------------------------
{
  "id": "room_0_0",
  "name": "Ancient Ruins' Entrance",
  "description": "The entrance to the ancient dungeon, dimly lit by flickering torches and the eerie glow of ancient runes etched into the stone walls. The air is thick with the scent of damp stone and forgotten magic, and the floor creaks softly as if it's alive with unseen threats. The walls are lined with narrow, winding corridors that seem to twist and turn without end, each corner filled with secrets waiting to be discovered. ",
  "is_entrance": true,
  "is_exit": false,
  "coordinates": {
    "x": 0,
    "y": 0
  },
  "visited": true,
  "has_monster": false,
  "has_non_player_character": false,
  "has_treasure": false,
  "gold_coins": 0,
  "has_magic_potion": false,
  "regeneration_health": 0
}
--------------------------------------------------
---[DM RESPONSE]----------------------------------
**Ancient Ruins' Entrance**

- **Description**: The entrance to the ancient dungeon, dimly lit by flickering torches and the eerie glow of ancient runes etched into the stone walls. The air is thick with the scent of damp stone and forgotten magic, and the floor creaks softly as if it's alive with unseen threats. The walls are lined with narrow, winding corridors that seem to twist and turn without end, each corner filled with secrets waiting to be discovered.
- **Is Entrance**: Yes
- **Is Exit**: No
- **Coordinates**: x=0, y=0
- **Visited**: Yes
- **Has Monster**: No
- **Has Non-Player Character**: No
- **Has Treasure**: No
- **Gold Coins**: 0
- **Has Magic Potion**: No
- **Regeneration Health**: 0

What would you like to do next?
--------------------------------------------------
```

> ✋ Plan to be able to show or hide the MCP Response.

#### `move_by_direction`

**You can now move around the dungeon. For example, to go north:**
```raw
┃ 🤖 (/bye to exit) [Zephyr]>                                                                                                  
┃ I want to move to the north 
```
> Theoretically I can even say "I want to go to the north, then to the east, then to the north again" and Zephyr should understand and make the 3 movements. (🐛 fix: Zephyr detects the 3 movements well, performs them but displays the room information of the arrival of the 1st movement. [TODO: to fix - not priority])

**Zephyr's response:**
```raw
< Zephyr speaking...>
⠧ Tools detection.....🟢 move_player with arguments: {"direction":"north"}
Do you want to execute this function? (y)es (n)o (a)bort (y/n/a) [y]: 
```

**Then:**
```raw
⠦ Tools detection.....✅ Tool executed successfully
---[MCP RESPONSE]---------------------------------
✅ Moved north to position (0, 1).
🏠 Room name:The Rune Altar
📝 Description:The room is dimly lit by the flickering torches hanging from the ceiling. Cobwebs and dust cover the stone walls, and a strong, musty smell permeates the air. In the center of the room lies a pedestal, adorned with ancient runes and a glowing amulet.
🧪 There is a magic potion here that can restore 14 health points!
--------------------------------------------------
---[DM RESPONSE]----------------------------------
You have moved to the **Rune Altar** room. The room is dimly lit by flickering torches, with cobwebs and dust covering the stone walls. A strong musty smell is in the air. In the center of the room is a pedestal with ancient runes and a glowing amulet. There is also a **magic potion** here that can restore **14 health points**!

Would you like to **collect the magic potion** or **investigate the pedestal and amulet**?
--------------------------------------------------
```

You can see that it's possible to collect the magic potion present in the room.

#### `collect_magic_potion`

**To collect the magic potion, just ask for it:**
```raw
┃ 🤖 (/bye to exit) [Zephyr]>
┃ I want to collect the magic potion
```

**Zephyr's response:**
```raw
< Zephyr speaking...>
⠹ Tools detection.....🟢 collect_magic_potion with arguments: {}
Do you want to execute this function? (y)es (n)o (a)bort (y/n/a) [y]:                                                          
```

**Then:**
> 👋 In this example, we can see that sometimes the agent suggests other tools, but it's possible to exit the loop
```raw
⠸ Tools detection.....✅ Tool executed successfully
⠦ Tools detection.....🟢 collect_gold with arguments: {}
Do you want to execute this function? (y)es (n)o (a)bort (y/n/a) [y]: n                                                                                
⠋ Tools detection.....🟢 collect_gold with arguments: {}
Do you want to execute this function? (y)es (n)o (a)bort (y/n/a) [y]: n                                                                                
---[MCP RESPONSE]---------------------------------
🧪 You collected a magic potion from The Rune Altar! You gained 14 health points. Your current health: 114
--------------------------------------------------
---[DM RESPONSE]----------------------------------
It seems that there are no gold coins in the current room. Let's check the room again to see if there are any other items or if we can move to a different room where gold might be located. What would you like to do next?
--------------------------------------------------
```

#### `get_player_info`

**To check your character's status, you can ask for their information:**
```raw
┃ 🤖 (/bye to exit) [Zephyr]>
┃ Give me information about myself
```

**Zephyr's response:**
```raw
< Zephyr speaking...>
⠹ Tools detection.....🟢 get_player_info with arguments: {}
Do you want to execute this function? (y)es (n)o (a)bort (y/n/a) [y]:                                                                                  
```

**Then:**
```raw
⠹ Tools detection.....✅ Tool executed successfully
---[MCP RESPONSE]---------------------------------
{
  "name": "Bob",
  "level": 1,
  "class": "warrior",
  "race": "dwarf",
  "position": {
    "x": 0,
    "y": 1
  },
  "room_id": "room_0_1",
  "health": 114,
  "strength": 10,
  "experience": 0,
  "gold_coins": 0,
  "is_dead": false
}
--------------------------------------------------
---[DM RESPONSE]----------------------------------
You are Bob, a level 1 warrior dwarf. You are currently in room_0_1, located at coordinates (0, 1). You have 114 health points and a strength of 10. You have 0 gold coins and 0 experience. You are not dead. What would you like to do next?
--------------------------------------------------
```

#### `get_dungeon_map`

**To display the dungeon map, you can use the `get_dungeon_map` tool:**
```raw
┃ 🤖 (/bye to exit) [Zephyr]>
┃ Show me the dungeon map
```

**Zephyr's response:**
```raw
< Zephyr speaking...>
⠙ Tools detection.....🟢 get_dungeon_map with arguments: {}
Do you want to execute this function? (y)es (n)o (a)bort (y/n/a) [y]:   
```

**Then:**
```raw
⠹ Tools detection.....✅ Tool executed successfully
---[MCP RESPONSE]---------------------------------
THE SQUARE DUNGEON OF COMPOSE-AND-DRAGONS
=========================================

    0       1       2       3       
  ┌───────┬───────┬───────┬───────┐
3 │ ???   │ ???   │ ???   │ ???   │
  │       │       │       │       │
  │       │       │       │       │
  ├───────┼───────┼───────┼───────┤
2 │       │       │       │ ???   │
  │ [G]   │ [G]   │ [G]   │       │
  │  ✓    │  ✓    │  ✓    │       │
  ├───────┼───────┼───────┼───────┤
1 │       │ ???   │       │ ???   │
  │       │       │ [@][+]│       │
  │  ✓    │       │  ✓    │       │
  ├───────┼───────┼───────┼───────┤
0 │       │ ???   │ ???   │ ???   │
  │ [E]   │       │       │       │
  │  ✓    │       │       │       │
  └───────┴───────┴───────┴───────┘

LEGEND:
=======
[@] - Player (Bob the Warrior)
[E] - Entrance
[G] - Goblin (Blue Moon)
[+] - Healer (Liora - Half-Elf)
 ✓  - Visited room
??? - Unvisited/Empty room

ROOM DETAILS:
=============
(0,0) Ancient Ruins' Entrance - ENTRANCE
(0,1) The Rune Altar
(0,2) Stairway to the Abyss - Has Guard
(1,2) The Rune-Lit Foggery - Has Goblin
(2,2) Tormented Hall - Has Goblin
(2,1) The Whispering Shadows - Has Healer (Current Location)

PLAYER STATUS:
==============
Name: Bob
Class: Warrior (Dwarf)
Level: 1
Health: 114/100
Strength: 10
Experience: 0
Gold: 0

Current Position: (2,1) - The Whispering Shadows

--------------------------------------------------
---[DM RESPONSE]----------------------------------
The dungeon map has been generated, and here it is:

    0       1       2       3       
  ┌───────┬───────┬───────┬───────┐
3 │ ???   │ ???   │ ???   │ ???   │
  │       │       │       │       │
  │       │       │       │       │
  ├───────┼───────┼───────┼───────┤
2 │       │       │       │ ???   │
  │ [G]   │ [G]   │ [G]   │       │
  │  ✓    │  ✓    │  ✓    │       │
  ├───────┼───────┼───────┼───────┤
1 │       │ ???   │       │ ???   │
  │       │       │ [@][+]│       │
  │  ✓    │       │  ✓    │       │
  ├───────┼───────┼───────┼───────┤
0 │       │ ???   │ ???   │ ???   │
  │ [E]   │       │       │       │
  │  ✓    │       │       │       │
  └───────┴───────┴───────┴───────┘


Legend:
- [@] - Player (Bob the Warrior)
- [E] - Entrance
- [G] - Goblin (Blue Moon)
- [+] - Healer (Liora - Half-Elf)
- ✓ - Visited room
- ??? - Unvisited/Empty room

Room Details:
- (0,0) Ancient Ruins' Entrance - ENTRANCE
- (0,1) The Rune Altar
- (0,2) Stairway to the Abyss - Has Guard
- (1,2) The Rune-Lit Foggery - Has Goblin
- (2,2) Tormented Hall - Has Goblin
- (2,1) The Whispering Shadows - Has Healer (Current Location)

Player Status:
- Name: Bob
- Class: Warrior (Dwarf)
- Level: 1
- Health: 114/100
- Strength: 10
- Experience: 0
- Gold: 0

Current Position: (2,1) - The Whispering Shadows

What would you like to do next?
--------------------------------------------------
```

**Etc. ...**

### Talking to an NPC

**[🚧 for now you can invoke a character from anywhere]** TODO: verify that the NPC is in the current room before being able to talk to them.

To talk to an NPC, just address them by name. For example, to talk to the merchant "Galdor":
```raw
┃ 🤖 (/bye to exit) [Zephyr]>
┃ I want to speak to Galdor
```

**Zephyr's response:**
```raw
< Zephyr speaking...>
⠸ Tools detection.....🟢 speak_to_somebody with arguments: {"name":"Galdor"}
Do you want to execute this function? (y)es (n)o (a)bort (y/n/a) [y]:     
```

**Then:**
```raw
---[MCP RESPONSE]---------------------------------
{"result": "😃 You speak to {"name":"Galdor"}. They greet you warmly and are eager to assist you on your quest."}
--------------------------------------------------
---[DM RESPONSE]----------------------------------
Galdor is a friendly and helpful character in the dungeon. He has offered to assist you on your quest. What would you like to do next?
--------------------------------------------------

┃ 🙂 (/bye to exit /dm to go back to the DM) [Galdor]>                                                                                                  
┃ Type your command here...                                                                                                                             
┃                                                                                                                                                       
┃                                                   
```

Now you can converse with "Galdor". To return to "Zephyr", just type the `/dm` command.

```raw
┃ 🙂 (/bye to exit /dm to go back to the DM) [Galdor]>                                                                                                  
┃ Hello, I'm Bob, tell me something about your family   
```

> The agent will perform similarity searches in its vector store to answer the question. This allows providing a lot of information to the model without overloading the prompt/context.

**Galdor's response:**

```raw
< Galdor speaking...>
🔍 Searching for similar chunks to 'Hello, I'm Bob, tell me something about your family'
--------------------------------------------------------------------------------
📝 Similarities found: 2
✅ CosineSimilarity: 0.5777657278112519 Chunk: ## Family
Galdor comes from a family of craftsmen. His father was a renowned blacksmith, his mother a jeweler. He has two brothers who run the family forge in the mountain stronghold of Khaz Ankor.
✅ CosineSimilarity: 0.5398634678994042 Chunk: ## Quote
"Good coin for good goods - that's the foundation of honest trade, and honest trade builds kingdoms."
--------------------------------------------------------------------------------
Hello Bob, thank you for asking. Galdor is descended from a family of craftsmen, including his father a blacksmith and mother a jeweler. His family has a long history of trade and craftsmanship, passed down through generations. His brother, who runs the family forge in Khaz Ankor, is a skilled artisan himself.
```