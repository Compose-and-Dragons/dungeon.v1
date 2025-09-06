# End-of-Level Boss Agent Architecture

This document explains the architecture and flow of the End-of-Level Boss Agent application using various diagrams and explanations.

## System Architecture Overview

```mermaid
graph TB
    subgraph "Main Application Flow"
        A["main Function"] --> B[Configuration Setup]
        A --> C[OpenAI Client Setup]
        A --> D[Boss Agent Creation]
        A --> E[Agent Card Setup]
        A --> F[Streaming Callback Definition]
        A --> G[A2A Server Setup]
    end
    
    B --> B1["Environment Variables<br/>LLM_URL, MCP_HOST<br/>SIMILARITY_LIMIT, HTTP_PORT"]
    
    C --> C1["OpenAI Client<br/>Base URL, API Key"]
    
    D --> D1["agents.GetBossAgent<br/>with context and client"]
    
    E --> E1["AgentCard Configuration<br/>Name, Description, URL<br/>Skills: ask_for_something"]
    
    F --> F1[Task Request Processing]
    F1 --> F2[RAG Similarity Search]
    F2 --> F3[Streaming Response]
    
    G --> G1["A2A Server Start<br/>Port 8080 with streaming"]
    
    style A fill:#e1d5e7
    style B fill:#fff2cc
    style C fill:#dae8fc
    style D fill:#f8cecc
    style E fill:#d5e8d4
    style F fill:#e6ccff
    style G fill:#ffcccc
```

### System Architecture Explanation

This diagram shows the initialization flow of the End-of-Level Boss Agent system. The main function orchestrates seven primary components: Configuration loading from environment variables, OpenAI client setup for LLM communication, Boss agent creation using the agents package, Agent card configuration for A2A protocol, Streaming callback definition for handling requests, and A2A server setup for HTTP communication. This creates a complete Agent-to-Agent (A2A) server that can handle streaming conversations with RAG (Retrieval-Augmented Generation) capabilities.

## Agent Streaming Callback Flow

```mermaid
sequenceDiagram
    participant Client as A2A Client
    participant Server as A2A Server
    participant Callback as Stream Callback
    participant RAG as RAG Function
    participant Agent as Boss Agent
    participant Search as Similarity Search
    
    Client->>Server: POST /stream with TaskRequest
    Server->>Callback: agentStreamCallback()
    
    Callback->>Callback: Extract user message
    Callback->>Callback: Check skill metadata
    
    alt skill == "ask_for_something"
        Callback->>Callback: userPrompt = userMessage
    else invalid skill
        Callback->>Callback: userPrompt = error message
    end
    
    Callback->>RAG: GeneratePromptMessagesWithSimilarities()
    RAG->>Search: agents.SearchSimilarities()
    Search-->>RAG: similarities array
    
    alt similarities found
        RAG-->>Callback: [SystemMessage(context), UserMessage(input)]
    else no similarities
        RAG-->>Callback: [UserMessage(input)]
    end
    
    Callback->>Agent: bossAgent.RunStream()
    Agent-->>Callback: streaming content
    Callback->>Server: streamFunc(content)
    Server-->>Client: streamed response
    
    opt debug mode
        Callback->>Agent: msg.DisplayHistory()
    end
```

### Agent Streaming Callback Explanation

This sequence diagram illustrates the complete flow when a client makes a streaming request to the A2A server. The process begins with task request processing, validates the skill type, performs RAG similarity search to gather relevant context, and then streams the boss agent's response back to the client. The RAG system enhances responses by finding similar historical interactions, while the streaming mechanism provides real-time response delivery. Debug functionality allows inspection of the conversation history when requested.

## RAG Similarity Search Process

```mermaid
stateDiagram-v2
    [*] --> SearchInitiated
    
    SearchInitiated --> SearchSimilarities : Call agents.SearchSimilarities()
    
    SearchSimilarities --> ProcessingResults : Get similarities array
    
    ProcessingResults --> SimilaritiesFound : len(similarities) > 0
    ProcessingResults --> NoSimilarities : len(similarities) == 0
    
    SimilaritiesFound --> BuildContext : Create context message
    BuildContext --> CreateSystemMessage : Add system message with context
    CreateSystemMessage --> AddUserMessage : Add user message
    AddUserMessage --> ReturnMessages : Return message array
    
    NoSimilarities --> CreateUserOnly : Create user message only
    CreateUserOnly --> ReturnMessages : Return single message
    
    ReturnMessages --> [*]
```

### RAG Similarity Search Explanation

This state diagram shows how the RAG (Retrieval-Augmented Generation) system processes user input to find relevant contextual information. The system searches for similar chunks of text based on the input, applying similarity limits and maximum result constraints. When similarities are found, they are formatted into a context message that becomes part of the system prompt, enhancing the agent's response with relevant background information. This approach allows the boss agent to provide more informed and contextually appropriate responses based on previous interactions or knowledge.

