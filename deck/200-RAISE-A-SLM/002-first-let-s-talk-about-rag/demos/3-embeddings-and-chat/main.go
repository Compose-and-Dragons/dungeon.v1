package main

import (
	"context"
	"embeddings-chat/rag"
	"fmt"
	"log"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

var chunks = []string{
	`# Orcs
	Orcs are savage, brutish humanoids with dark green skin and prominent tusks. 
	These fierce warriors inhabit dense forests where they hunt in packs, 
	using crude but effective weapons forged from scavenged metal and bone. 
	Their tribal society revolves around strength and combat prowess, 
	making them formidable opponents for any adventurer brave enough to enter their woodland domain.`,

	`# Dragons
	Dragons are magnificent and ancient creatures of immense power, soaring through the skies on massive wings. 
	These intelligent beings possess scales that shimmer like precious metals and breathe devastating elemental attacks. 
	Known for their vast hoards of treasure and centuries of accumulated knowledge, 
	dragons command both fear and respect throughout the realm. 
	Their aerial dominance makes them nearly untouchable in their celestial domain.`,

	`# Goblins
	Goblins are small, cunning creatures with mottled green skin and sharp, pointed ears. 
	Despite their diminutive size, they are surprisingly agile swimmers who have adapted to life around ponds and marshlands. 
	These mischievous beings are known for their quick wit and tendency to play pranks on unwary travelers. 
	They build elaborate underwater lairs connected by hidden tunnels beneath the murky pond waters.`,

	`# Krakens
	Krakens are colossal sea monsters with massive tentacles that can crush entire ships with ease. 
	These legendary creatures dwell in the deepest ocean trenches, surfacing only to hunt or when disturbed. 
	Their intelligence rivals that of the wisest sages, and their tentacles can stretch for hundreds of feet. 
	Sailors speak in hushed tones of these maritime titans, whose very presence can create devastating whirlpools 
	and tidal waves that reshape entire coastlines.`,
}

func main() {
	ctx := context.Background()

	baseURL := "http://localhost:12434/engines/llama.cpp/v1/"
	embeddingsModel := "ai/mxbai-embed-large"
	chatModel := "ai/qwen2.5:0.5B-F16"

	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey(""),
	)

	// -------------------------------------------------
	// Create a vector store
	// -------------------------------------------------
	store := rag.MemoryVectorStore{
		Records: make(map[string]rag.VectorRecord),
	}

	// -------------------------------------------------
	// STEP 1: Create and save the embeddings from the chunks
	// -------------------------------------------------
	fmt.Println("⏳ Creating the embeddings...")

	for _, chunk := range chunks {
		// EMBEDDING COMPLETION:
		embeddingsResponse, err := client.Embeddings.New(ctx, openai.EmbeddingNewParams{
			Input: openai.EmbeddingNewParamsInputUnion{
				OfString: openai.String(chunk),
			},
			Model: embeddingsModel,
		})

		if err != nil {
			fmt.Println(err)
		} else {
			_, errSave := store.Save(rag.VectorRecord{
				Prompt:    chunk,
				Embedding: embeddingsResponse.Data[0].Embedding,
			})
			if errSave != nil {
				fmt.Println("😡:", errSave)
			}
		}
	}

	fmt.Println("✋", "Embeddings created, total of records", len(store.Records))
	fmt.Println()

	// -------------------------------------------------
	// Search for similarities
	// -------------------------------------------------

	// USER MESSAGE:
	userQuestion := "Tell me something about the dragons"

	fmt.Println("⏳ Searching for similarities...")

	// -------------------------------------------------
	// STEP 2: EMBEDDING COMPLETION:
	// Create embedding from the user question
	// -------------------------------------------------
	embeddingsResponse, err := client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: openai.String(userQuestion),
		},
		Model: embeddingsModel,
	})
	if err != nil {
		log.Fatal("😡:", err)
	}

	// -------------------------------------------------
	// STEP 3: SIMILARITY SEARCH: use the vector store to find similar chunks
	// -------------------------------------------------
	// Create a vector record from the user embedding
	embeddingFromUserQuestion := rag.VectorRecord{
		Embedding: embeddingsResponse.Data[0].Embedding,
	}

	similarities, _ := store.SearchTopNSimilarities(embeddingFromUserQuestion, 0.6, 2)

	documentsContent := "Documents:\n"

	for _, similarity := range similarities {
		fmt.Println("✅ CosineSimilarity:", similarity.CosineSimilarity, "Chunk:", similarity.Prompt)
		documentsContent += similarity.Prompt
	}
	documentsContent += "\n"
	fmt.Println("\n✋", "Similarities found, total of records", len(similarities))
	fmt.Println()

	// -------------------------------------------------
	// STEP 4: Generate CHAT COMPLETION:
	// -------------------------------------------------
	messages := []openai.ChatCompletionMessageParamUnion{
		// SYSTEM MESSAGE:
		openai.SystemMessage(`
			You are a useful AI agent and RPG expert. 
			Use only the following documents, 
			like a dungeon master, 
			create an history to answer the user question:
		`),
		// SIMILARITIES:
		openai.SystemMessage(documentsContent),
		// USER MESSAGE:
		openai.UserMessage(userQuestion),
	}

	param := openai.ChatCompletionNewParams{
		Messages:    messages,
		Model:       chatModel,
		Temperature: openai.Opt(0.8),
	}

	stream := client.Chat.Completions.NewStreaming(ctx, param)

	for stream.Next() {
		chunk := stream.Current()
		// Stream each chunk as it arrives
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			fmt.Print(chunk.Choices[0].Delta.Content)
		}
	}

	if err := stream.Err(); err != nil {
		log.Fatalln("😡:", err)
	}

	fmt.Println()
	store.SaveJSONToFile("vectorstore.json")
	fmt.Println("\n✋", "Vector store saved to vectorstore.json")
}
