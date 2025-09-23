# Boss Agent Implementation

⬅️ **Back to:** [Boss Agent Schema](100-boss-agent-schema.md)

## Boss Agent Overview

The Boss Agent (`agents/boss.agent.go`) implements a singleton-pattern agent that serves as the end-of-level boss character. It integrates LLM capabilities, RAG (Retrieval-Augmented Generation), and fallback mechanisms.

```mermaid
flowchart TD
    BossAgentFile[boss.agent.go<br/><a href='/dungeon-end-of-level-boss/agents/boss.agent.go#L14'>boss.agent.go:14</a>]:::main
    BossAgentFile --> Constructor[Agent Constructor<br/><a href='/dungeon-end-of-level-boss/agents/boss.agent.go#L28'>boss.agent.go:28-78</a]:::constructor
    BossAgentFile --> RAGInit[RAG Initialization]:::rag
    BossAgentFile --> Fallback[Fallback Mechanism]:::fallback


    Constructor --> EnvVars["Environment Variables<br/><a href='/dungeon-end-of-level-boss/agents/boss.agent.go#L30'>boss.agent.go:30-33</a>"]:::env
    Constructor --> SystemInstructions["System Instructions Loading<br/><a href='/dungeon-end-of-level-boss/agents/boss.agent.go#L46'>boss.agent.go:46-59</a>"]:::instructions

    RAGInit --> EmbeddingGen["GenerateEmbeddings Call<br/><a href='/dungeon-end-of-level-boss/agents/boss.agent.go#L36'>boss.agent.go:36</a>"]:::embedding

    Fallback --> GhostAgent["NewGhostAgent Fallback<br/><a href='/dungeon-end-of-level-boss/agents/boss.agent.go#L73'>boss.agent.go:73</a>"]:::ghost

    classDef main fill:#e1f5fe,stroke:#01579b,stroke-width:3px,color:#000
    classDef singleton fill:#f3e5f5,stroke:#4a148c,stroke-width:2px,color:#000
    classDef constructor fill:#e8f5e8,stroke:#1b5e20,stroke-width:2px,color:#000
    classDef config fill:#fff3e0,stroke:#e65100,stroke-width:2px,color:#000
    classDef rag fill:#fce4ec,stroke:#880e4f,stroke-width:2px,color:#000
    classDef fallback fill:#e0f2f1,stroke:#00695c,stroke-width:2px,color:#000
    classDef vars fill:#f9fbe7,stroke:#827717,stroke-width:2px,color:#000
    classDef getter fill:#e8eaf6,stroke:#283593,stroke-width:2px,color:#000
    classDef create fill:#e4f7ff,stroke:#0277bd,stroke-width:2px,color:#000
    classDef agent fill:#f1f8e9,stroke:#33691e,stroke-width:2px,color:#000
    classDef env fill:#fff8e1,stroke:#f57f17,stroke-width:2px,color:#000
    classDef instructions fill:#c8e6c9,stroke:#388e3c,stroke-width:2px,color:#000
    classDef embedding fill:#ffcdd2,stroke:#d32f2f,stroke-width:2px,color:#000
    classDef error fill:#ffebee,stroke:#c62828,stroke-width:2px,color:#000
    classDef ghost fill:#e1f5fe,stroke:#0288d1,stroke-width:2px,color:#000
```

## Singleton Pattern Implementation

### Global Variables
```go
var (
    sorcererAgentInstance mu.Agent
    sorcererAgentOnce     sync.Once
)
```

**Design Notes**:
- **Variable Naming**: Uses `sorcererAgent` naming (legacy from original wizard boss)
- **Thread Safety**: `sync.Once` ensures single initialization in concurrent environment
- **Instance Storage**: Global variable holds the singleton instance

### GetBossAgent Function
```go
func GetBossAgent(ctx context.Context, client openai.Client) mu.Agent {
    sorcererAgentOnce.Do(func() {
        sorcererAgentInstance = createBossAgent(ctx, client)
    })
    return sorcererAgentInstance
}
```

