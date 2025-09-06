# Dungeon Master Application Architecture

This document explains the architecture and flow of the Dungeon Master application (`main.go`) using various diagrams and explanations.

## System Architecture Overview

```mermaid
graph TB
    Main[Main Application] --> Config[Configuration]
    Main --> OpenAI[OpenAI Client]
    Main --> MCP[MCP Client]
    Main --> Tools[Tools System]
    Main --> Agents[Agents Team]
    Main --> GameLoop[Game Loop]
    
    Config --> |Environment Variables| OpenAI
    Config --> |Environment Variables| MCP
    
    Tools --> |Tool Index| GameLoop
    MCP --> |External Tools| Tools
    
    Agents --> DM[Dungeon Master]
    Agents --> Ghost[Ghost Agent]
    Agents --> Guard[Guard Agent]
    Agents --> Sorcerer[Sorcerer Agent]
    Agents --> Merchant[Merchant Agent]
    Agents --> Healer[Healer Agent]
    Agents --> Boss[Boss Agent]
    
    GameLoop --> RAG[RAG System]
    GameLoop --> Commands[Command System]
    GameLoop --> Executor[Function Executor]
    
    style Main fill:#dae8fc
    style DM fill:#ffcccc
    style Boss fill:#ff9999
    style RAG fill:#cce5ff
    style Tools fill:#d5e8d4
```

### Architecture Explanation

The Dungeon Master application is built around a multi-agent system that provides an interactive text-based gaming experience. The main application initializes various components including OpenAI clients for LLM communication, MCP (Model Context Protocol) clients for external tool integration, and a team of specialized agents representing different characters in the dungeon.

## Application Initialization Flow

```mermaid
sequenceDiagram
    participant Main as Main Function
    participant Config as Configuration
    participant OpenAI as OpenAI Client
    participant MCP as MCP Client
    participant Tools as Tools System
    participant Agents as Agents Team
    
    Main->>Config: Load environment variables
    Main->>OpenAI: Initialize client with LLM URL
    Main->>MCP: Initialize MCP client
    MCP-->>Main: Return tool index
    Main->>Tools: Add speak_to_somebody tool
    Main->>Agents: Create all agents
    Agents-->>Main: Return agents team map
    Main->>Main: Start game loop
```

### Initialization Explanation

The application starts by loading configuration from environment variables, then initializes the OpenAI client for LLM communication and the MCP client for external tool integration. It creates a comprehensive tools system and assembles a team of agents, each with specific roles and capabilities. Finally, it enters the main game loop to handle user interactions.

## Agents Team Structure

```mermaid
graph LR
    subgraph "Agents Team"
        DM[Dungeon Master<br/>- Tool Detection<br/>- System Instructions<br/>- Main Controller]
        Ghost[Ghost Agent<br/>- Test Agent<br/>- Simple Response]
        Guard[Guard Agent<br/>- RAG Enabled<br/>- Memory History]
        Sorcerer[Sorcerer Agent<br/>- RAG Enabled<br/>- Magic Knowledge]
        Merchant[Merchant Agent<br/>- RAG Enabled<br/>- Trading Logic]
        Healer[Healer Agent<br/>- RAG Enabled<br/>- Healing Knowledge]
        Boss[Boss Agent<br/>- Remote Agent<br/>- HTTP Endpoint]
    end
    
    DM -.-> |Controls| Ghost
    DM -.-> |Controls| Guard
    DM -.-> |Controls| Sorcerer
    DM -.-> |Controls| Merchant
    DM -.-> |Controls| Healer
    DM -.-> |Controls| Boss
    
    style DM fill:#ffcccc
    style Ghost fill:#e6ccff
    style Guard fill:#ccffcc
    style Sorcerer fill:#ccccff
    style Merchant fill:#ffffcc
    style Healer fill:#ffccff
    style Boss fill:#ff9999
```

### Agents Team Explanation

The agents team consists of seven distinct characters, each with unique capabilities. The Dungeon Master serves as the primary controller with tool detection capabilities. Most agents (Guard, Sorcerer, Merchant, Healer) are enhanced with RAG (Retrieval-Augmented Generation) for contextual responses. The Ghost agent is used for testing, while the Boss agent operates as a remote service via HTTP endpoints.

## Main Game Loop Flow

