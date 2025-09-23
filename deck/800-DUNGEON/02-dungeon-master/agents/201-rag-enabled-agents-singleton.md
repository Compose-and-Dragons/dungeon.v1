# RAG-Enabled Agents - Singleton Pattern

⬅️ **Back to:** [RAG-Enabled Agents Schema](200-rag-enabled-agents-schema.md)

## Singleton Pattern Implementation

All RAG-enabled agents (Guard, Sorcerer, Merchant, Healer) use the singleton pattern for efficient resource management and thread-safe initialization.

```mermaid
flowchart TD
    Singleton[Singleton Pattern]:::singleton
    Singleton --> GuardAgent["Guard Agent (Huey)<br/><a href='/dungeon-master/agents/guard.agent.go#L20'>guard.agent.go:20</a>"]:::agent
    Singleton --> SorcererAgent["Sorcerer Agent<br/><a href='/dungeon-master/agents/sorcerer.agent.go#L20'>sorcerer.agent.go:20</a>"]:::agent
    Singleton --> MerchantAgent["Merchant Agent<br/><a href='/dungeon-master/agents/merchant.agent.go#L20'>merchant.agent.go:20</a>"]:::agent
    Singleton --> HealerAgent["Healer Agent<br/><a href='/dungeon-master/agents/healer.agent.go#L20'>healer.agent.go:20</a>"]:::agent

    classDef singleton fill:#f3e5f5,stroke:#4a148c,stroke-width:2px,color:#000
    classDef agent fill:#f9fbe7,stroke:#827717,stroke-width:2px,color:#000
```

### Code Structure

```go
var (
    agentInstance mu.Agent
    agentOnce     sync.Once
)

func GetAgent(ctx context.Context, client openai.Client) mu.Agent {
    agentOnce.Do(func() {
        agentInstance = createAgent(ctx, client)
    })
    return agentInstance
}
```

### Implementation Details

#### Global Variables
- **`agentInstance`**: Stores the single instance of the agent
- **`agentOnce`**: `sync.Once` ensures initialization happens only once

#### GetAgent Function
- **Thread Safety**: Uses `sync.Once` for safe concurrent access
- **Lazy Loading**: Agent created only when first requested
- **Context Preservation**: Creation context maintained throughout lifecycle

### Benefits

#### Resource Efficiency
- **Single Instance**: One agent instance per type across the application
- **Memory Management**: Avoids duplicate embeddings and model loading
- **Initialization Cost**: One-time setup per agent type

#### Performance Optimizations
- **Concurrent Safety**: Multiple goroutines can safely call `GetAgent`
- **Initialization Guarantee**: Agent creation happens exactly once
- **Resource Conservation**: No unnecessary re-initialization

### Agent Instances

The singleton pattern is implemented in:
- [Guard Agent (Huey)](guard.agent.go:20)
- [Sorcerer Agent](sorcerer.agent.go:20)
- [Merchant Agent](merchant.agent.go:20)
- [Healer Agent](healer.agent.go:20)

---

➡️ **Next:** [Configuration](202-rag-enabled-agents-configuration.md)