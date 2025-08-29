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

var AgentsStores = make(map[string]rag.MemoryVectorStore)

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
		fmt.Println("🔶 Error creating embedding agent", err)
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
	AgentsStores[name] = rag.MemoryVectorStore{
		Records: make(map[string]rag.VectorRecord),
	}
	store := AgentsStores[name]

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

// SearchSimilarities searches for similar content in the agent's vector store
// based on the input question and returns the top N similar records
func SearchSimilarities(ctx context.Context, client *openai.Client, agentName string, input string, threshold float64, topN int) ([]rag.VectorRecord, error) {
	store := AgentsStores[agentName]

	embeddingAgent, err := mu.NewAgent(ctx, "vector-agent",
		mu.WithClient(*client),
		mu.WithEmbeddingParams(
			openai.EmbeddingNewParams{
				Model: helpers.GetEnvOrDefault("EMBEDDING_MODEL", "ai/mxbai-embed-large:latest"),
			},
		),
	)
	if err != nil {
		fmt.Println("🔶 Error creating embedding agent", err)
		return nil, err
	}

	fmt.Println(strings.Repeat("-", 80))
	questionEmbeddingVector, err := embeddingAgent.GenerateEmbeddingVector(input)
	if err != nil {
		return nil, err
	}

	questionRecord := rag.VectorRecord{Embedding: questionEmbeddingVector}

	similarities, err := store.SearchTopNSimilarities(questionRecord, threshold, topN)
	if err != nil {
		return nil, err
	}

	fmt.Println("📝 Similarities found:", len(similarities))

	for _, similarity := range similarities {
		fmt.Println("✅ CosineSimilarity:", similarity.CosineSimilarity, "Chunk:", similarity.Prompt)
	}

	fmt.Println(strings.Repeat("-", 80))

	return similarities, nil
}

