# Boss A2A Server Implementation

⬅️ **Back to:** [Boss Agent Schema](100-boss-agent-schema.md)

## A2A Server Overview

The A2A (Agent-to-Agent) server implementation enables the End-of-Level Boss service to communicate with other game services through a standardized HTTP protocol with streaming capabilities.

```mermaid
flowchart TD
    A2AServer[A2A Server with Streaming]:::main
    A2AServer --> AgentCard[Agent Card Registration <br/><a href='/dungeon-end-of-level-boss/main.go#L39'>main.go:39</a>]:::card
    A2AServer --> StreamingCallback[Streaming Callback Handler]:::callback

    StreamingCallback --> RequestExtraction["Task Request Parsing<br/><a href='/dungeon-end-of-level-boss/main.go#L57'>main.go:57</a>"]:::extract
    StreamingCallback --> SkillRouting["Skill-based Routing<br/><a href='/dungeon-end-of-level-boss/main.go#L65'>main.go:65</a>"]:::routing
    StreamingCallback --> RAGIntegration["RAG Similarity Search<br/><a href='/dungeon-end-of-level-boss/main.go#L78'>main.go:78</a>"]:::rag
    StreamingCallback --> AgentExecution["Boss Agent Stream Execution<br/><a href='/dungeon-end-of-level-boss/main.go#L85'>main.go:85</a>"]:::execution



    classDef main fill:#e1f5fe,stroke:#01579b,stroke-width:3px,color:#000
    classDef card fill:#f3e5f5,stroke:#4a148c,stroke-width:2px,color:#000
    classDef callback fill:#e8f5e8,stroke:#1b5e20,stroke-width:2px,color:#000
    classDef init fill:#fff3e0,stroke:#e65100,stroke-width:2px,color:#000
    classDef pipeline fill:#fce4ec,stroke:#880e4f,stroke-width:2px,color:#000
    classDef cardinfo fill:#e0f2f1,stroke:#00695c,stroke-width:2px,color:#000
    classDef skills fill:#f9fbe7,stroke:#827717,stroke-width:2px,color:#000
    classDef extract fill:#e8eaf6,stroke:#283593,stroke-width:2px,color:#000
    classDef routing fill:#e4f7ff,stroke:#0277bd,stroke-width:2px,color:#000
    classDef rag fill:#f1f8e9,stroke:#33691e,stroke-width:2px,color:#000
    classDef execution fill:#fff8e1,stroke:#f57f17,stroke-width:2px,color:#000
    classDef port fill:#c8e6c9,stroke:#388e3c,stroke-width:2px,color:#000
    classDef start fill:#ffcdd2,stroke:#d32f2f,stroke-width:2px,color:#000
    classDef protocol fill:#ffebee,stroke:#c62828,stroke-width:2px,color:#000
    classDef realtime fill:#e1f5fe,stroke:#0288d1,stroke-width:2px,color:#000
```

## Agent Card Configuration

The Agent Card serves as the service's public interface definition, advertising capabilities to other services.

### Card Structure
```go
agentCard := a2a.AgentCard{
    Name:        bossAgent.GetName(),
    Description: bossAgent.GetDescription(),
    URL:         "http://localhost:" + httpPort,
    Version:     "1.0.0",
    Skills: []map[string]any{
        {
            "id":          "ask_for_something",
            "name":        "Ask for something",
            "description": bossAgent.GetName() + " is using a small language model to answer questions",
        },
    },
}
```

### Service Metadata
```mermaid
flowchart LR
    ServiceMeta[Service Metadata]:::main
    ServiceMeta --> ServiceName["Name<br/>From boss agent configuration"]:::field
    ServiceMeta --> ServiceDesc["Description<br/>Boss agent description"]:::field
    ServiceMeta --> ServiceURL["URL<br/>http://localhost:{port}"]:::field
    ServiceMeta --> ServiceVersion["Version<br/>1.0.0"]:::field

    ServiceMeta --> SkillsArray[Skills Array]:::skills
    SkillsArray --> AskSkill["ask_for_something<br/>General interaction capability"]:::skill

    classDef main fill:#e1f5fe,stroke:#01579b,stroke-width:3px,color:#000
    classDef field fill:#f3e5f5,stroke:#4a148c,stroke-width:2px,color:#000
    classDef skills fill:#e8f5e8,stroke:#1b5e20,stroke-width:2px,color:#000
    classDef skill fill:#fff3e0,stroke:#e65100,stroke-width:2px,color:#000
```

**Metadata Properties**:
- **Dynamic Name**: Uses boss agent's configured name
- **Dynamic Description**: Reflects agent's character description
- **Service URL**: Constructed from configured HTTP port
- **Version**: Static version for API compatibility
- **Skills**: Array of available capabilities

## Streaming Callback Implementation

The streaming callback function handles incoming requests and orchestrates the response generation process.

