# Dungeon MCP Server Architecture

This document explains the architecture and flow of the Dungeon MCP Server application using various diagrams and explanations.

## System Architecture Overview

```mermaid
graph TB
    subgraph "Main Application Flow"
        A["main Function"] --> B[MCP Server Creation]
        A --> C[Dungeon Agent Setup]
        A --> D[Game Initialization]
        A --> E[Tools Registration]
        A --> F[HTTP Server Setup]
    end
    
    B --> B1["NewMCPServer<br/>dungeon-mcp-server, 0.0.0"]
    
    C --> C1[OpenAI Client Setup]
    C1 --> C2[Micro Agent Creation]
    
    D --> D1[Player Creation]
    D --> D2[Dungeon Creation]
    D --> D3[Entrance Room Generation]
    
    E --> E1[Register 10 MCP Tools]
    
    F --> F1[HTTP ServeMux Setup]
    F1 --> F2[Server Start on Port 9090]
    
    style A fill:#e1d5e7
    style B fill:#dae8fc
    style C fill:#fff2cc
    style D fill:#f8cecc
    style E fill:#e6ccff
    style F fill:#ffe6cc
```

### System Architecture Explanation

This diagram shows the high-level initialization flow of the Dungeon MCP Server. The main function orchestrates six primary components: MCP Server creation, AI Agent setup, Game state initialization, MCP Tools registration, and HTTP Server configuration. Each component is responsible for a specific aspect of the dungeon crawler game infrastructure, working together to provide a complete Model Context Protocol (MCP) server that can generate and manage dungeon exploration experiences.

## Agent Setup and Configuration Flow

```mermaid
sequenceDiagram
    participant Main as Main Function
    participant Config as Environment Config
    participant OpenAI as OpenAI Client
    participant Agent as Dungeon Agent
    participant Schema as JSON Schema
    
    Main->>Config: Get MODEL_RUNNER_BASE_URL
    Main->>Config: Get DUNGEON_MODEL
    Main->>Config: Get DUNGEON_MODEL_TEMPERATURE
    
    Main->>OpenAI: NewClient(baseURL, apiKey)
    OpenAI-->>Main: client instance
    
    Main->>Schema: GetRoomSchema()
    Schema-->>Main: JSON schema for room response
    
    Main->>Agent: "mu.NewAgent with dungeon-agent"
    Note over Agent: Configured with:<br/>OpenAI client, Model parameters,<br/>JSON schema response format
    Agent-->>Main: dungeonAgent instance
```

### Agent Setup Explanation

This sequence diagram illustrates how the AI agent is configured and initialized. The system first retrieves environment variables for the AI model configuration, creates an OpenAI-compatible client, and then initializes a "micro agent" with specific parameters. The agent is configured to respond in a structured JSON format using a predefined schema for room generation, ensuring consistent and parseable responses when generating dungeon content.

## Game State Initialization Process

```mermaid
stateDiagram-v2
    [*] --> ConfigLoading
    
    ConfigLoading --> PlayerCreation : Load environment variables
    ConfigLoading --> DungeonCreation : Load dungeon parameters
    
    PlayerCreation --> PlayerReady : Create Player with Name Unknown
    
    DungeonCreation --> DungeonParametersSet : Set width, height, entrance, exit
    DungeonParametersSet --> EntranceRoomGeneration : Initialize empty dungeon
    
    EntranceRoomGeneration --> AgentCall : Call dungeonAgent.Run()
    AgentCall --> JSONParsing : Parse AI response
    JSONParsing --> RoomCreation : Create entrance room object
    RoomCreation --> DungeonReady : Add room to dungeon
    
    PlayerReady --> GameInitialized
    DungeonReady --> GameInitialized
    GameInitialized --> [*]
```

### Game State Initialization Explanation

This state diagram shows how the game initializes its core state components. The process begins by loading configuration from environment variables, then creates both a player object and dungeon structure in parallel. The dungeon initialization includes a special step where the AI agent generates the entrance room dynamically, calling the language model to create appropriate name and description content. This ensures that each game session can have unique, contextually appropriate entrance rooms while maintaining the required game structure.

## MCP Tools Registration Architecture

