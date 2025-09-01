# Compose and Dragons

## Setup de VSCode

Le dossier `.vscode` est **obligatoire**. Certains réglages et extensions sont utilisés pour améliorer l'affichage des documents et du code source et pour faciliter la présentation.

## Lancement du projet Docker Compose

À la racine du projet, lancer la commande suivante pour démarrer les conteneurs Docker :
```bash
docker compose up --build -d
docker compose logs -f mcp-dungeon mcp-gateway dungeon-end-of-level-boss
```
> Ou utiliser `build-start-dungeon.sh`

### Qu'est-ce qui est lancé ?

```bash
docker compose ps --services
```

- `mcp-dungeon` (projet: `dungeon-crawler-mcp-server`) : le serveur MCP personnalisé pour le jeu de donjon. (transport utilise: Streamable HTTP)
  - Le serveur est lancé automatiquement
- `mcp-gateway` : le serveur MCP Gateway qui permet de se connecter au serveur MCP `mcp-dungeon` et potentiellement à d'autres serveurs MCP.
  - La gateway est lancée automatiquement
- `dungeon-end-of-level-boss` (projet: `dungeon-end-of-level-boss`) : un agent d'IA **"fonctionnel"** "indépendant" qui utilise 2 agents techniques (1 pour la complétion de chat, l'autre pour le RAG).
  - C'est un PNJ "Boss de fin de niveau", nommé **"Shesepankh*"** qui sera utilisé à la sortie du donjon.
  - Les données de configuration de l'agent sont dans le dossier `dungeon-end-of-level-boss/data`
  - L'agent est lancé automatiquement: au 1er démarrage, l'agent va construire le vector store dans un fichier JSON à partir des données de configuration. Cela peut prendre un petit moment. Ce sera ensuite instantané aux prochains démarrages.
- `dungeon-master` (projet: `dungeon-master`) : il est composé de **plusieurs** agents d'IA **"fonctionnels"**:
  - **"Zephyr"** : on peut considérer que c'est le "maître du donjon" (Dungeon Master) qui gère l'état du jeu et interagit avec le joueur.
    - Le container est démarré mais le programme `dungeon-master` n'est **pas** lancé automatiquement.
    - Il va détecter les commandes du joueur et faire appel aux autres agents **"fonctionnels"** si nécessaire.
    - Il va utiliser les MCP Tools du serveur MCP `mcp-dungeon` pour répondre aux actions du joueur.
    - Il utilise 1 agent technique (1 pour la complétion de chat & la détection des tool calls).
  - Les autres agents **"fonctionnels"** (les PNJ) sont invoqués par "Zephyr" **à la demande** du joueur.
    - **"Galdor"** : un PNJ Marchand
    - **"Elara"** : une PNJ Sorcière
    - **"Thrain"** : un PNJ Guardien
    - **"Liora"** : une PNJ Guérisseuse
    - Chacun d'eux utilise 2 agents techniques (1 pour la complétion de chat, l'autre pour le RAG).
    - ✋👻 *il existe aussi un "Ghost agent" qui est un fake d'agent IA à destination de tests - pas d'utilité pour le gameplay*.

### Lancement de l'agent "Zephyr" (le Dungeon Master avec l'interface utilisateur)

Pour lancer l'inteface utilisateur, il faut exécuter la commande suivante dans un nouveau terminal :
```bash
docker attach $(docker compose ps -q dungeon-master)
```
Au 1er démarrage, les 4 agents PNJ vont chacun construire un vector store dans un fichier JSON à partir des données de configuration (présentes dans le dossier `dungeon-master/data`). Cela peut prendre un petit moment. Ce sera ensuite instantané aux prochains démarrages.

## Principe de fonctionnement du jeu

- Le donjon est extrêmement simple: un carré de 4x4. (paramétrable dans le fichier `compose.yml`)
- Les 4 PMJs ainsi que le Boss de fin de niveau sont placé à l'avance dans des pièces spécifiques.
  - Ce paramétrage de positions se fait dans le fichier `compose.yml`
- On utilise le principe du "Dungeon Crawler" (exploration de donjon) avec des pièces générées dynamiquement:
  - Une fois une pièce générée, elle est stockée dans le serveur MCP `mcp-dungeon` (outil: `store_room`) et ne sera plus modifiée (le joueur peut revenir dans la pièce).
  - Lors de la génération d'une pièce, on peut aussi génèrer :
    - objets (potions, pièces d'or) qui sont stockés dans le serveur MCP `mcp-dungeon`.
    - ennemis (des monstres) qui sont aussi stockés dans le serveur MCP `mcp-dungeon`.
    - Le joueur peut interagir avec les objets et les ennemis (combat, ramasser des objets, etc.)
- Le joueur peut se déplacer dans les 4 directions (nord, sud, est, ouest) et interagir avec les PNJs.
- Le but du jeu est de trouver et vaincre le Boss de fin de niveau (**Shesepankh**) qui se trouve dans une pièce spécifique.
- **[🚧 pas encore implémenté (1)]** Le joueur va collecter auprès des PNJs des informations qui lui serviront lors de sa discussion avec **Shesepankh**. Il faudra donc que le joueur rencontre tous les PNJs et leur pose les bonnes questions.
- **[🚧 pas encore implémenté (2)]** Le joueur devra donner les informations collectées à **Shesepankh** pour qu'elle accepte de le laisser sortir.

> - (1): cela doit pouvoir se faire uniquement à base de prompt.
> - (2): cela doit pouvoir se faire à base de prompt + MCP tool(s).

## Déroulement d'une partie

1. Lancer le projet Docker Compose (voir plus haut)
2. ⏳ Patienter 
3. Lancer l'agent "Zephyr" (le Dungeon Master avec l'interface utilisateur) (voir plus haut)
4. Le jeu commence, le joueur se trouve dans la pièce de départ (0,0)

