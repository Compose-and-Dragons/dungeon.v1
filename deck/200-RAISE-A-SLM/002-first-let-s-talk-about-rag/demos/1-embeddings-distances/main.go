package main

import (
	"context"
	"embeddings-demo/rag"
	"fmt"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

var chunks = []string{
	`Orcs run in the forest`,
	`Dragons fly in the sky`,
	`Goblins swim in the pond`,
	`Krakens swim in the sea`,
}

func main() {
	ctx := context.Background()

	baseURL := "http://localhost:12434/engines/llama.cpp/v1/"
	embeddingsModel := "ai/mxbai-embed-large"

	fmt.Println(baseURL)
	fmt.Println(embeddingsModel)

	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey(""),
	)

	// -------------------------------------------------
	// Generate embeddings from user question
	// -------------------------------------------------
	// USER MESSAGE:
	userQuestion := "Which creatures swim?"

	fmt.Println("⏳ Creating embeddings from user question...:", userQuestion)

	embeddingsFromUserQuestion, err := client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: openai.String(userQuestion),
		},
		Model: embeddingsModel,
	})
	if err != nil {
		fmt.Println(err)
	}

	// -------------------------------------------------
	// Generate embeddings from chunks
	// -------------------------------------------------
	fmt.Println("⏳ Creating embeddings from chunks...")

	for _, chunk := range chunks {
		embeddingsResponse, err := client.Embeddings.New(ctx, openai.EmbeddingNewParams{
			Input: openai.EmbeddingNewParamsInputUnion{
				OfString: openai.String(chunk),
			},
			Model: embeddingsModel,
		})

		if err != nil {
			fmt.Println(err)
		} else {
			cosineSimilarity := rag.CosineSimilarity(
				embeddingsResponse.Data[0].Embedding,
				embeddingsFromUserQuestion.Data[0].Embedding,
			)
			fmt.Println("🔗 Cosine similarity with ", chunk, "=", cosineSimilarity, IsGoodCosineSimilarity(cosineSimilarity))
		}
	}
}

func IsGoodCosineSimilarity(cosineSimilarity float64) string {
	if cosineSimilarity > 0.65 {
		return "✅"
	} else {
		return "❌"
	}
}
