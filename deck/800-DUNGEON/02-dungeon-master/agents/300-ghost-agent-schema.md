# Ghost Agent Schema

⬅️ **Back to:** [NPC Agents System](../003-schema-npc-agents-system.md)

## Overview

The Ghost Agent is a test/simulation agent that mimics AI behavior without actual AI processing. It's designed for development and testing purposes, providing predictable responses based on keyword matching.

```mermaid
flowchart TD
    GhostAgent[Ghost Agent]:::main
    GhostAgent --> Structure[Agent Structure]:::struct
    GhostAgent --> Simulation[Response Simulation]:::simulation
    GhostAgent --> Interface[Interface Implementation]:::interface

    Structure --> Fields["name: string<br/>messages: []ChatCompletionMessage<br/>responseFormat: ResponseFormat"]:::field

    Simulation --> RunStream("RunStream Method<br/><a href='/dungeon-master/agents/ghost.agent.go#L144'>ghost.agent.go:144</a>"):::stream
    Simulation --> KeywordMatching("Keyword Matching<br/><a href='/dungeon-master/agents/ghost.agent.go#L212'>ghost.agent.go:212</a>"):::keywords
    Simulation --> StreamingSimulation("Simulate Streaming<br/><a href='/dungeon-master/agents/ghost.agent.go#L161'>ghost.agent.go:161</a>"):::streaming

    Interface --> Implemented["Basic Methods<br/>GetName, SetName<br/>GetMessages, SetMessages<br/>RunStream"]:::impl
    Interface --> Unimplemented["Most mu.Agent methods<br/>panic('unimplemented')"]:::unimpl

    classDef main fill:#e1f5fe,stroke:#01579b,stroke-width:3px,color:#000
    classDef struct fill:#f3e5f5,stroke:#4a148c,stroke-width:2px,color:#000
    classDef simulation fill:#fff3e0,stroke:#e65100,stroke-width:2px,color:#000
    classDef interface fill:#fce4ec,stroke:#880e4f,stroke-width:2px,color:#000
    classDef field fill:#f9fbe7,stroke:#827717,stroke-width:2px,color:#000
    classDef stream fill:#ffebee,stroke:#c62828,stroke-width:2px,color:#000
    classDef keywords fill:#f3e5f5,stroke:#4a148c,stroke-width:2px,color:#000
    classDef streaming fill:#fff8e1,stroke:#f57f17,stroke-width:2px,color:#000
    classDef impl fill:#c8e6c9,stroke:#388e3c,stroke-width:2px,color:#000
    classDef unimpl fill:#ffcdd2,stroke:#d32f2f,stroke-width:2px,color:#000
```

## Key Features

### Structure
Simple struct with three fields for basic agent functionality:
- **name**: Agent identifier
- **messages**: Message history storage
- **responseFormat**: Response format configuration

### Response Simulation
- **Keyword Matching**: Predefined responses for common keywords (hello, weather, code, etc.)
- **Streaming Simulation**: Word-by-word delivery with 50ms delays
- **Default Fallback**: Always provides a response when no keywords match

### Interface Implementation
- **Implemented**: Basic methods like GetName, SetName, GetMessages, RunStream
- **Unimplemented**: Most mu.Agent methods panic with "unimplemented"

## Usage

**Development Tool**: Safe testing without AI costs or dependencies
**Predictable Behavior**: Consistent responses for automated testing
**Fallback Agent**: Used when real agents fail to create ([guard.agent.go:73](guard.agent.go:73))

## Code References

- **Constructor**: [ghost.agent.go:19](ghost.agent.go:19)
- **Main Method**: [ghost.agent.go:144](ghost.agent.go:144)
- **Keyword Responses**: [ghost.agent.go:212](ghost.agent.go:212)

---

⬅️ **Back to:** [NPC Agents System](../003-schema-npc-agents-system.md)