### Request Processing Pipeline

```mermaid
flowchart TD
    IncomingRequest[Incoming A2A Request]:::input
    IncomingRequest --> TaskParsing[Task Request Parsing]:::parse
    TaskParsing --> MessageExtraction[User Message Extraction]:::extract
    MessageExtraction --> SkillValidation[Skill Validation]:::validate
    SkillValidation --> ValidSkill[Valid Skill: ask_for_something]:::valid
    SkillValidation --> InvalidSkill[Invalid Skill]:::invalid
    ValidSkill --> RAGSearch[RAG Similarity Search]:::rag
    InvalidSkill --> ErrorResponse[Error Response Generation]:::error
    RAGSearch --> ContextEnhancement[Context Enhancement]:::context
    ContextEnhancement --> AgentStream[Boss Agent Streaming]:::stream
    AgentStream --> RealTimeDelivery[Real-time Response Delivery]:::delivery
    RealTimeDelivery --> ResponseComplete[Response Complete]:::complete

    classDef input fill:#e1f5fe,stroke:#01579b,stroke-width:3px,color:#000
    classDef parse fill:#f3e5f5,stroke:#4a148c,stroke-width:2px,color:#000
    classDef extract fill:#e8f5e8,stroke:#1b5e20,stroke-width:2px,color:#000
    classDef validate fill:#fff3e0,stroke:#e65100,stroke-width:2px,color:#000
    classDef valid fill:#e0f2f1,stroke:#00695c,stroke-width:2px,color:#000
    classDef invalid fill:#ffcdd2,stroke:#d32f2f,stroke-width:2px,color:#000
    classDef rag fill:#f9fbe7,stroke:#827717,stroke-width:2px,color:#000
    classDef error fill:#ffebee,stroke:#c62828,stroke-width:2px,color:#000
    classDef context fill:#e8eaf6,stroke:#283593,stroke-width:2px,color:#000
    classDef stream fill:#e4f7ff,stroke:#0277bd,stroke-width:2px,color:#000
    classDef delivery fill:#f1f8e9,stroke:#33691e,stroke-width:2px,color:#000
    classDef complete fill:#fff8e1,stroke:#f57f17,stroke-width:2px,color:#000
```

### 1. Task Request Parsing
```go
fmt.Printf("🟢 Processing streaming task request: %s\n", taskRequest.ID)
userMessage := taskRequest.Params.Message.Parts[0].Text
fmt.Printf("🔵 UserMessage: %s\n", userMessage)
fmt.Printf("🟡 TaskRequest Metadata: %v\n", taskRequest.Params.MetaData)
```

**Parsing Process**:
- **Request ID**: Unique identifier for tracking
- **Message Extraction**: Retrieves user's text from message parts
- **Metadata Access**: Extracts routing and context information
- **Logging**: Comprehensive debug output for monitoring

### 2. Skill-based Routing
```go
var userPrompt string

switch taskRequest.Params.MetaData["skill"] {
case "ask_for_something":
    userPrompt = userMessage
default:
    userPrompt = "Be nice, and explain that " + fmt.Sprintf("%v", taskRequest.Params.MetaData["skill"]) + " is not a valid task ID."
}
```

**Routing Logic**:
- **Valid Skill**: `ask_for_something` processes user message directly
- **Invalid Skill**: Generates helpful error message
- **Extensible**: Easy to add new skills by extending the switch statement
- **User-Friendly**: Polite error handling for invalid requests

### 3. RAG Integration
```go
bossAgentMessages, err := GeneratePromptMessagesWithSimilarities(ctx, &client, bossAgent.GetName(), userPrompt, similaritySearchLimit, similaritySearchMaxResults)

if err != nil {
    ui.Println(ui.Red, "Error:", err)
}
```

**RAG Enhancement**:
- **Context Search**: Finds relevant background information
- **Similarity Threshold**: Configurable relevance filtering
- **Result Limiting**: Controls context volume
- **Error Resilience**: Continues without RAG if search fails

### 4. Boss Agent Streaming Execution
```go
_, err = bossAgent.RunStream(
    bossAgentMessages,
    func(content string) error {
        if content != "" {
            fmt.Print(content)         // Print to console for debugging
            return streamFunc(content) // Stream to client
        }
        return nil // Continue streaming
    })
```

**Streaming Execution**:
- **Real-time Processing**: Content delivered as it's generated
- **Dual Output**: Console logging and client streaming
- **Error Handling**: Stream errors propagated to client
- **Content Filtering**: Empty content ignored

## Server Initialization

### A2A Server Creation
```go
a2aServer := a2a.NewA2AServerWithStreaming(helpers.StringToInt(httpPort), agentCard, agentStreamCallback)
fmt.Println("🚀 Starting A2A server with streaming support on port:", httpPort)
if err := a2aServer.Start(); err != nil {
    fmt.Printf("❌ Failed to start A2A server: %v\n", err)
}
```

