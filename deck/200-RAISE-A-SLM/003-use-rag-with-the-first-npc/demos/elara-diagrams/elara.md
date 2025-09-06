# Elara Sorcerer Agent Architecture

This document explains the architecture and flow of the Elara Sorcerer Agent application using various diagrams and explanations.

## System Architecture Overview

```mermaid
graph TB
    subgraph "Main Application Flow"
        A["main Function"] --> B[Initialization]
        B --> C[OpenAI Client Setup]
        B --> D[Sorcerer Agent Creation]
        A --> E[Interactive Chat Loop]
    end
    
    B --> B1["Environment Configuration<br/>LLM_URL, SIMILARITY_LIMIT<br/>SIMILARITY_MAX_RESULTS"]
    
    C --> C1["OpenAI Client<br/>Base URL, Empty API Key"]
    
    D --> D1["agents.GetSorcererAgent<br/>with context and client"]
    
    E --> E1[User Input Processing]
    E1 --> E2[Command Processing]
    E2 --> E3[RAG Enhancement]
    E3 --> E4[Agent Response Generation]
    E4 --> E1
    
    style A fill:#e1d5e7
    style B fill:#fff2cc
    style C fill:#dae8fc
    style D fill:#f8cecc
    style E fill:#d5e8d4
```

### System Architecture Explanation

This diagram shows the architecture of the Elara Sorcerer Agent, an interactive chat application with RAG capabilities. The system initializes with environment configuration, creates an OpenAI client and sorcerer agent, then enters an infinite chat loop. The loop processes user input, handles special commands, enhances prompts with RAG similarity search, and generates streaming responses. This creates a conversational AI experience where the agent can reference previous interactions and contextual information to provide more informed responses.

## Interactive Chat Loop Flow

```mermaid
stateDiagram-v2
    [*] --> PromptUser
    
    PromptUser --> ProcessInput : User enters input
    
    ProcessInput --> CheckByeCommand : Parse input
    CheckByeCommand --> ExitApplication : input starts with "/bye"
    CheckByeCommand --> CheckMemoryCommand : not "/bye"
    
    CheckMemoryCommand --> DisplayHistory : input starts with "/memory"
    CheckMemoryCommand --> ProcessChat : regular input
    
    DisplayHistory --> PromptUser : Continue loop
    
    ProcessChat --> RAGProcessing : Enhance with similarities
    RAGProcessing --> AgentResponse : Generate response
    AgentResponse --> PromptUser : Continue loop
    
    ExitApplication --> [*] : Print goodbye message
```

### Interactive Chat Loop Explanation

This state diagram illustrates the main interactive loop of the Elara application. The system continuously prompts the user for input, then processes it through different paths based on the content. Special commands like "/bye" exit the application, "/memory" displays conversation history, while regular input goes through RAG processing and agent response generation. The loop structure ensures continuous interaction until the user explicitly exits, providing a persistent conversational experience with the sorcerer agent.

## RAG Similarity Search Process

```mermaid
sequenceDiagram
    participant User as User Input
    participant Main as Main Loop
    participant RAG as RAG Function
    participant Search as Similarity Search
    participant Agent as Sorcerer Agent
    participant Stream as Stream Handler
    
    User->>Main: Enter chat message
    Main->>RAG: GeneratePromptMessagesWithSimilarities()
    RAG->>Search: agents.SearchSimilarities()
    
    alt Similarities found
        Search-->>RAG: similarity results array
        RAG->>RAG: Build context message from similarities
        RAG-->>Main: [SystemMessage(context), UserMessage(input)]
    else No similarities
        Search-->>RAG: empty results
        RAG-->>Main: [UserMessage(input)]
    end
    
    Main->>Agent: RunStream with enhanced messages
    Agent->>Stream: Stream response chunks
    
    loop For each content chunk
        Stream->>Main: Print content to console
    end
    
    Main->>Main: Add newlines and continue loop
```

### RAG Similarity Search Explanation

This sequence diagram shows how the RAG (Retrieval-Augmented Generation) system enhances user prompts with contextual information. When a user enters a message, the system searches for similar previous interactions or knowledge chunks. If similarities are found above the configured threshold, they are formatted into a context message that becomes part of the system prompt. This enhanced prompt is then sent to the sorcerer agent, which can reference the contextual information to provide more informed and relevant responses.

