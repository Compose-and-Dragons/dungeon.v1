---
marp: true
theme: default
paginate: true
---
# Micro Agent - 🤖 µAgent
> Only a wrapper around the OpenAI Golang SDK
## Especially created for this presentation
---
## OpenAI Client Initialization

```golang
ctx := context.Background()
// Initialize OpenAI client
client := openai.NewClient(
    option.WithBaseURL("http://localhost:12434/engines/llama.cpp/v1"),
    option.WithAPIKey(""),
)
```
---

## Agent and Parameters

```golang
chatAgent, err := mu.NewAgent(ctx, "Bob",
    mu.WithClient(client),
    mu.WithParams(openai.ChatCompletionNewParams{
        Model:       "ai/qwen2.5:1.5B-F16",
        Temperature: openai.Opt(0.0),
    }),
)
```

---

## Messages and Completion Request

```golang
response, err := chatAgent.Run([]openai.ChatCompletionMessageParamUnion{
    openai.SystemMessage("Your name is Bob. You are a helpful AI assistant."),
    openai.UserMessage("Hello what is your name?"),
})

println("Response:", response)
```

[← Previous](docker-mcp-gateway.md) | [Next →](model-context-protocol.md)