```mermaid
graph LR
    subgraph "MCP Server"
        Server[MCP Server Instance]
    end
    
    subgraph "Game State"
        Player[Player Object]
        Dungeon[Dungeon Object]
        Agent[Dungeon Agent]
    end
    
    subgraph "MCP Tools"
        T1[CreatePlayer]
        T2[GetPlayerInfo]
        T3[GetDungeonInfo]
        T4[MovePlayer]
        T5[GetCurrentRoom]
        T6[GetDungeonMap]
        T7[CollectGold]
        T8[CollectPotion]
        T9[FightMonster]
        T10[IsPlayerInSameRoomAsNPC]
    end
    
    T1 --> Server
    T2 --> Server
    T3 --> Server
    T4 --> Server
    T5 --> Server
    T6 --> Server
    T7 --> Server
    T8 --> Server
    T9 --> Server
    T10 --> Server
    
    T1 -.-> Player
    T1 -.-> Dungeon
    T2 -.-> Player
    T2 -.-> Dungeon
    T3 -.-> Dungeon
    T4 -.-> Player
    T4 -.-> Dungeon
    T4 -.-> Agent
    T5 -.-> Player
    T5 -.-> Dungeon
    T6 -.-> Dungeon
    T7 -.-> Player
    T7 -.-> Dungeon
    T8 -.-> Player
    T8 -.-> Dungeon
    T9 -.-> Player
    T9 -.-> Dungeon
    T10 -.-> Player
    T10 -.-> Dungeon
```

### MCP Tools Registration Explanation

This diagram illustrates how MCP (Model Context Protocol) tools are registered with the server and how they interact with the game state. Each tool is registered with the MCP server and has access to shared references to the player object, dungeon object, and in some cases the AI agent. The dotted lines show the data dependencies - for example, the MovePlayer tool needs access to all three components because it updates player position, modifies dungeon state, and may call the AI agent to generate new rooms as the player explores.

## HTTP Server Request Flow

```mermaid
sequenceDiagram
    participant Client as MCP Client
    participant HTTP as HTTP Server
    participant Mux as ServeMux
    participant Health as Health Handler
    participant MCP as MCP Handler
    participant Tools as MCP Tools
    
    Note over HTTP: Server starts on port 9090
    
    Client->>HTTP: GET /health
    HTTP->>Mux: Route request
    Mux->>Health: healthCheckHandler
    Health-->>Mux: status healthy response
    Mux-->>HTTP: JSON response
    HTTP-->>Client: 200 OK
    
    Client->>HTTP: POST /mcp
    HTTP->>Mux: Route request
    Mux->>MCP: StreamableHTTPServer
    MCP->>Tools: Execute tool based on request
    Tools-->>MCP: Tool response
    MCP-->>Mux: MCP protocol response
    Mux-->>HTTP: Response
    HTTP-->>Client: MCP response
```

### HTTP Server Request Flow Explanation

This sequence diagram shows how HTTP requests are handled by the server. The server uses a custom ServeMux to route requests to two different endpoints: `/health` for health checks (useful for Docker Compose deployments) and `/mcp` for actual MCP protocol communication. When MCP requests are received, they are processed by the StreamableHTTPServer which then delegates to the appropriate registered tools based on the request content. This architecture allows the same server to handle both infrastructure monitoring and the core dungeon crawler game logic.

## Data Flow During Room Generation

```mermaid
flowchart TD
    A[Player wants to move] --> B{Room exists at target coordinates?}
    B -->|No| C["Call dungeonAgent.Run"]
    B -->|Yes| D[Move to existing room]
    
    C --> C1["System Instruction:<br/>DUNGEON_AGENT_ROOM_SYSTEM_INSTRUCTION"]
    C1 --> C2["User Message:<br/>Generate room at coordinates X,Y"]
    C2 --> C3[AI Model Processing]
    C3 --> C4["JSON Response:<br/>{name: string, description: string}"]
    
    C4 --> E[json.Unmarshal response]
    E --> F[Create Room object]
    F --> G["Set room properties:<br/>ID, coordinates,<br/>Generated name/description,<br/>Game flags"]
    G --> H[Add room to dungeon.Rooms]
    H --> I[Move player to new room]
    
    D --> J[Update player position]
    I --> J
    J --> K[Return movement result]
```

### Room Generation Data Flow Explanation

This flowchart demonstrates the dynamic room generation process that occurs during gameplay. When a player attempts to move to an unexplored area, the system first checks if a room already exists at those coordinates. If not, it triggers the AI agent to generate appropriate room content. The AI receives system instructions and a user prompt, then returns structured JSON containing the room's name and description. This response is parsed and used to create a full room object with all necessary game properties before being added to the dungeon and allowing the player to enter.