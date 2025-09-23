# Dungeon Master Application

```mermaid
flowchart TD
    X[Dungeon Master Application]:::main
    X --> Initialize(Initialize):::process

    Initialize --> Environment("Load Environment Variables<br/><a href='/dungeon-master/main.go#L30'>main.go:30</a>"):::config
    Environment --> LLMConfig[LLM URL, MCP Host, Model]:::config
    Environment --> RAGConfig[Similarity Search Config]:::config

    Initialize --> OpenAIClient("Create OpenAI Client<br/><a href='/dungeon-master/main.go#L41'>main.go:41</a>"):::client
    Initialize --> MCPClient("Create MCP Client<br/><a href='/dungeon-master/main.go#L46'>main.go:46</a>"):::client

    MCPClient --> ToolsCatalog("Get Tools from MCP<br/><a href='/dungeon-master/main.go#L56'>main.go:56</a>"):::tools
    ToolsCatalog --> CustomTools("Add Custom Tools<br/><a href='/dungeon-master/main.go#L61'>main.go:61</a>"):::tools
    CustomTools --> SpeakToAgentTool[speak_to_somebody tool]:::tool

    OpenAIClient --> DungeonMasterAgent("Create Dungeon Master Agent<br/><a href='/dungeon-master/main.go#L86'>main.go:86</a>"):::agent
    DungeonMasterAgent --> SystemInstructions("System Instructions<br/><a href='/dungeon-master/main.go#L103'>main.go:103</a>"):::config

    Initialize --> NPCAgents("Create NPC Agents<br/><a href='003-schema-npc-agents-system.md'>See NPC Agents Schema</a>"):::process
    NPCAgents --> AgentsTeam("Create Agents Team Map<br/><a href='/dungeon-master/main.go#L153'>main.go:153</a>"):::team
    AgentsTeam --> DefaultSelection("Select Dungeon Master as Default<br/><a href='/dungeon-master/main.go#L162'>main.go:162</a>"):::selection

    DefaultSelection --> MainLoop("Main Game Loop<br/><a href='004-schema-main-loop.md'>See Game Loop Schema</a>"):::loop

    classDef main fill:#e1f5fe,stroke:#01579b,stroke-width:3px,color:#000
    classDef process fill:#f3e5f5,stroke:#4a148c,stroke-width:2px,color:#000
    classDef config fill:#fff3e0,stroke:#e65100,stroke-width:2px,color:#000
    classDef client fill:#e8f5e8,stroke:#1b5e20,stroke-width:2px,color:#000
    classDef tools fill:#e8eaf6,stroke:#283593,stroke-width:2px,color:#000
    classDef tool fill:#f9fbe7,stroke:#827717,stroke-width:2px,color:#000
    classDef agent fill:#f1f8e9,stroke:#33691e,stroke-width:2px,color:#000
    classDef team fill:#e0f2f1,stroke:#00695c,stroke-width:2px,color:#000
    classDef selection fill:#fff8e1,stroke:#f57f17,stroke-width:2px,color:#000
    classDef loop fill:#e4f7ff,stroke:#0277bd,stroke-width:3px,color:#000
```

## Main Function Flow (<a href="/dungeon-master/main.go#L26">main()</a>)

The main function orchestrates the entire Dungeon Master application:

1. **Environment Setup** (<a href="/dungeon-master/main.go#L30-L39">lines 30-39</a>)
   - Load LLM URL, MCP Host, and Model configuration
   - Configure similarity search parameters

2. **Client Initialization** (<a href="/dungeon-master/main.go#L41-L51">lines 41-51</a>)
   - Create OpenAI client for LLM interactions
   - Initialize MCP client for tool communication

3. **Tools Configuration** (<a href="/dungeon-master/main.go#L56-L78">lines 56-78</a>)
   - Retrieve tools catalog from MCP server
   - Add custom `speak_to_somebody` tool

4. **Agent Creation** (<a href="/dungeon-master/main.go#L86-L162">lines 86-162</a>)
   - **Dungeon Master Agent**: Main agent with tools capability
   - **NPC Agents**: See [NPC Agents Schema](003-schema-npc-agents-system.md) for detailed information
   - Create agents team map for easy access

5. **Game Loop** - See [Game Loop Schema](004-schema-main-loop.md) for detailed information
   - Handle user input and commands
   - Route conversations to appropriate agents
   - Process tool calls and game logic

## Key Components

### Agent Types
- **Dungeon Master** (<a href="/dungeon-master/main.go#L231-L264">lines 231-264</a>): Main agent with full tool access
- **NPC Agents**: See [NPC Agents Schema](003-schema-npc-agents-system.md) for detailed information about Ghost, Guard, Sorcerer, Merchant, Healer, and Boss agents

### Tool Execution Handler

See [Game Loop Schema](004-schema-main-loop.md) for detailed tool execution flow. Main categories:
1. **Custom Tools**: `speak_to_somebody` with room validation
2. **MCP Tools**: Delegated to MCP server

### RAG Integration (<a href="/dungeon-master/main.go#L611-L637">GeneratePromptMessagesWithSimilarities</a>)

Enhances NPC interactions with contextual information through similarity search. See [NPC Agents Schema](003-schema-npc-agents-system.md) for detailed RAG implementation.

### Game Commands

See [Game Loop Schema](004-schema-main-loop.md) for detailed command implementation:
- `/bye`: Exit the game
- `/back-to-dm` or `/dm`: Return to Dungeon Master
- `/memory`: Display conversation history
- `/agents`: List all available agents
- `/tools`: Display available tools

## Agent Team Structure

The agents team structure is detailed in [NPC Agents Schema](003-schema-npc-agents-system.md). The main structure includes:
- **Dungeon Master**: Main agent with tools
- **NPC Agents**: Ghost, Guard, Sorcerer, Merchant, Healer, Boss