---
marp: true
theme: default
paginate: true
---
# Chat Completion with Docker Model Runner & **Golang**
> OpenAI Golang SDK
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

## Messages and Parameters

```golang
messages := []openai.ChatCompletionMessageParamUnion{
    openai.SystemMessage("You are a useful AI agent expert with TV series."),
    openai.UserMessage("Tell me about the English series called The Avengers?"),
}

param := openai.ChatCompletionNewParams{
    Messages: messages,
    Model:    model,
    Temperature: openai.Opt(0.8),
}
```

---
## Completion Request

```golang
completion, err := client.Chat.Completions.New(ctx, param)

fmt.Println(completion.Choices[0].Message.Content)
```
