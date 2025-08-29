package agents

import (
	"context"
	"fmt"
	"strings"

	"github.com/micro-agent/micro-agent-go/agent/helpers"
	"github.com/micro-agent/micro-agent-go/agent/mu"
	"github.com/micro-agent/micro-agent-go/agent/rag"
	"github.com/micro-agent/micro-agent-go/agent/ui"

	"github.com/openai/openai-go/v2"
)

var agentsStores = make(map[string]rag.MemoryVectorStore)

// GenerateEmbeddings reads a context file, splits it into chunks, generates embeddings,
// and stores them in the vector store for the specified agent
func GenerateEmbeddings(ctx context.Context, client *openai.Client, name string, contextInstructionsContentPath string) error {

	ui.Println(ui.Green, strings.Repeat("─", 80))
	ui.Println(ui.Green,"🚧 Generating embeddings for agent:", name)
	ui.Println(ui.Green, strings.Repeat("─", 80))


	embeddingAgent, err := mu.NewAgent(ctx, "vector-agent",
		mu.WithClient(*client),
		mu.WithEmbeddingParams(
			openai.EmbeddingNewParams{
				Model: helpers.GetEnvOrDefault("EMBEDDING_MODEL", "ai/mxbai-embed-large:latest"),
			},
		),
	)
	if err != nil {
		fmt.Println("🔶 Error creating embedding agent, creating ghost agent instead:", err)
		return err
	}

	fmt.Println("✅ Embedding agent created successfully")

	if contextInstructionsContentPath == "" {
		fmt.Println("🔶 No context path provided, using default instructions.")
		return fmt.Errorf("no context path provided")
	}

	// Read the content of the file at contextInstructionsContentPath
	contextInstructionsContent, err := helpers.ReadTextFile(contextInstructionsContentPath)
	if err != nil {
		fmt.Println("🔶 Error reading the file, using default instructions:", err)
		return err
	}

	// [RAG] Initialize the vector store for the agent
	agentsStores[name] = rag.MemoryVectorStore{
		Records: make(map[string]rag.VectorRecord),
	}
	store := agentsStores[name]

	chunks := rag.SplitMarkdownBySections(contextInstructionsContent)

	for idx, chunk := range chunks {
		fmt.Println("🔶 Chunk", idx, ":", chunk)
		embeddingVector, err := embeddingAgent.GenerateEmbeddingVector(chunk)
		if err != nil {
			return err
		}
		_, errSave := store.Save(rag.VectorRecord{
			Prompt:    chunk,
			Embedding: embeddingVector,
		})

		if errSave != nil {
			fmt.Println("🔴 When saving the vector", errSave)
			return errSave
		}
		fmt.Println("✅ Chunk", idx, "saved with embedding:", len(embeddingVector))
	}
	fmt.Println("📝 Total records in the vector store:", len(store.Records))
	ui.Println(ui.Green, strings.Repeat("─", 80))
	fmt.Println()

	return nil
}