## Command Processing and Routing

```mermaid
flowchart TD
    A[User Input] --> B{Check Input Prefix}
    
    B -->|"/bye"| C[Exit Command]
    B -->|"/memory"| D[Memory Command]
    B -->|Regular Input| E[Chat Processing]
    
    C --> C1[Print Goodbye Message]
    C1 --> C2[Break from Loop]
    C2 --> F[Application Exit]
    
    D --> D1[Call msg.DisplayHistory]
    D1 --> D2[Show Agent Memory]
    D2 --> G[Continue Loop]
    
    E --> E1[RAG Processing]
    E1 --> E2[Agent Response]
    E2 --> E3[Stream to Console]
    E3 --> G
    
    G --> A
    
    style C fill:#ffcccc
    style D fill:#e6ccff
    style E fill:#d5e8d4
    style F fill:#f8d7da
```

### Command Processing Explanation

This flowchart demonstrates how the application routes different types of user input. The system uses prefix matching to identify special commands: "/bye" triggers application exit, "/memory" displays the agent's conversation history, while regular input goes through the full RAG and response generation pipeline. This command system provides users with control over the application state and debugging capabilities while maintaining the primary chat functionality.

## Environment Configuration and Setup

```mermaid
graph LR
    subgraph "Environment Variables"
        ENV1[MODEL_RUNNER_BASE_URL]
        ENV2[SIMILARITY_LIMIT]
        ENV3[SIMILARITY_MAX_RESULTS]
    end
    
    subgraph "Configuration Processing"
        PROC1[llmURL Configuration]
        PROC2[RAG Parameters]
    end
    
    subgraph "System Components"
        COMP1[OpenAI Client]
        COMP2[Sorcerer Agent]
        COMP3[RAG System]
    end
    
    ENV1 --> PROC1
    ENV2 --> PROC2
    ENV3 --> PROC2
    
    PROC1 --> COMP1
    PROC1 --> COMP2
    PROC2 --> COMP3
    
    COMP1 --> COMP2
    COMP2 --> COMP3
    
    style ENV1 fill:#fff2cc
    style ENV2 fill:#fff2cc
    style ENV3 fill:#fff2cc
    style COMP1 fill:#dae8fc
    style COMP2 fill:#f8cecc
    style COMP3 fill:#ffe6cc
```

### Environment Configuration Explanation

This diagram shows how environment variables configure the various system components. The MODEL_RUNNER_BASE_URL sets the LLM endpoint for both the OpenAI client and sorcerer agent creation. SIMILARITY_LIMIT and SIMILARITY_MAX_RESULTS control the RAG system's behavior, determining the threshold for similarity matching and the maximum number of similar chunks to include in the context. This configuration approach allows the system to adapt to different deployment environments and fine-tune the RAG system's performance.

## Agent Memory and Context Management

```mermaid
sequenceDiagram
    participant User as User
    participant UI as User Interface
    participant Agent as Sorcerer Agent
    participant Memory as Agent Memory
    participant RAG as RAG System
    
    User->>UI: Enter message
    UI->>RAG: Search for similarities
    RAG-->>UI: Enhanced messages with context
    
    UI->>Agent: RunStream with enhanced messages
    Agent->>Memory: Store conversation turn
    Agent->>UI: Stream response
    
    opt User enters "/memory"
        User->>UI: "/memory" command
        UI->>Memory: msg.DisplayHistory(agent)
        Memory-->>UI: Formatted conversation history
        UI-->>User: Display history
    end
    
    Note over Memory: Persistent across chat turns
    Note over RAG: Searches historical context
```

### Agent Memory and Context Management Explanation

This sequence diagram illustrates how the sorcerer agent manages conversation memory and context. Each interaction through RunStream automatically adds the conversation turn to the agent's memory, creating a persistent history across chat sessions. The RAG system can search through this historical context to find relevant information for enhancing responses. Users can inspect this memory using the "/memory" command, providing transparency into what the agent remembers from previous conversations.