**Initialization Parameters**:
- **Port Configuration**: HTTP port from environment variable
- **Agent Card**: Service metadata and capabilities
- **Streaming Callback**: Request handling function
- **Error Handling**: Startup failure reporting

### Server Capabilities

```mermaid
flowchart LR
    ServerCapabilities[A2A Server Capabilities]:::main
    ServerCapabilities --> HTTPServer[HTTP Server]:::http
    ServerCapabilities --> StreamingSupport[Streaming Support]:::stream
    ServerCapabilities --> JSONRPCCompliance[JSON-RPC 2.0 Compliance]:::jsonrpc
    ServerCapabilities --> ServiceDiscovery[Service Discovery]:::discovery

    HTTPServer --> PortBinding["Port Binding<br/>Configurable HTTP port"]:::binding
    HTTPServer --> RequestRouting["Request Routing<br/>Endpoint mapping"]:::routing

    StreamingSupport --> RealTimeResponse["Real-time Response<br/>Chunk-by-chunk delivery"]:::realtime
    StreamingSupport --> CallbackIntegration["Callback Integration<br/>Custom response handlers"]:::callback

    JSONRPCCompliance --> ProtocolStandard["Protocol Standard<br/>JSON-RPC 2.0 spec"]:::protocol
    JSONRPCCompliance --> MessageFormat["Message Format<br/>Structured requests/responses"]:::format

    ServiceDiscovery --> AgentCardExposure["Agent Card Exposure<br/>Capability advertisement"]:::exposure
    ServiceDiscovery --> SkillRegistration["Skill Registration<br/>Available operations"]:::skills

    classDef main fill:#e1f5fe,stroke:#01579b,stroke-width:3px,color:#000
    classDef http fill:#f3e5f5,stroke:#4a148c,stroke-width:2px,color:#000
    classDef stream fill:#e8f5e8,stroke:#1b5e20,stroke-width:2px,color:#000
    classDef jsonrpc fill:#fff3e0,stroke:#e65100,stroke-width:2px,color:#000
    classDef discovery fill:#fce4ec,stroke:#880e4f,stroke-width:2px,color:#000
    classDef binding fill:#e0f2f1,stroke:#00695c,stroke-width:2px,color:#000
    classDef routing fill:#f9fbe7,stroke:#827717,stroke-width:2px,color:#000
    classDef realtime fill:#e8eaf6,stroke:#283593,stroke-width:2px,color:#000
    classDef callback fill:#e4f7ff,stroke:#0277bd,stroke-width:2px,color:#000
    classDef protocol fill:#f1f8e9,stroke:#33691e,stroke-width:2px,color:#000
    classDef format fill:#fff8e1,stroke:#f57f17,stroke-width:2px,color:#000
    classDef exposure fill:#c8e6c9,stroke:#388e3c,stroke-width:2px,color:#000
    classDef skills fill:#ffcdd2,stroke:#d32f2f,stroke-width:2px,color:#000
```

## Communication Protocol

### JSON-RPC 2.0 Message Format
The server expects requests in JSON-RPC 2.0 format:

```json
{
  "id": "task-1699123456",
  "jsonrpc": "2.0",
  "method": "message/send",
  "params": {
    "message": {
      "role": "user",
      "parts": [
        {
          "text": "Hello, boss! What challenges await?",
          "type": "text"
        }
      ]
    },
    "metadata": {
      "skill": "ask_for_something"
    }
  }
}
```

### Streaming Response Protocol
Responses are delivered in real-time chunks:

```json
{
  "id": "task-1699123456",
  "result": {
    "status": {
      "state": "completed"
    },
    "history": [
      {
        "role": "assistant",
        "parts": [
          {
            "text": "Greetings, adventurer! I am the guardian of this realm...",
            "type": "text"
          }
        ]
      }
    ]
  }
}
```

## Debug Mode Integration

### Debug Command Handling
```go
if strings.HasPrefix(userPrompt, "/debug") {
    msg.DisplayHistory(bossAgent)
}
```

**Debug Features**:
- **Trigger**: Commands starting with `/debug`
- **History Display**: Shows conversation history
- **Development Aid**: Helps troubleshoot agent behavior
- **Non-intrusive**: Doesn't affect normal operation

## Integration Benefits

### Service Architecture
- **Microservice Design**: Independent boss service
- **Standard Protocol**: A2A compatibility with other services
- **Configurable Deployment**: Environment-based configuration
- **Monitoring Support**: Comprehensive logging and debug output

### Extensibility
- **Skill Addition**: Easy to add new capabilities
- **Custom Handlers**: Flexible callback system
- **Protocol Compliance**: Standard JSON-RPC 2.0 interface
- **Client Agnostic**: Works with any A2A-compatible client

---

⬅️ **Back to:** [Boss Agent Schema](100-boss-agent-schema.md)