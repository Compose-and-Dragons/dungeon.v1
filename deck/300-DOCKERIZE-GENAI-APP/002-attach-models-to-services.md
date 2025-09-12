---
marp: true
theme: default
paginate: true
---

# How attaching models to services
The Compose way easy and straight forward

---
## `model` attribute of the service definition
- attach one or more models
- choose to configure or not the info passed to the container

```yaml
models:
  # simple definition
  non-player-character-model:
    model: ai/qwen2.5:1.5B-F16

services:
  sorcerer-agent:
  ...
  models:
    - non-player-character-model
```
---
## Customize Env variables
- By default, Compose send variables with names `[SERVICE]_MODEL` & `[SERVICE]_URL` to containers
- You can override those values to use the one you want with `model_var` & `endpoint_var`

```yaml
models:
  # simple definition
  non-player-character-model:
    model: ai/qwen2.5:1.5B-F16

services:
  sorcerer-agent:
  ...
  models:
    non-player-character-model:
      model_var: MY_CUSTOM_MODEL_VAR
      endpoint_var: MY_CUSTOM_URL_VAR
```

---
## How Compose handle interaction with Docker Model Runner

```mermaid
flowchart TD
    Start([Compose up]) --> ParseCompose[Parse Compose files]
    ParseCompose --> ComposeStartService[Start services]
    ComposeStartService --> CheckModels{"Check `models` "}
    CheckModels -->|No| ContinueComposeUP(["Continue up process"])
    CheckModels --> |Yes| CallMDR[Call Model Runner]
    CallMDR --> CheckPull{"Need to pull missing models?"}
    CheckPull --> |No| PullModels[Pull models]
    CheckPull --> |Yes| ConfigureModels[Configure Models]
    PullModels --> ConfigureModels
    ConfigureModels --> SendInfoCompose[Send back model name and url to Compose]
    SendInfoCompose --> PrepareEnvForService[Set env variables for service containers]
    PrepareEnvForService --> CheckCustomEnv{"Custom env variables?"}
    CheckCustomEnv --> |No| SetDefaultEnv["Set [SERVICE]_MODEL & [SERVICE]_URL"]
    CheckCustomEnv --> |Yes| SetCustomEnv["Use model_var & endpoint_var to set env var"]
    SetDefaultEnv --> PassEnvsToContainers[Pass env variables to containers]
    SetCustomEnv --> PassEnvsToContainers
    PassEnvsToContainers --> StartServices[Start service containers]
    StartServices --> ContinueComposeUP
    
    style Start fill:#90EE90
    style ContinueComposeUP fill:#FFB6C1
```