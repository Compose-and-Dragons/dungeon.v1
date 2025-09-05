package main

import (
	"context"
	"os"

	"fmt"
	"log"

	//mcp_golang "github.com/metoro-io/mcp-golang"
	//mcp_http "github.com/metoro-io/mcp-golang/transport/http"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/shared"
)

func main() {
	ctx := context.Background()

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

	fmt.Println("🛠️  Available tools:")

	for _, tool := range openAITools {
		fmt.Println(" -", tool.GetFunction().Name, tool.GetFunction().Description, tool.GetFunction().Parameters)
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
