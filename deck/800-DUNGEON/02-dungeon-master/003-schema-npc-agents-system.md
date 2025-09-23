# NPC Agents System

⬅️ **Back to:** [Dungeon Master Schema](./002-schema-dungeon-master.md)

```mermaid
flowchart TD
    NPCAgents[Create NPC Agents]:::process

    NPCAgents --> GhostAgent("Ghost Agent Test<br/><a href='/dungeon-master/main.go#L111'>main.go:111</a>"):::npc
    NPCAgents --> GuardAgent("Guard Agent + RAG<br/><a href='/dungeon-master/main.go#L117'>main.go:117</a>"):::npc
    NPCAgents --> SorcererAgent("Sorcerer Agent + RAG<br/><a href='/dungeon-master/main.go#L122'>main.go:122</a>"):::npc
    NPCAgents --> MerchantAgent("Merchant Agent + RAG<br/><a href='/dungeon-master/main.go#L127'>main.go:127</a>"):::npc
    NPCAgents --> HealerAgent("Healer Agent + RAG<br/><a href='/dungeon-master/main.go#L132'>main.go:132</a>"):::npc
    NPCAgents --> BossAgent("Boss Agent Remote<br/><a href='/dungeon-master/main.go#L137'>main.go:137</a>"):::npc

    NPCAgents --> AgentsTeam("Create Agents Team Map<br/><a href='/dungeon-master/main.go#L153'>main.go:153</a>"):::team

    GhostAgent --> AgentExecution[Agent Execution Flows]:::execution
    GuardAgent --> AgentExecution
    SorcererAgent --> AgentExecution
    MerchantAgent --> AgentExecution
    HealerAgent --> AgentExecution
    BossAgent --> AgentExecution

    AgentExecution --> GhostExecution("Ghost Agent Execution<br/><a href='/dungeon-master/main.go#L269'>main.go:269</a>"):::npc_exec
    AgentExecution --> RAGExecution[Guard/Sorcerer/Merchant/Healer with RAG]:::npc_exec
    AgentExecution --> BossExecution("Boss Agent Execution<br/><a href='/dungeon-master/main.go#L423'>main.go:423</a>"):::npc_exec

    RAGExecution --> SimilaritySearch("Generate Prompt with Similarities<br/><a href='/dungeon-master/main.go#L611'>main.go:611</a>"):::rag
    SimilaritySearch --> SearchSimilarities[agents.SearchSimilarities]:::rag

    BossExecution --> GameEndLogic("Game End Detection<br/><a href='/dungeon-master/main.go#L438'>main.go:438</a>"):::game_end
    GameEndLogic --> DefeatCondition("You are trapped - Game Over<br/><a href='/dungeon-master/main.go#L438'>main.go:438</a>"):::game_end
    GameEndLogic --> VictoryCondition("You are free - Victory<br/><a href='/dungeon-master/main.go#L454'>main.go:454</a>"):::game_end

    classDef process fill:#f3e5f5,stroke:#4a148c,stroke-width:2px,color:#000
    classDef npc fill:#fce4ec,stroke:#880e4f,stroke-width:2px,color:#000
    classDef team fill:#e0f2f1,stroke:#00695c,stroke-width:2px,color:#000
    classDef execution fill:#f3e5f5,stroke:#4a148c,stroke-width:2px,color:#000
    classDef npc_exec fill:#fce4ec,stroke:#880e4f,stroke-width:2px,color:#000
    classDef rag fill:#e0f2f1,stroke:#00695c,stroke-width:2px,color:#000
    classDef game_end fill:#ffebee,stroke:#c62828,stroke-width:2px,color:#000
```

## NPC Agent Types

### Ghost Agent (<a href="/dungeon-master/main.go#L111-L112">lines 111-112</a>)
- **Purpose**: Test agent for development
- **Implementation**: Simple fake agent
- **Usage**: Always available for testing
- **Execution**: Direct streaming without tools (<a href="/dungeon-master/main.go#L269-L286">lines 269-286</a>)

