<!-- 
Prompt:
I'd like to create an ASCII representation in an README.md file like the one shown at https://d2lang.com/blog/ascii/ that explains how the main.go program works, as well as the executeFunction and generatePromptMessagesWithSimilarities functions.
-->

# Dungeon Master Program Architecture

## Main Program Flow

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                            DUNGEON MASTER MAIN PROGRAM                          │
└─────────────────────────────────────────────────────────────────────────────────┘

1. INITIALIZATION PHASE
   ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
   │   Environment   │    │   OpenAI Client │    │   MCP Client    │
   │   Variables     │ -> │   Creation      │ -> │   Creation      │
   │                 │    │   (LLM URL)     │    │   (Tools Host)  │
   └─────────────────┘    └─────────────────┘    └─────────────────┘
                                  │
                                  v
   ┌─────────────────────────────────────────────────────────────────────────────┐
   │                        TOOLS INDEX CREATION                                 │
   │                                                                             │
   │  ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────────┐      │
   │  │   MCP Tools     │    │  Custom Tools   │    │   Combined Index    │      │
   │  │   (External)    │ +  │  (speak_to_     │ =  │   (All Available    │      │
   │  │                 │    │   somebody)     │    │    Functions)       │      │
   │  └─────────────────┘    └─────────────────┘    └─────────────────────┘      │
   └─────────────────────────────────────────────────────────────────────────────┘
                                  │
                                  v
2. AGENTS CREATION PHASE
   ┌─────────────────────────────────────────────────────────────────────────────┐
   │                         AGENTS TEAM ASSEMBLY                                │
   │                                                                             │
   │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
   │  │   Dungeon   │  │    Ghost    │  │    Guard    │  │  Sorcerer   │         │
   │  │   Master    │  │    Agent    │  │    Agent    │  │    Agent    │         │
   │  │  (w/Tools)  │  │  (Testing)  │  │   (w/RAG)   │  │   (w/RAG)   │         │
   │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘         │
   │                                                                             │
   │  ┌─────────────┐  ┌─────────────┐                                           │
   │  │  Merchant   │  │   Healer    │                                           │
   │  │   Agent     │  │   Agent     │                                           │
   │  │  (w/RAG)    │  │  (w/RAG)    │                                           │
   │  └─────────────┘  └─────────────┘                                           │
   └─────────────────────────────────────────────────────────────────────────────┘
                                  │
                                  v
3. MAIN INTERACTION LOOP
   ┌─────────────────────────────────────────────────────────────────────────────┐
   │                         MAIN GAME LOOP                                      │
   │                                                                             │
   │  ┌─────────────┐    ┌─────────────┐    ┌─────────────────────────────────┐  │
   │  │    User     │    │   Command   │    │        Route to Agent           │  │
   │  │   Input     │ -> │  Parsing    │ -> │      (Based on Selected)        │  │
   │  │             │    │   (/bye,    │    │                                 │  │
   │  │             │    │   /dm, etc) │    │                                 │  │
   │  └─────────────┘    └─────────────┘    └─────────────────────────────────┘  │
   │                                                   │                         │
   │                   ┌───────────────────────────────┼───────────────────┐     │
   │                   │                               │                   │     │
   │                   v                               v                   v     │
   │        ┌─────────────────┐           ┌─────────────────┐   ┌─────────────┐  │
   │        │  Dungeon Master │           │  NPCs w/ RAG    │   │ Ghost Agent │  │
   │        │   (w/ Tools)    │           │  (Guard, Sorc,  │   │ (Testing)   │  │
   │        │                 │           │  Merchant, etc) │   │             │  │
   │        └─────────────────┘           └─────────────────┘   └─────────────┘  │
   └─────────────────────────────────────────────────────────────────────────────┘
