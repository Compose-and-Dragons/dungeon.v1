package main

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"fmt"
	"log"

	//mcp_golang "github.com/metoro-io/mcp-golang"
	//mcp_http "github.com/metoro-io/mcp-golang/transport/http"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	"github.com/openai/openai-go/v2/shared"
)

func main() {
	ctx := context.Background()

	// Docker Model Runner base URL
	chatURL := os.Getenv("MODEL_RUNNER_BASE_URL")
	model := os.Getenv("MODEL_RUNNER_LLM_TOOLS")

	openAIClient := openai.NewClient(
		option.WithBaseURL(chatURL),
		option.WithAPIKey(""),
	)

	// BEGIN: MCP client initialization
	fmt.Println("🚀 Initializing MCP StreamableHTTP client...")
	// Create HTTP transport
	httpURL := os.Getenv("MCP_HTTP_SERVER_URL")
	httpTransport, err := transport.NewStreamableHTTP(httpURL)
	if err != nil {
		log.Fatalf("😡 Failed to create HTTP transport: %v", err)
	}
	// Create client with the transport
	mcpClient := client.NewClient(httpTransport)
	// Start the client
	if err := mcpClient.Start(ctx); err != nil {
		log.Fatalf("😡 Failed to start client: %v", err)
	}

	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "MCP-Go Simple Client Example",
		Version: "1.0.0",
	}
	initRequest.Params.Capabilities = mcp.ClientCapabilities{}

	//serverInfo, err := mcpClient.Initialize(ctx, initRequest)

	_, err = mcpClient.Initialize(ctx, initRequest)
	if err != nil {
		log.Fatalf("😡 Failed to initialize: %v", err)
	}
	// END:  of MCP client initialization

	// TOOLS:
	toolsRequest := mcp.ListToolsRequest{}
	// STEP 1: make the MCP Request
	toolsResult, err := mcpClient.ListTools(ctx, toolsRequest)
	if err != nil {
		log.Fatalf("😡 Failed to list tools: %v", err)
	}
	// STEP 2: Convert the MCP tools to to OpenAI tools
	openAITools := ConvertMCPToolsToOpenAITools(toolsResult)

	// USER MESSAGE:
	userQuestion := openai.UserMessage(`
		I want to speak to an elf.
		I want to talk about magic.
		I want to speak with a Dwarf
	`)

	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			userQuestion,
		},
		Tools:       openAITools,
		ParallelToolCalls: openai.Bool(true), // Sequential tool calls
		Model:       model,
		Temperature: openai.Opt(0.0),
	}

	// Make [COMPLETION] request
	completion, err := openAIClient.Chat.Completions.New(ctx, params)
	if err != nil {
		panic(err)
	}
	// TOOL CALLS: Extract tool calls from the response
	toolCalls := completion.Choices[0].Message.ToolCalls

	// Return early if there are no tool calls
	if len(toolCalls) == 0 {
		fmt.Println("😡 No function call")
		fmt.Println()
		return
	}

	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("🛠️  Tool calls detected by the LLM:")
	fmt.Println(strings.Repeat("=", 50))
	// Execute TOOL CALLS: Display the tool calls
	for _, toolCall := range toolCalls {
		fmt.Println("✅", toolCall.Function.Name, toolCall.Function.Arguments)

		// Parse the tool arguments from JSON string
		var args map[string]any
		args, _ = JsonStringToMap(toolCall.Function.Arguments)

		// NOTE: Call the MCP tool with the arguments
		request := mcp.CallToolRequest{}
		request.Params.Name = toolCall.Function.Name
		request.Params.Arguments = args

		toolResponse, err := mcpClient.CallTool(ctx, request)
		if err != nil {
			log.Fatalf("😡 Failed to call tool %s: %v", toolCall.Function.Name, err)
		}
		if toolResponse == nil || len(toolResponse.Content) == 0 {
			log.Fatalf("😡 No response from tool %s", toolCall.Function.Name)
		}

		fmt.Println("🛠️  Tool response:", toolResponse.Content[0].(mcp.TextContent).Text)

	}

}

// ConvertMCPToolsToOpenAITools transforms MCP tool definitions into OpenAI tool format
func ConvertMCPToolsToOpenAITools(tools *mcp.ListToolsResult) []openai.ChatCompletionToolUnionParam {
	openAITools := make([]openai.ChatCompletionToolUnionParam, len(tools.Tools))
	for i, tool := range tools.Tools {

		openAITools[i] = openai.ChatCompletionFunctionTool(shared.FunctionDefinitionParam{
			Name:        tool.Name,
			Description: openai.String(tool.Description),
			Parameters: shared.FunctionParameters{
				"type":       "object",
				"properties": tool.InputSchema.Properties,
				"required":   tool.InputSchema.Required,
			},
		},
		)
	}
	return openAITools
}

// JsonStringToMap parses a JSON string and converts it to a map with string keys and any values
func JsonStringToMap(jsonString string) (map[string]any, error) {
	var result map[string]any
	err := json.Unmarshal([]byte(jsonString), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
