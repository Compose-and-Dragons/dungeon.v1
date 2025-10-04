#!/bin/bash
export MODEL_RUNNER_BASE_URL="http://localhost:12434/engines/llama.cpp/v1"
#export MODEL_RUNNER_LLM_TOOLS="hf.co/menlo/jan-nano-gguf:q4_k_m"
# ---------------------------------------------------------
# Boss settings
# ---------------------------------------------------------
export BOSS_MODEL="hf.co/menlo/lucy-gguf:q8_0"
export EMBEDDING_MODEL="ai/granite-embedding-multilingual:latest"
export BOSS_REMOTE_AGENT_HTTP_PORT=8888
export BOSS_MODEL_TEMPERATURE=0.7
export BOSS_NAME="Shesepankh"
export BOSS_RACE="Sphinx"
export BOSS_DESCRIPTION="An ancient and wise Sphinx who guards the exit of the dungeon. Known for posing riddles to those who seek passage."
export BOSS_SYSTEM_INSTRUCTIONS_PATH="./data/boss_system_instructions.md"
export BOSS_CONTEXT_PATH="./data/boss_background_and_personality.md"
# ---------------------------------------------------------
# Similarity search settings
# ---------------------------------------------------------
export SIMILARITY_LIMIT=0.55
export SIMILARITY_MAX_RESULTS=2
export VECTOR_STORES_PATH="./data"

go run main.go







