# Vector Stores - Global Agent Stores

⬅️ **Back to:** [Vector Stores Schema](100-stores-schema.md)

## Global Agent Stores

The RAG system uses a global registry to manage vector stores for all agents through a centralized map structure.

```mermaid
flowchart TD
    GlobalStore[Global Agent Stores]:::global
    GlobalStore --> AgentsStores("AgentsStores Map<br/><a href='/dungeon-master/agents/stores.go#L17'>stores.go:17</a>"):::map
    AgentsStores --> StoreStructure["map[string]rag.MemoryVectorStore"]:::structure

    StoreStructure --> AgentKeys["Agent Names<br/>(huey, sorcerer, merchant, healer)"]:::keys
    StoreStructure --> VectorStores["MemoryVectorStore<br/>Records: map[string]rag.VectorRecord"]:::stores

    classDef global fill:#f3e5f5,stroke:#4a148c,stroke-width:2px,color:#000
    classDef map fill:#fce4ec,stroke:#880e4f,stroke-width:2px,color:#000
    classDef structure fill:#e0f2f1,stroke:#00695c,stroke-width:2px,color:#000
    classDef keys fill:#fff3e0,stroke:#e65100,stroke-width:2px,color:#000
    classDef stores fill:#e8f5e8,stroke:#1b5e20,stroke-width:2px,color:#000
```

## Implementation Details

### Global Variable Declaration
```go
var AgentsStores = make(map[string]rag.MemoryVectorStore)
```

### Structure and Purpose

#### Map Structure
- **Key**: Agent name (string) - lowercase agent identifier
- **Value**: `rag.MemoryVectorStore` containing embedded knowledge
- **Scope**: Package-level variable accessible to all functions

#### Agent Registration
Each agent gets its own vector store entry:
- **"huey"**: Guard agent's knowledge base
- **"sorcerer"**: Sorcerer agent's magical knowledge
- **"merchant"**: Merchant agent's trading information
- **"healer"**: Healer agent's medical knowledge

#### Memory Vector Store Structure
```go
rag.MemoryVectorStore{
    Records: map[string]rag.VectorRecord{
        // Vector records with embeddings and content
    }
}
```

### Usage Pattern

#### Store Initialization
```go
AgentsStores[name] = rag.MemoryVectorStore{
    Records: make(map[string]rag.VectorRecord),
}
```

#### Store Retrieval
```go
store := AgentsStores[agentName]
```

#### Store Persistence
Each store can be saved to JSON files for persistence between application runs.

---

➡️ **Next:** [Embedding Generation](102-stores-embedding-generation.md)