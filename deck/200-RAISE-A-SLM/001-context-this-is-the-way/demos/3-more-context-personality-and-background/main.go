package main

import (
	"bufio"
	"context"
	"encoding/json"
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

	// ✋ NOTE: load the system instructions from a file
	systemInstructions, err := os.ReadFile("./sorcerer_system_instructions.md")
	if err != nil {
		log.Fatal("😡:", err)
	}
	backgroundAndPersonality, err := os.ReadFile("./sorcerer_background_and_personality.md")
	if err != nil {
		log.Fatal("😡:", err)
	}

	// NOTE: initialize the messages slice with a system message to set the behavior of the assistant
	// MEMORY:
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(string(systemInstructions)),
		openai.SystemMessage(string(backgroundAndPersonality)),
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("🤖🧠 [%s](%s) ask me something - /bye to exit> ", "Elara", model)
		userMessage, _ := reader.ReadString('\n')

		if strings.HasPrefix(userMessage, "/bye") {
			fmt.Println("👋 Bye!")
			break
		}

		if strings.HasPrefix(userMessage, "/memory") {
			DisplayConversationalMemory(messages)
			continue
		}

		// NOTE: append the USER MESSAGE: to the messages slice
		messages = append(messages, openai.UserMessage(userMessage))

		// Keep the context of the conversation by appending the user message
		// and the assistant response to the messages slice
		// In a real application, you might want to limit the size of this context
		// to avoid exceeding token limits (or model limits + machine limits ex: RPI).

		param := openai.ChatCompletionNewParams{
			Messages:    messages,
			Model:       model,
			Temperature: openai.Opt(1.8),
			TopP:        openai.Opt(0.9),
		}

		stream := client.Chat.Completions.NewStreaming(ctx, param)

		fmt.Println()

		answer := ""
		for stream.Next() {
			chunk := stream.Current()
			// Stream each chunk as it arrives
			if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
				content := chunk.Choices[0].Delta.Content
				// NOTE: accumulate the content of the assistant's response
				answer += content
				fmt.Print(content)
			}
		}

		if err := stream.Err(); err != nil {
			log.Fatalln("😡:", err)
		}

		// NOTE: Append the ASSISTANT MESSAGE: (response) to the messages slice
		messages = append(messages, openai.AssistantMessage(answer))

		fmt.Println("\n\n", strings.Repeat("-", 80))

	}

}

// MessageToMap converts an OpenAI chat message to a map with string keys and values
func MessageToMap(message openai.ChatCompletionMessageParamUnion) (map[string]string, error) {
	jsonData, err := message.MarshalJSON()
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil, err
	}

	stringMap := make(map[string]string)
	for key, value := range result {
		if str, ok := value.(string); ok {
			stringMap[key] = str
		}
	}

	return stringMap, nil
}

func DisplayConversationalMemory(messages []openai.ChatCompletionMessageParamUnion) {
	// remove the /debug part from the input
	fmt.Println()
	fmt.Println("📝 Messages history / Conversational memory:")
	for i, message := range messages {
		printableMessage, err := MessageToMap(message)
		if err != nil {
			fmt.Printf("😡 Error converting message to map: %v\n", err)
			continue
		}
		fmt.Print("-", i, " ")
		fmt.Print(printableMessage["role"], ": ")
		fmt.Println(printableMessage["content"])
	}
	fmt.Println("📝 End of the messages")
	fmt.Println()
}