```

## executeFunction Flow

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                            executeFunction WORKFLOW                             │
└─────────────────────────────────────────────────────────────────────────────────┘

INPUT: functionName (string), arguments (string)
   │
   v
┌─────────────────────────────────────────────────────────────────────────────────┐
│                        FUNCTION CALL DETECTION                                  │
│                                                                                 │
│  ┌─────────────────────────────────────────────────────────────────────────┐    │
│  │  1. Display function name and arguments                                 │    │
│  │  2. Pause thinking controller                                           │    │
│  │  3. Ask user for confirmation: (y)es (n)o (a)bort                       │    │
│  │  4. Resume thinking controller                                          │    │
│  └─────────────────────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────────────────────┘
   │
   v
┌─────────────────────────────────────────────────────────────────────────────────┐
│                           USER CHOICE ROUTING                                   │
│                                                                                 │
│        ┌─────────┐         ┌─────────┐         ┌─────────────────┐              │
│        │   "n"   │         │   "a"   │         │  "y" (default)  │              │
│        │  (No)   │         │ (Abort) │         │     (Yes)       │              │
│        └─────────┘         └─────────┘         └─────────────────┘              │
│            │                   │                       │                        │
│            v                   v                       v                        │
│  ┌─────────────────┐  ┌─────────────────┐   ┌─────────────────────────────┐     │
│  │ Return: Function│  │ Return: Function│   │    FUNCTION EXECUTION       │     │
│  │ not executed    │  │ not executed    │   │         ROUTING             │     │
│  │                 │  │ + Exit Error    │   │                             │     │
│  └─────────────────┘  └─────────────────┘   └─────────────────────────────┘     │
└─────────────────────────────────────────────────────────────────────────────────┘
                                                           │
                                                           v
┌─────────────────────────────────────────────────────────────────────────────────┐
│                         FUNCTION TYPE DETECTION                                 │
│                                                                                 │
│  ┌─────────────────────────────────────┐   ┌─────────────────────────────────┐  │
│  │        "speak_to_somebody"          │   │         OTHER FUNCTIONS         │  │
│  │        (Custom Tool)                │   │         (MCP Tools)             │  │
│  └─────────────────────────────────────┘   └─────────────────────────────────┘  │
│                    │                                       │                    │
│                    v                                       v                    │
│  ┌─────────────────────────────────────┐   ┌─────────────────────────────────┐  │
│  │ 1. Parse JSON arguments             │   │ 1. Call mcpClient.CallTool()    │  │
│  │ 2. Check if agent exists in team    │   │ 2. Extract result content       │  │
│  │ 3. Set selectedAgent if found       │   │ 3. Return tool result           │  │
│  │ 4. Return success/error message     │   │                                 │  │
│  └─────────────────────────────────────┘   └─────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

## generatePromptMessagesWithSimilarities Flow

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                  generatePromptMessagesWithSimilarities WORKFLOW                │
└─────────────────────────────────────────────────────────────────────────────────┘

INPUT: ctx, client, agentName, input, similarityLimit, maxResults
   │
   v
┌─────────────────────────────────────────────────────────────────────────────────┐
│                          RAG SIMILARITY SEARCH                                  │
│                                                                                 │
│  ┌─────────────────────────────────────────────────────────────────────────┐    │
│  │  Call: agents.SearchSimilarities()                                      │    │
│  │  Parameters:                                                            │    │
│  │  - Context                                                              │    │
│  │  - OpenAI Client                                                        │    │
│  │  - Agent Name                                                           │    │
│  │  - User Input (search query)                                            │    │
│  │  - Similarity Threshold (e.g., 0.5)                                     │    │
│  │  - Max Results (e.g., 2)                                                │    │
│  └─────────────────────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────────────────────┘
   │
   v
┌─────────────────────────────────────────────────────────────────────────────────┐
│                           SEARCH RESULTS ROUTING                                │
│                                                                                 │
│  ┌─────────────────────────────┐       ┌─────────────────────────────────────┐  │
│  │      SIMILARITIES FOUND     │       │       NO SIMILARITIES FOUND         │  │
│  │      (len(similarities)>0)  │       │       (len(similarities)==0)        │  │
│  └─────────────────────────────┘       └─────────────────────────────────────┘  │
│               │                                         │                       │
│               v                                         v                       │
│  ┌─────────────────────────────┐       ┌─────────────────────────────────────┐  │
│  │ BUILD CONTEXT MESSAGE:      │       │ RETURN SIMPLE MESSAGE:              │  │
│  │                             │       │                                     │  │
│  │ 1. Create context header    │       │ Return: []ChatCompletionMessage{    │  │
│  │    "Here is some context    │       │   UserMessage(input)                │  │
│  │     that might be useful:"  │       │ }                                   │  │
│  │                             │       │                                     │  │
│  │ 2. For each similarity:     │       │                                     │  │
│  │    - Append "- {prompt}\n"  │       │                                     │  │
│  │                             │       │                                     │  │
│  │ 3. Return: []ChatCompletion │       │                                     │  │
│  │    Message{                 │       │                                     │  │
│  │      SystemMessage(context) │       │                                     │  │
│  │      UserMessage(input)     │       │                                     │  │
│  │    }                        │       │                                     │  │
│  └─────────────────────────────┘       └─────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

## Agent Interaction Pattern

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                        AGENT CONVERSATION PATTERNS                              │
└─────────────────────────────────────────────────────────────────────────────────┘

1. DUNGEON MASTER (with Tools)
   ┌─────────────────────────────────────────────────────────────────────────┐
   │  User Input -> System Instructions + User Message                       │
   │       │                                                                 │
   │       v                                                                 │
   │  DetectToolCalls() with executeFunction callback                        │
   │       │                                                                 │
   │       v                                                                 │
   │  Tool Detection -> User Confirmation -> Tool Execution                  │
   │       │                                                                 │
   │       v                                                                 │
   │  Assistant Response (with tool results)                                 │
   └─────────────────────────────────────────────────────────────────────────┘

2. NPCs with RAG (Guard, Sorcerer, Merchant, Healer)
   ┌─────────────────────────────────────────────────────────────────────────┐
   │  User Input -> generatePromptMessagesWithSimilarities()                 │
   │       │                                                                 │
   │       v                                                                 │
   │  Similarity Search -> Context Building                                  │
   │       │                                                                 │
   │       v                                                                 │
   │  [SystemMessage(context)] + UserMessage(input)                          │
   │       │                                                                 │
   │       v                                                                 │
   │  agent.RunStream() -> Streaming Response                                │
   └─────────────────────────────────────────────────────────────────────────┘

3. Ghost Agent (Testing)
   ┌─────────────────────────────────────────────────────────────────────────┐
   │  User Input -> Simple System Message + User Message                     │
   │       │                                                                 │
   │       v                                                                 │
   │  agent.RunStream() -> Direct Streaming Response                         │
   └─────────────────────────────────────────────────────────────────────────┘
```

