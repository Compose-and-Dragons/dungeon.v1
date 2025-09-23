# Vector Stores - Embedding Generation

⬅️ **Back to:** [Global Agent Stores](101-stores-global-stores.md)

## Embedding Generation Process

The `GenerateEmbeddings` function creates and manages vector embeddings for agent knowledge bases, with intelligent caching and persistence.

```mermaid
flowchart TD
    EmbeddingGeneration[Embedding Generation]:::embedding
    EmbeddingGeneration --> GenerateEmbeddings("GenerateEmbeddings Function<br/><a href='/dungeon-master/agents/stores.go#L21'>stores.go:21</a>"):::function

    GenerateEmbeddings --> StoreInitialization("Vector Store Initialization<br/><a href='/dungeon-master/agents/stores.go#L24'>stores.go:24</a>"):::init
    GenerateEmbeddings --> FileLoading("Load Existing Store<br/><a href='/dungeon-master/agents/stores.go#L36'>stores.go:36</a>"):::load
    GenerateEmbeddings --> FreshCreation("Create Fresh Store<br/><a href='/dungeon-master/agents/stores.go#L41'>stores.go:41</a>"):::fresh

    FileLoading --> JSONPath("JSON Store File Path<br/><a href='/dungeon-master/agents/stores.go#L30'>stores.go:30</a>"):::path
    FileLoading --> LoadResult{File Exists?}:::decision
    LoadResult -->|Yes| LoadSuccess("Load Successful<br/><a href='/dungeon-master/agents/stores.go#L120'>stores.go:120</a>"):::success
    LoadResult -->|No| CreateNew("Create New Store<br/><a href='/dungeon-master/agents/stores.go#L42'>stores.go:42</a>"):::create

    FreshCreation --> EmbeddingAgent("Create Embedding Agent<br/><a href='/dungeon-master/agents/stores.go#L49'>stores.go:49</a>"):::agent
    FreshCreation --> ContentReading("Read Context File<br/><a href='/dungeon-master/agents/stores.go#L70'>stores.go:70</a>"):::read
    FreshCreation --> ChunkSplitting("Split into Chunks<br/><a href='/dungeon-master/agents/stores.go#L77'>stores.go:77</a>"):::chunks
    FreshCreation --> EmbeddingLoop("Generate Embeddings Loop<br/><a href='/dungeon-master/agents/stores.go#L79'>stores.go:79</a>"):::loop
    FreshCreation --> StorePersistence("Persist Store to File<br/><a href='/dungeon-master/agents/stores.go#L99'>stores.go:99</a>"):::persist

    classDef embedding fill:#fff3e0,stroke:#e65100,stroke-width:2px,color:#000
    classDef function fill:#f9fbe7,stroke:#827717,stroke-width:2px,color:#000
    classDef init fill:#e8eaf6,stroke:#283593,stroke-width:2px,color:#000
    classDef load fill:#ffebee,stroke:#c62828,stroke-width:2px,color:#000
    classDef fresh fill:#e4f7ff,stroke:#0277bd,stroke-width:2px,color:#000
    classDef path fill:#fff8e1,stroke:#f57f17,stroke-width:2px,color:#000
    classDef decision fill:#ffecb3,stroke:#f57c00,stroke-width:2px,color:#000
    classDef success fill:#c8e6c9,stroke:#388e3c,stroke-width:2px,color:#000
    classDef create fill:#ffcdd2,stroke:#d32f2f,stroke-width:2px,color:#000
    classDef agent fill:#e1f5fe,stroke:#0288d1,stroke-width:2px,color:#000
    classDef read fill:#f3e5f5,stroke:#4a148c,stroke-width:2px,color:#000
    classDef chunks fill:#e0f2f1,stroke:#00695c,stroke-width:2px,color:#000
    classDef loop fill:#fff3e0,stroke:#e65100,stroke-width:2px,color:#000
    classDef persist fill:#e8f5e8,stroke:#1b5e20,stroke-width:2px,color:#000
```

## Function Signature

```go
func GenerateEmbeddings(ctx context.Context, client *openai.Client, name string, contextInstructionsContentPath string) error
```

### Parameters
- **ctx**: Context for operation control
- **client**: OpenAI client for embedding generation
- **name**: Agent name (used for store identification)
- **contextInstructionsContentPath**: Path to agent's knowledge base file

## Process Flow

### 1. Store Initialization
```go
AgentsStores[name] = rag.MemoryVectorStore{
    Records: make(map[string]rag.VectorRecord),
}
store := AgentsStores[name]
```

### 2. File Path Construction
```go
jsonStoreFilePath := helpers.GetEnvOrDefault("VECTOR_STORES_PATH", "./data") + "/" +
    strings.ToLower(name) + "_vector_store.json"
```

**Path Pattern**: `{VECTOR_STORES_PATH}/{agent_name}_vector_store.json`

### 3. Existing Store Loading
If a persisted store exists, it's loaded directly:
```go
if err == nil {
    fmt.Println("✅ Vector store loaded successfully with", len(store.Records), "records")
    return nil
}
```

### 4. Fresh Store Creation

#### Embedding Agent Setup
```go
embeddingAgent, err := mu.NewAgent(ctx, "vector-agent",
    mu.WithClient(*client),
    mu.WithEmbeddingParams(
        openai.EmbeddingNewParams{
            Model: helpers.GetEnvOrDefault("EMBEDDING_MODEL", "ai/mxbai-embed-large:latest"),
        },
    ),
)
```

#### Content Processing
```go
contextInstructionsContent, err := helpers.ReadTextFile(contextInstructionsContentPath)
chunks := rag.SplitMarkdownBySections(contextInstructionsContent)
```

#### Embedding Generation Loop
```go
for idx, chunk := range chunks {
    embeddingVector, err := embeddingAgent.GenerateEmbeddingVector(chunk)
    if err != nil {
        return err
    }
    _, errSave := store.Save(rag.VectorRecord{
        Prompt:    chunk,
        Embedding: embeddingVector,
    })
    if errSave != nil {
        return errSave
    }
}
```

#### Store Persistence
```go
err = store.Persist(jsonStoreFilePath)
if err != nil {
    fmt.Println("🔶 Error saving vector store:", err)
    return err
}
```

## Configuration

### Environment Variables
- **`VECTOR_STORES_PATH`**: Directory for storing vector databases (default: "./data")
- **`EMBEDDING_MODEL`**: Model for generating embeddings (default: "ai/mxbai-embed-large:latest")

---

➡️ **Next:** [Similarity Search](103-stores-similarity-search.md)