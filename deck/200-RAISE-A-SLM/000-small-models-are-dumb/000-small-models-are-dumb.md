---
marp: true
theme: default
paginate: true
---
## Small Models are babies 🐣
> You need to explain everything to them

### Let's play with **`ai/qwen2.5:0.5B-F16`**

Qwen2.5 family of models:
| MODEL NAME          | PARAMETERS | QUANTIZATION   | SIZE       |
|---------------------|------------|----------------|------------|
| ai/qwen2.5:latest   | 7.62 B     | IQ2_XXS/Q4_K_M | 4.36 GiB   |
| ai/qwen2.5:3B-F16   | 3.09 B     | F16            | 5.75 GiB   |
| ai/qwen2.5:1.5B-F16 | 1.54 B     | F16            | 2.88 GiB   |
| ai/qwen2.5:0.5B-F16 | 494.03 M   | F16            | 942.43 MiB |

<!--

-->

[← Previous](../../100-DMR/103-stream-completion-go/000-dmr-chat-stream-completion-openai-go.md) | [Next →](001-parameters.md)
