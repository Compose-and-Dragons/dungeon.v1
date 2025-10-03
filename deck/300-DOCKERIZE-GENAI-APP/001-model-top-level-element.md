---
marp: true
theme: default
paginate: true
---

# The top level element `models`
Define and tune your models for your containerized app

---
## How to define models
- you can define a list of models
- a model definition is easy
  * a name
  * `model`: OCI artifact name of your model
  * `context_size`: maximum token context size for the model
  * `runtime_flags`: raw flags passed to the inference runtime such as temperature, verbose mode...

---
## How to define models
```yaml
models:
  # simple definition
  embedding-model:
    model: ai/granite-embedding-multilingual:latest
    
  # more detailed definition
  non-player-character-model:
    model: ai/qwen2.5:1.5B-F16
    context_size: 4096
    runtime_flags: 
      - "--temp"                # Temperature
      - "0.1"
      - "--verbose-prompt"
```

[← Previous](000-agentic-compose.md) | [Next →](002-attach-models-to-services.md)
