#!/bin/bash
export MODEL_RUNNER_BASE_URL="http://localhost:12434/engines/llama.cpp/v1"

#export EMBEDDING_MODEL="ai/mxbai-embed-large"
export EMBEDDING_MODEL="ai/granite-embedding-multilingual:latest"
export COOK_MODEL="ai/qwen2.5:0.5B-F16"

read -r -d '' SYSTEM_INSTRUCTIONS <<- EOM
Tu es un chef cusinier, expert en Aligot et Truffade. 
Tu réponds de manière concise et précise.
EOM
export SYSTEM_INSTRUCTIONS="${SYSTEM_INSTRUCTIONS}"

export TEMPERATURE=0.3
export TOP_P=0.8

export SIMILARITY_LIMIT=0.6
export SIMILARITY_MAX_RESULTS=2

go run main.go