## A2A Server Configuration and Startup

```mermaid
graph LR
    subgraph "Server Configuration"
        Config[Server Config]
        Port[HTTP Port 8080]
        Card[Agent Card]
        Callback[Stream Callback]
    end
    
    subgraph "Agent Card Properties"
        Name[Agent Name]
        Desc[Description]
        URL[Agent URL]
        Version[Version 1.0.0]
        Skills[Skills Array]
    end
    
    subgraph "Skills Definition"
        Skill1["ask_for_something<br/>Using small LLM<br/>to answer questions"]
    end
    
    subgraph "Server Runtime"
        HTTP[HTTP Server]
        Stream[Streaming Support]
        Endpoints[A2A Endpoints]
    end
    
    Config --> Port
    Config --> Card
    Config --> Callback
    
    Card --> Name
    Card --> Desc
    Card --> URL
    Card --> Version
    Card --> Skills
    
    Skills --> Skill1
    
    Config --> HTTP
    HTTP --> Stream
    HTTP --> Endpoints
```

### A2A Server Configuration Explanation

This diagram illustrates the configuration and structure of the A2A (Agent-to-Agent) server. The server is configured with specific properties including the HTTP port, agent card metadata, and streaming callback function. The agent card defines the boss agent's capabilities and available skills, currently supporting the "ask_for_something" skill for general question answering. The server provides streaming support for real-time response delivery and implements the A2A protocol endpoints for agent communication.

## Environment Configuration and Dependencies

```mermaid
flowchart TD
    A[Environment Variables] --> B[MODEL_RUNNER_BASE_URL]
    A --> C[MCP_HOST]
    A --> D[SIMILARITY_LIMIT]
    A --> E[SIMILARITY_MAX_RESULTS]
    A --> F[BOSS_REMOTE_AGENT_HTTP_PORT]
    
    B --> G[OpenAI Client Configuration]
    C --> H[MCP Service Connection]
    D --> I[RAG Search Threshold]
    E --> I
    F --> J[Server Port Configuration]
    
    G --> K[Boss Agent Creation]
    H --> L[External Service Integration]
    I --> M[Similarity Search Parameters]
    J --> N[HTTP Server Setup]
    
    K --> O[Agent Ready]
    L --> O
    M --> O
    N --> O
    
    O --> P[A2A Server Start]
    
    style A fill:#fff2cc
    style G fill:#dae8fc
    style H fill:#f8cecc
    style I fill:#ffe6cc
    style J fill:#ffcccc
    style P fill:#d5e8d4
```

### Environment Configuration Explanation

This flowchart demonstrates how environment variables configure the various components of the boss agent system. The configuration process reads multiple environment variables to set up the LLM connection, MCP service integration, RAG search parameters, and server settings. Each configuration aspect feeds into specific system components: the OpenAI client for LLM communication, MCP host for external service integration, similarity search parameters for RAG functionality, and HTTP port for the A2A server. This flexible configuration approach allows the system to adapt to different deployment environments and requirements.

## Task Processing and Response Generation

```mermaid
sequenceDiagram
    participant Request as Task Request
    participant Processor as Task Processor
    participant Validator as Skill Validator
    participant RAG as RAG System
    participant Agent as Boss Agent
    participant Stream as Stream Handler
    
    Request->>Processor: TaskRequest with message and metadata
    Processor->>Processor: Extract user message from Parts[0].Text
    Processor->>Validator: Check metadata["skill"]
    
    alt skill == "ask_for_something"
        Validator-->>Processor: Valid skill, use userMessage
    else invalid skill
        Validator-->>Processor: Invalid skill, generate error prompt
    end
    
    Processor->>RAG: GeneratePromptMessagesWithSimilarities()
    Note over RAG: Search for similar chunks<br/>Apply similarity limits<br/>Build context messages
    RAG-->>Processor: Enhanced message array
    
    Processor->>Agent: RunStream with enhanced messages
    Agent->>Stream: Stream response chunks
    
    loop For each content chunk
        Stream->>Stream: Print to console (debug)
        Stream->>Request: streamFunc(content) to client
    end
    
    opt userPrompt starts with "/debug"
        Agent->>Agent: msg.DisplayHistory()
    end
```

### Task Processing and Response Generation Explanation

This sequence diagram shows the detailed process of handling incoming task requests and generating streaming responses. The system extracts and validates the user message, checks skill metadata to ensure the request is valid, enhances the prompt using RAG similarity search, and then generates a streaming response through the boss agent. The streaming mechanism allows real-time response delivery while maintaining debug capabilities for development and troubleshooting. Error handling ensures that invalid skills receive appropriate error messages rather than causing system failures.