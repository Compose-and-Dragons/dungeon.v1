#!/bin/bash
export MODEL_RUNNER_BASE_URL="http://localhost:12434/engines/llama.cpp/v1"

#export SORCERER_MODEL="ai/qwen2.5:1.5B-F16"
export SORCERER_MODEL="ai/qwen2.5:0.5B-F16"

# The chat agent
export SORCERER_NAME="Elara"
export SORCERER_MODEL_TEMPERATURE=0.5
export SORCERER_SYSTEM_INSTRUCTIONS_PATH="./data/sorcerer_system_instructions.md"
export SORCERER_CONTEXT_PATH="./data/sorcerer_background_and_personality.md"

# The RAG configuration for the embedding agent
export EMBEDDING_MODEL="ai/mxbai-embed-large:latest"
export SIMILARITY_LIMIT=0.5
export SIMILARITY_MAX_RESULTS=2
export VECTOR_STORES_PATH="./data"


go run main.go