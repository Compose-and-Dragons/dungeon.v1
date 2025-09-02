---
marp: true
theme: default
paginate: true
---
# Chat Stream Completion with Docker Model Runner & **Golang**
> OpenAI Golang SDK
---
## Completion Request

```golang
stream := client.Chat.Completions.NewStreaming(ctx, param)

for stream.Next() {
    chunk := stream.Current()
    // Stream each chunk as it arrives
    if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
        fmt.Print(chunk.Choices[0].Delta.Content)
    }
}
```