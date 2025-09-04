package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

func main() {
	// Docker Model Runner Chat base URL
	baseURL := "http://localhost:12434/engines/llama.cpp/v1/"
	model := "ai/qwen2.5:0.5B-F16"

	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey(""),
	)

	ctx := context.Background()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("🤖 [%s](%s) ask me something - /bye to exit> ", "Elara", model)
		userMessage, _ := reader.ReadString('\n')

		if strings.HasPrefix(userMessage, "/bye") {
			fmt.Println("👋 Bye!")
			break
		}

		messages := []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(`
			You are an expert of medieval role playing games
			Your name is Elara, Weaver of the Arcane
		`),
			openai.UserMessage(userMessage),
		}

		param := openai.ChatCompletionNewParams{
			Messages:    messages,
			Model:       model,
			Temperature: openai.Opt(1.8),
			TopP:        openai.Opt(0.9),
		}

		stream := client.Chat.Completions.NewStreaming(ctx, param)

		fmt.Println()

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
		fmt.Println("\n\n", strings.Repeat("-", 80))

	}

}