### Au lancement du jeu (du Dungeon Master)

```
docker attach $(docker compose ps -q dungeon-master)
/app # ./dungeon-master
```

**Si tout va bien**, vous devriez avoir un affichage similaire à celui-ci:
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
> - Par défaut l'agent sélectionné pour converser avec vous est "Zephyr"
> - 👋 Pour le moment ça log un paquet de messages de debug (niveau DEBUG) - à nettoyer... Ou pas.

### Création du personnage

Il faut commencer par créer un personnage en entrant son nom. Par exemple:

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

### Liste des outils MCP disponibles

🎉 Vous pouvez maintenant vous déplacer dans le donjon ! Et vous avez accès à plusieurs outils MCP pour interagir avec le donjon:

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


> Vous pouvez taper la commande `/tools` pour voir la liste des outils disponibles.

#### `get_current_room_info`

**Vous êtes à l'entrée du donjon, vous pouvez commencer par regarder autour de vous:**
```raw
┃ 🤖 (/bye to exit) [Zephyr]>                                                                                                  
┃ give me information about the room      
```

**Réponse de Zephyr:**
```raw
< Zephyr speaking...>
⠧ Tools detection.....🟢 get_current_room_info with arguments: {}
Do you want to execute this function? (y)es (n)o (a)bort (y/n/a) [y]:   
```

**Puis:**
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

> ✋ Prévoir de pouvoir afficher ou masquer la MCP Response.

#### `move_by_direction`

**Vous pouvez maintenant vous déplacer dans le donjon. Par exemple, pour aller au nord:**
```raw
┃ 🤖 (/bye to exit) [Zephyr]>                                                                                                  
┃ I want to move to the north 
```
> Théoriquement je peux même dire "I want to go to the north, then to the east, then to the north again" et Zephyr doit comprendre et faire les 3 déplacements. (🐛 fix: Zephyr détecte bien les 3 mouvements, les effectue mais affiche les informations de la pièce d'arrivée du 1er mouvement. [TODO: à fixer - pas prioritaire])

**Réponse de Zephyr:**
```raw
< Zephyr speaking...>
⠧ Tools detection.....🟢 move_player with arguments: {"direction":"north"}
Do you want to execute this function? (y)es (n)o (a)bort (y/n/a) [y]: 
```

**Puis:**
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

On peut voir qu'il est possible de ramasser la potion magique présente dans la pièce.

#### `collect_magic_potion`

**Pour ramasser la potion magique, il suffit de le demander:**
```raw
┃ 🤖 (/bye to exit) [Zephyr]>
┃ I want to collect the magic potion
```

**Réponse de Zephyr:**
```raw
< Zephyr speaking...>
⠹ Tools detection.....🟢 collect_magic_potion with arguments: {}
Do you want to execute this function? (y)es (n)o (a)bort (y/n/a) [y]:                                                          
```

**Puis:**
> 👋 Dans cet exemple, on voit que parfois l'agent propose d'autres outils, mais il est possible de sortir de la boucle
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

**Pour vérifier l'état de votre personnage, vous pouvez demander ses informations:**
```raw
┃ 🤖 (/bye to exit) [Zephyr]>
┃ Give me information about myself
```

**Réponse de Zephyr:**
```raw
< Zephyr speaking...>
⠹ Tools detection.....🟢 get_player_info with arguments: {}
Do you want to execute this function? (y)es (n)o (a)bort (y/n/a) [y]:                                                                                  
```

**Puis:**
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

**Pour afficher la carte du donjon, vous pouvez utiliser l'outil `get_dungeon_map`:**
```raw
┃ 🤖 (/bye to exit) [Zephyr]>
┃ Show me the dungeon map
```

**Réponse de Zephyr:**
```raw
< Zephyr speaking...>
⠙ Tools detection.....🟢 get_dungeon_map with arguments: {}
Do you want to execute this function? (y)es (n)o (a)bort (y/n/a) [y]:   
```

**Puis:**
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

### Parler avec un PNJ

**[🚧 pour le moment on peut invoquer un personnage de n'importe où]** TODO: vérifier que le PNJ est dans la pièce courrante avant de pouvoir lui parler.

Pour parler avec un PNJ, il suffit de lui adresser la parole par son nom. Par exemple, pour parler avec le marchand "Galdor":
```raw
┃ 🤖 (/bye to exit) [Zephyr]>
┃ I want to speak to Galdor
```

**Réponse de Zephyr:**
```raw
< Zephyr speaking...>
⠸ Tools detection.....🟢 speak_to_somebody with arguments: {"name":"Galdor"}
Do you want to execute this function? (y)es (n)o (a)bort (y/n/a) [y]:     
```

**Puis:**
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

Maintenant vous pouvez converser avec "Galdor". Pour revenir à "Zephyr", il suffit de taper la commande `/dm`.

```raw
┃ 🙂 (/bye to exit /dm to go back to the DM) [Galdor]>                                                                                                  
┃ Hello, I'm Bob, tell me something about your family   
```

> L'agent va faire des recherche de similarités dans son vector store pour répondre à la question. Cela permet de fournir beaucoup d'information au modèle sans surcharger le prompt/contexte.

**Réponse de Galdor:**

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