**Function Behavior**:
- **Lazy Initialization**: Agent created only when first requested
- **Thread-Safe**: `sync.Once` guarantees single execution
- **Context Passing**: Forwards context for dependent operations
- **Client Injection**: OpenAI client dependency injection

## Environment Configuration

The agent constructor reads configuration from environment variables:

```mermaid
flowchart LR
    EnvConfig[Environment Configuration]:::main
    EnvConfig --> BossName["BOSS_NAME<br/>Default: 'Louie'"]:::env
    EnvConfig --> BossDesc["BOSS_DESCRIPTION<br/>Default: 'A wise and powerful Sphinx in a fantasy world.'"]:::env
    EnvConfig --> BossModel["BOSS_MODEL<br/>Default: 'ai/qwen2.5:1.5B-F16'"]:::env
    EnvConfig --> Temperature["BOSS_MODEL_TEMPERATURE<br/>Default: '0.0'"]:::env
    EnvConfig --> ContextPath["BOSS_CONTEXT_PATH<br/>Default: ''"]:::env
    EnvConfig --> InstructionsPath["BOSS_SYSTEM_INSTRUCTIONS_PATH<br/>Default: ''"]:::env

    classDef main fill:#e1f5fe,stroke:#01579b,stroke-width:3px,color:#000
    classDef env fill:#f3e5f5,stroke:#4a148c,stroke-width:2px,color:#000
```

### Configuration Loading
```go
name := helpers.GetEnvOrDefault("BOSS_NAME", "Louie")
description := helpers.GetEnvOrDefault("BOSS_DESCRIPTION", "A wise and powerful Sphinx in a fantasy world.")
model := helpers.GetEnvOrDefault("BOSS_MODEL", "ai/qwen2.5:1.5B-F16")
temperature := helpers.StringToFloat(helpers.GetEnvOrDefault("BOSS_MODEL_TEMPERATURE", "0.0"))
```

**Configuration Features**:
- **Fallback Values**: Each variable has sensible defaults
- **Type Conversion**: Automatic string-to-float conversion for temperature
- **Character Setup**: Name and description define the boss personality
- **Model Selection**: Configurable LLM model for different capabilities

## RAG System Integration

### Embeddings Generation
```go
errEmbedding := GenerateEmbeddings(ctx, &client, name, helpers.GetEnvOrDefault("BOSS_CONTEXT_PATH", ""))
if errEmbedding != nil {
    fmt.Println("🔶 Error generating embeddings for sorcerer agent:", errEmbedding)
}
```

**RAG Setup Process**:
1. **Context File**: Reads from `BOSS_CONTEXT_PATH` environment variable
2. **Embedding Generation**: Creates vector embeddings for context chunks
3. **Error Handling**: Non-fatal errors - agent continues without RAG if embedding fails
4. **Agent Storage**: Associates embeddings with the boss agent name

### Vector Store Benefits
- **Context Awareness**: Boss can reference relevant lore and game context
- **Enhanced Responses**: Similarity search provides relevant background information
- **Configurable Content**: External files allow easy content updates
- **Performance**: Pre-computed embeddings for fast similarity search

## System Instructions Management

