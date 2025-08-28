package agents

import (
	"github.com/micro-agent/micro-agent-go/agent/mu"
	"github.com/openai/openai-go/v2"
)

type GhostAgent struct {
	name           string
	messages       []openai.ChatCompletionMessageParamUnion
	responseFormat openai.ChatCompletionNewParamsResponseFormatUnion
}

// GetName implements mu.Agent.
func (g *GhostAgent) GetName() string {
	panic("unimplemented")
}

// SetName implements mu.Agent.
func (g *GhostAgent) SetName(name string) {
	panic("unimplemented")
}

// DetectToolCalls implements mu.Agent.
func (g *GhostAgent) DetectToolCalls(messages []openai.ChatCompletionMessageParamUnion, toolCallBack func(functionName string, arguments string) (string, error)) (string, []string, string, error) {
	panic("unimplemented")
}

// DetectToolCallsStream implements mu.Agent.
func (g *GhostAgent) DetectToolCallsStream(messages []openai.ChatCompletionMessageParamUnion, toolCallback func(functionName string, arguments string) (string, error), streamCallback func(content string) error) (string, []string, string, error) {
	panic("unimplemented")
}

// GenerateEmbeddingVector implements mu.Agent.
func (g *GhostAgent) GenerateEmbeddingVector(content string) ([]float64, error) {
	panic("unimplemented")
}

// GetMessages implements mu.Agent.
func (g *GhostAgent) GetMessages() []openai.ChatCompletionMessageParamUnion {
	panic("unimplemented")
}

// GetResponseFormat implements mu.Agent.
func (g *GhostAgent) GetResponseFormat() openai.ChatCompletionNewParamsResponseFormatUnion {
	panic("unimplemented")
}

// Run implements mu.Agent.
func (g *GhostAgent) Run(Messages []openai.ChatCompletionMessageParamUnion) (string, error) {
	panic("unimplemented")
}

// RunStream implements mu.Agent.
func (g *GhostAgent) RunStream(Messages []openai.ChatCompletionMessageParamUnion, callBack func(content string) error) (string, error) {
	panic("unimplemented")
}

// RunStreamWithReasoning implements mu.Agent.
func (g *GhostAgent) RunStreamWithReasoning(Messages []openai.ChatCompletionMessageParamUnion, contentCallback func(content string) error, reasoningCallback func(reasoning string) error) (string, string, error) {
	panic("unimplemented")
}

// RunWithReasoning implements mu.Agent.
func (g *GhostAgent) RunWithReasoning(Messages []openai.ChatCompletionMessageParamUnion) (string, string, error) {
	panic("unimplemented")
}

// SetMessages implements mu.Agent.
func (g *GhostAgent) SetMessages(messages []openai.ChatCompletionMessageParamUnion) {
	panic("unimplemented")
}

// SetResponseFormat implements mu.Agent.
func (g *GhostAgent) SetResponseFormat(format openai.ChatCompletionNewParamsResponseFormatUnion) {
	g.responseFormat = format
}

// NewFakeAgent creates a new fake agent instance
func NewGhostAgent(name string) mu.Agent {
	return &GhostAgent{
		name:     name,
		messages: []openai.ChatCompletionMessageParamUnion{},
	}
}
