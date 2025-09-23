package agents

/*
This is a stub implementation of a BossAgent that communicates with a remote agent via HTTP.
It implements the mu.Agent interface.
*/

import (
	"fmt"
	"time"

	"github.com/micro-agent/micro-agent-go/agent/experimental/a2a"
	"github.com/micro-agent/micro-agent-go/agent/mu"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/shared"
)

type BossAgent struct {
	name      string
	messages  []openai.ChatCompletionMessageParamUnion
	remoteURL string
	client    *a2a.A2AClient
	connected bool
}

// NewBossAgent creates a new boss agent instance
func NewBossAgent(name string, remoteURL string) mu.Agent {

	// Initialize the A2A client
	client := a2a.NewA2AClient(remoteURL)

	// First, ping the agent to verify connection
	fmt.Println("🔍 Pinging agent...")
	agentCard, err := client.PingAgent()
	if err != nil {
		fmt.Printf("❌ Failed to ping agent: %v\n", err)
	}

	fmt.Printf("✅ Connected to agent: %s\n", agentCard.Name)
	fmt.Printf("📝 Description: %s\n", agentCard.Description)
	fmt.Printf("🔧 Available skills: %v\n", len(agentCard.Skills))
	fmt.Println()

	return &BossAgent{
		name:      name,
		messages:  []openai.ChatCompletionMessageParamUnion{},
		remoteURL: remoteURL,
		client:    client,
	}
}

func (b *BossAgent) IsConnected() bool {
	return b.connected
}

// AddMessage implements mu.Agent.
func (b *BossAgent) AddMessage(message openai.ChatCompletionMessageParamUnion) {
	panic("unimplemented")
}

// AddMessages implements mu.Agent.
func (b *BossAgent) AddMessages(messages []openai.ChatCompletionMessageParamUnion) {
	panic("unimplemented")
}

// DetectToolCalls implements mu.Agent.
func (b *BossAgent) DetectToolCalls(messages []openai.ChatCompletionMessageParamUnion, toolCallBack func(functionName string, arguments string) (string, error)) (string, []string, string, error) {
	panic("unimplemented")
}

// DetectToolCallsStream implements mu.Agent.
func (b *BossAgent) DetectToolCallsStream(messages []openai.ChatCompletionMessageParamUnion, toolCallback func(functionName string, arguments string) (string, error), streamCallback func(content string) error) (string, []string, string, error) {
	panic("unimplemented")
}

// GenerateEmbeddingVector implements mu.Agent.
func (b *BossAgent) GenerateEmbeddingVector(content string) ([]float64, error) {
	panic("unimplemented")
}

// GetDescription implements mu.Agent.
func (b *BossAgent) GetDescription() string {
	panic("unimplemented")
}

// GetFirstNMessages implements mu.Agent.
func (b *BossAgent) GetFirstNMessages(n int) []openai.ChatCompletionMessageParamUnion {
	panic("unimplemented")
}

// GetLastMessage implements mu.Agent.
func (b *BossAgent) GetLastMessage() (openai.ChatCompletionMessageParamUnion, bool) {
	panic("unimplemented")
}

// GetLastNMessages implements mu.Agent.
func (b *BossAgent) GetLastNMessages(n int) []openai.ChatCompletionMessageParamUnion {
	panic("unimplemented")
}

// GetMessages implements mu.Agent.
func (b *BossAgent) GetMessages() []openai.ChatCompletionMessageParamUnion {
	panic("unimplemented")
}

// GetMetaData implements mu.Agent.
func (b *BossAgent) GetMetaData() any {
	panic("unimplemented")
}

// GetModel implements mu.Agent.
func (b *BossAgent) GetModel() shared.ChatModel {
	return "Remote Model" 
}

// GetName implements mu.Agent.
func (b *BossAgent) GetName() string {
	return b.name
}

// GetResponseFormat implements mu.Agent.
func (b *BossAgent) GetResponseFormat() openai.ChatCompletionNewParamsResponseFormatUnion {
	panic("unimplemented")
}