## Key Data Structures

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                            GLOBAL STATE                                         │
└─────────────────────────────────────────────────────────────────────────────────┘

agentsTeam: map[string]mu.Agent
┌─────────────────────────────────────────────────────────────────────────────┐
│  Key: agentId (lowercase)     │  Value: Agent Instance                      │
│  ──────────────────────────── │  ─────────────────────                      │
│  "sam"                        │  dungeonMasterToolsAgent                    │
│  "casper"                     │  ghostAgent                                 │
│  "guard_name"                 │  guardAgent                                 │
│  "sorcerer_name"              │  sorcererAgent                              │
│  "merchant_name"              │  merchantAgent                              │
│  "healer_name"                │  healerAgent                                │
└─────────────────────────────────────────────────────────────────────────────┘

selectedAgent: mu.Agent (currently active agent)

conversationalMemory: []openai.ChatCompletionMessageParamUnion
┌─────────────────────────────────────────────────────────────────────────────┐
│  [0]: SystemMessage (Dungeon Master instructions)                           │
│  [1]: UserMessage (user input 1)                                            │
│  [2]: AssistantMessage (DM response 1)                                      │
│  [3]: UserMessage (user input 2)                                            │
│  [4]: AssistantMessage (DM response 2)                                      │
│  ...                                                                        │
└─────────────────────────────────────────────────────────────────────────────┘
```