### Guard Agent (<a href="/dungeon-master/main.go#L117">line 117</a>)
- **Type**: RAG-enabled NPC
- **Implementation**: `agents.GetGuardAgent(ctx, client)`
- **Features**: Similarity search integration
- **Execution**: With RAG context (<a href="/dungeon-master/main.go#L291-L320">lines 291-320</a>)

### Sorcerer Agent (<a href="/dungeon-master/main.go#L122">line 122</a>)
- **Type**: RAG-enabled NPC
- **Implementation**: `agents.GetSorcererAgent(ctx, client)`
- **Features**: Similarity search integration
- **Execution**: With RAG context (<a href="/dungeon-master/main.go#L324-L353">lines 324-353</a>)

### Merchant Agent (<a href="/dungeon-master/main.go#L127">line 127</a>)
- **Type**: RAG-enabled NPC
- **Implementation**: `agents.GetMerchantAgent(ctx, client)`
- **Features**: Similarity search integration
- **Execution**: With RAG context (<a href="/dungeon-master/main.go#L357-L386">lines 357-386</a>)

### Healer Agent (<a href="/dungeon-master/main.go#L132">line 132</a>)
- **Type**: RAG-enabled NPC
- **Implementation**: `agents.GetHealerAgent(ctx, client)`
- **Features**: Similarity search integration
- **Execution**: With RAG context (<a href="/dungeon-master/main.go#L390-L419">lines 390-419</a>)

### Boss Agent (<a href="/dungeon-master/main.go#L137-L141">lines 137-141</a>)
- **Type**: Remote agent
- **Implementation**: `agents.NewBossAgent(name, url)`
- **Features**: Game-ending logic
- **Remote URL**: Configurable via `BOSS_REMOTE_AGENT_URL`
- **Execution**: With victory/defeat detection (<a href="/dungeon-master/main.go#L423-L475">lines 423-475</a>)

## Agent Team Structure (<a href="/dungeon-master/main.go#L153-L161">lines 153-161</a>)

```go
agentsTeam = map[string]mu.Agent{
    idDungeonMasterToolsAgent: dungeonMasterToolsAgent,  // Main DM
    idGhostAgent:              ghostAgent,               // Test agent
    idGuardAgent:              guardAgent,               // RAG-enabled
    idSorcererAgent:           sorcererAgent,            // RAG-enabled
    idMerchantAgent:           merchantAgent,            // RAG-enabled
    idHealerAgent:             healerAgent,              // RAG-enabled
    idBossAgent:               bossAgent,                // Remote agent
}
```

## RAG Integration (<a href="/dungeon-master/main.go#L611-L637">GeneratePromptMessagesWithSimilarities</a>)

The RAG (Retrieval-Augmented Generation) system enhances NPC interactions:

1. **Similarity Search** (<a href="/dungeon-master/main.go#L614">line 614</a>)
   - Uses `agents.SearchSimilarities()` to find relevant context
   - Configured with similarity limit and max results

2. **Context Enhancement** (<a href="/dungeon-master/main.go#L623-L630">lines 623-630</a>)
   - Adds contextual information to system messages
   - Improves NPC response relevance

3. **Fallback Handling** (<a href="/dungeon-master/main.go#L632-L636">lines 632-636</a>
   - Graceful degradation when no similarities found
   - Maintains conversation flow

## Game End Logic (Boss Agent)

### Victory Condition (<a href="/dungeon-master/main.go#L454-L468">lines 454-468</a>)
- **Trigger**: Response contains "you are free"
- **Action**: Display victory message and player info
- **MCP Call**: Direct call to `get_player_info`

### Defeat Condition (<a href="/dungeon-master/main.go#L438-L452">lines 438-452</a>)
- **Trigger**: Response contains "you are trapped"
- **Action**: Display game over message and player info
- **MCP Call**: Direct call to `get_player_info`

---

⬅️ **Back to:** [Dungeon Master Schema](002-schema-dungeon-master.md)