// PrependMessage implements mu.Agent.
func (b *BossAgent) PrependMessage(message openai.ChatCompletionMessageParamUnion) {
	panic("unimplemented")
}

// PrependMessages implements mu.Agent.
func (b *BossAgent) PrependMessages(messages []openai.ChatCompletionMessageParamUnion) {
	panic("unimplemented")
}

// RemoveFirstMessage implements mu.Agent.
func (b *BossAgent) RemoveFirstMessage() {
	panic("unimplemented")
}

// RemoveLastMessage implements mu.Agent.
func (b *BossAgent) RemoveLastMessage() {
	panic("unimplemented")
}

// RemoveLastNMessages implements mu.Agent.
func (b *BossAgent) RemoveLastNMessages(n int) {
	panic("unimplemented")
}

// ResetMessages implements mu.Agent.
func (b *BossAgent) ResetMessages() {
	panic("unimplemented")
}

// Run implements mu.Agent.
func (b *BossAgent) Run(Messages []openai.ChatCompletionMessageParamUnion) (string, error) {
	panic("unimplemented")
}

// RunStream implements mu.Agent.
func (b *BossAgent) RunStream(Messages []openai.ChatCompletionMessageParamUnion, callBack func(content string) error) (string, error) {
	// NOTE: do not send all messages (history), only the last one
	lastMessage := Messages[len(Messages)-1]

	// Create a task request
	taskRequest := a2a.TaskRequest{
		ID:             fmt.Sprintf("task-%d", time.Now().Unix()),
		JSONRpcVersion: "2.0",
		Method:         "message/send",
		Params: a2a.AgentMessageParams{
			Message: a2a.AgentMessage{
				Role: "user",
				Parts: []a2a.TextPart{
					{
						Text: *lastMessage.GetContent().AsAny().(*string),
						Type: "text",
					},
				},
			},
			MetaData: map[string]any{
				"skill": "ask_for_something",
			},
		},
	}

	fmt.Printf("🚀 Sending streaming task request: %s\n", taskRequest.ID)
	fmt.Println("🌊 Streaming response:")

	// Send the streaming request
	response, err := b.client.SendToAgentStream(taskRequest, callBack)
	if err != nil {
		fmt.Printf("\n❌ Failed to send streaming request: %v\n", err)
		return "", err
	}
	fmt.Println()
	fmt.Printf("✅ Task completed: %s\n", response.ID)
	fmt.Printf("🎯 Status: %s\n", response.Result.Status.State)

	if len(response.Result.History) > 0 {
		fullText := response.Result.History[0].Parts[0].Text
		return fullText, nil
	} else {
		return "", nil
	}

}

// RunStreamWithReasoning implements mu.Agent.
func (b *BossAgent) RunStreamWithReasoning(Messages []openai.ChatCompletionMessageParamUnion, contentCallback func(content string) error, reasoningCallback func(reasoning string) error) (string, string, error) {
	panic("unimplemented")
}

// RunWithReasoning implements mu.Agent.
func (b *BossAgent) RunWithReasoning(Messages []openai.ChatCompletionMessageParamUnion) (string, string, error) {
	panic("unimplemented")
}

// SetDescription implements mu.Agent.
func (b *BossAgent) SetDescription(description string) {
	panic("unimplemented")
}

// SetMessages implements mu.Agent.
func (b *BossAgent) SetMessages(messages []openai.ChatCompletionMessageParamUnion) {
	panic("unimplemented")
}

// SetMetaData implements mu.Agent.
func (b *BossAgent) SetMetaData(metaData any) {
	panic("unimplemented")
}

// SetModel implements mu.Agent.
func (b *BossAgent) SetModel(model shared.ChatModel) {
	panic("unimplemented")
}

// SetName implements mu.Agent.
func (b *BossAgent) SetName(name string) {
	panic("unimplemented")
}

// SetResponseFormat implements mu.Agent.
func (b *BossAgent) SetResponseFormat(format openai.ChatCompletionNewParamsResponseFormatUnion) {
	panic("unimplemented")
}