```mermaid
flowchart TD
    Start([Start Game Loop]) --> Prompt[Display Prompt]
    Prompt --> Input[Get User Input]
    Input --> CheckBye{"/bye command?"}
    CheckBye -->|Yes| Exit([Exit Application])
    CheckBye -->|No| CheckDM{"/dm command?"}
    CheckDM -->|Yes| SwitchDM[Switch to Dungeon Master]
    CheckDM -->|No| CheckCommands{Check Other Commands}
    CheckCommands -->|/memory| ShowMemory[Display Memory]
    CheckCommands -->|/agents| ShowAgents[Display Agents]
    CheckCommands -->|/tools| ShowTools[Display Tools]
    CheckCommands -->|Normal Input| ProcessAgent[Process with Selected Agent]
    
    ProcessAgent --> IsDM{Is Dungeon Master?}
    IsDM -->|Yes| DMProcess[DM with Tools Processing]
    IsDM -->|No| NPCProcess[NPC Processing]
    
    DMProcess --> ToolDetection[Tool Detection]
    ToolDetection --> ToolExecution[Tool Execution]
    ToolExecution --> DMResponse[Display DM Response]
    
    NPCProcess --> RAGCheck{RAG Enabled?}
    RAGCheck -->|Yes| RAGProcess[RAG Processing]
    RAGCheck -->|No| SimpleResponse[Simple Response]
    RAGProcess --> NPCResponse[Display NPC Response]
    SimpleResponse --> NPCResponse
    
    SwitchDM --> Prompt
    ShowMemory --> Prompt
    ShowAgents --> Prompt
    ShowTools --> Prompt
    DMResponse --> Prompt
    NPCResponse --> Prompt
    
    style Start fill:#90EE90
    style Exit fill:#FFB6C1
    style DMProcess fill:#FFE4B5
    style RAGProcess fill:#E0E6FF
```

### Game Loop Explanation

The main game loop continuously processes user input and routes it to appropriate handlers. It supports various commands for navigation and debugging, and handles different agent types differently. The Dungeon Master agent uses tool detection and execution, while NPC agents may use RAG for enhanced contextual responses.

## Tool Execution System

```mermaid
sequenceDiagram
    participant User as User
    participant DM as Dungeon Master
    participant Executor as Function Executor
    participant MCP as MCP Client
    participant Agent as Target Agent
    
    User->>DM: User input with tool request
    DM->>Executor: Detect and prepare tool call
    Executor->>User: Request confirmation
    User->>Executor: Confirm execution
    
    alt speak_to_somebody tool
        Executor->>MCP: Check player room position
        MCP-->>Executor: Room validation result
        alt Same room
            Executor->>Agent: Switch to target agent
            Executor-->>DM: Success message
        else Different room
            Executor-->>DM: Cannot speak message
        end
    else MCP tool
        Executor->>MCP: Execute MCP tool
        MCP-->>Executor: Tool result
        Executor-->>DM: Tool execution result
    end
    
    DM-->>User: Display final response
```

### Tool Execution Explanation

The tool execution system provides a secure way to execute both custom tools (like `speak_to_somebody`) and MCP tools. It includes user confirmation for security, room validation for agent interaction, and proper error handling. The system ensures players can only interact with NPCs in the same room and provides appropriate feedback for all scenarios.

## RAG (Retrieval-Augmented Generation) System

```mermaid
flowchart TD
    UserInput[User Input] --> Search[Similarity Search]
    Search --> Database[(Knowledge Database)]
    Database --> Results{Found Similarities?}
    Results -->|Yes| Context[Generate Context Message]
    Results -->|No| Direct[Direct User Message]
    Context --> Combine[Combine Context + User Input]
    Direct --> Agent[Send to Agent]
    Combine --> Agent
    Agent --> Response[Generate Response]
    Response --> Display[Display to User]
    
    style Search fill:#E0E6FF
    style Database fill:#FFE4E1
    style Context fill:#F0FFF0
    style Agent fill:#FFF8DC
```

### RAG System Explanation

The RAG system enhances agent responses by searching for relevant context from a knowledge database. When a user interacts with RAG-enabled agents (Guard, Sorcerer, Merchant, Healer), the system performs similarity searches to find relevant information, then combines this context with the user's input to generate more informed and contextual responses.

## Command System

```mermaid
graph LR
    subgraph "Available Commands"
        Bye["/bye<br/>Exit Application"]
        DM["/dm, /back-to-dm<br/>Switch to Dungeon Master"]
        Memory["/memory<br/>Display Conversation History"]
        Agents["/agents<br/>List All Available Agents"]
        Tools["/tools<br/>Display Available Tools"]
        Debug["/debug<br/>Show Agent Memory<br/>(NPC agents only)"]
    end
    
    User[User Input] --> CommandParser{Command Parser}
    CommandParser --> Bye
    CommandParser --> DM
    CommandParser --> Memory
    CommandParser --> Agents
    CommandParser --> Tools
    CommandParser --> Debug
    CommandParser --> Normal[Normal Conversation]
    
    style Bye fill:#FFB6C1
    style DM fill:#FFE4B5
    style Memory fill:#E0E6FF
    style Debug fill:#FFF8DC
```

### Command System Explanation

The application provides a comprehensive command system for navigation and debugging. Users can exit the game, switch between agents, view conversation history, list available agents and tools, and access debug information. This system ensures smooth navigation and provides transparency about the application's capabilities and state.