```mermaid
flowchart TD
    SystemInstructions[System Instructions]:::main
    SystemInstructions --> PathCheck[Check Instructions Path]:::check
    SystemInstructions --> FileRead[Read Instructions File]:::read
    SystemInstructions --> FallbackHandling[Handle Fallback]:::fallback

    PathCheck --> PathEmpty["Path Empty?<br/><a href='/dungeon-end-of-level-boss/agents/boss.agent.go#L47-L49'>boss.agent.go:47-49</a>"]:::pathcheck
    PathCheck --> PathProvided["Path Provided<br/><a href='/dungeon-end-of-level-boss/agents/boss.agent.go#L52'>boss.agent.go:52</a>"]:::pathprovided

    FileRead --> ReadSuccess["File Read Successfully<br/><a href='/dungeon-end-of-level-boss/agents/boss.agent.go#L58'>boss.agent.go:58</a>"]:::success
    FileRead --> ReadError["File Read Error<br/><a href='/dungeon-end-of-level-boss/agents/boss.agent.go#L54-L56'>boss.agent.go:54-56</a>"]:::error

    FallbackHandling --> DefaultInstructions["Default Sphinx Instructions<br/><a href='/dungeon-end-of-level-boss/agents/boss.agent.go#L56'>boss.agent.go:56</a>"]:::default

    classDef main fill:#e1f5fe,stroke:#01579b,stroke-width:3px,color:#000
    classDef check fill:#f3e5f5,stroke:#4a148c,stroke-width:2px,color:#000
    classDef read fill:#e8f5e8,stroke:#1b5e20,stroke-width:2px,color:#000
    classDef fallback fill:#fff3e0,stroke:#e65100,stroke-width:2px,color:#000
    classDef pathcheck fill:#fce4ec,stroke:#880e4f,stroke-width:2px,color:#000
    classDef pathprovided fill:#e0f2f1,stroke:#00695c,stroke-width:2px,color:#000
    classDef success fill:#f9fbe7,stroke:#827717,stroke-width:2px,color:#000
    classDef error fill:#e8eaf6,stroke:#283593,stroke-width:2px,color:#000
    classDef default fill:#e4f7ff,stroke:#0277bd,stroke-width:2px,color:#000
```

### Instructions Loading Logic
```go
systemInstructionsContentPath := helpers.GetEnvOrDefault("BOSS_SYSTEM_INSTRUCTIONS_PATH", "")
if systemInstructionsContentPath == "" {
    fmt.Println("🔶 No BOSS_SYSTEM_INSTRUCTIONS_PATH provided, using default instructions.")
}

systemInstructionsContent, err := helpers.ReadTextFile(systemInstructionsContentPath)

if err != nil {
    fmt.Println("🔶 Error reading the file, using default instructions:", err)
    systemInstructions = openai.SystemMessage("You are a wise and powerful Sphinx in a fantasy world.")
} else {
    systemInstructions = openai.SystemMessage(systemInstructionsContent)
}
```

**Instruction Handling**:
- **File-based**: Allows external configuration of boss personality
- **Graceful Fallback**: Uses default Sphinx character if file unavailable
- **Non-blocking**: File errors don't prevent agent creation
- **Flexible**: Easy to update boss behavior without code changes

## Agent Creation Process

### LLM Agent Construction
```go
chatAgent, err := mu.NewAgentWithDescription(ctx, name, description,
    mu.WithClient(client),
    mu.WithParams(openai.ChatCompletionNewParams{
        Model:       model,
        Temperature: openai.Opt(temperature),
        Messages: []openai.ChatCompletionMessageParamUnion{
            systemInstructions,
        },
    }),
)
```

**Agent Parameters**:
- **Name & Description**: Character identity from environment
- **OpenAI Client**: Injected client for LLM communication
- **Model Configuration**: Configurable model and temperature
- **System Message**: Instructions for character behavior

### Agent Interface Compliance
The created agent implements the `mu.Agent` interface, providing:
- **RunStream**: Streaming response generation
- **Message Management**: Conversation history handling
- **Embedding Generation**: Vector embedding capabilities
- **Metadata Management**: Agent configuration and state

## Fallback Mechanism

### Error Handling Strategy
```go
if err != nil {
    fmt.Println("🔶 Error creating boss agent, creating ghost agent instead:", err)
    return NewGhostAgent("[Ghost] " + name)
}
```

**Fallback Benefits**:
- **Service Reliability**: Agent creation never fails completely
- **Graceful Degradation**: Ghost agent provides basic functionality
- **Debug Information**: Clear error logging for troubleshooting
- **Name Preservation**: Ghost agent maintains original name with prefix

### Ghost Agent Integration
The Ghost Agent serves as a lightweight fallback that:
- **Implements mu.Agent**: Same interface as full LLM agent
- **Simulates Responses**: Provides pre-defined responses for testing
- **Maintains Service**: Keeps the boss service operational
- **Debugging Aid**: Helps identify configuration issues

---

⬅️ **Back to:** [Boss Agent Schema](100-boss-agent-schema.md)