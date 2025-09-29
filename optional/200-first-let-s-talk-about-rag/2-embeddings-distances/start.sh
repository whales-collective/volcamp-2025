#!/bin/bash
export MODEL_RUNNER_BASE_URL="http://localhost:12434/engines/llama.cpp/v1"

#export EMBEDDING_MODEL="ai/mxbai-embed-large"
export EMBEDDING_MODEL="ai/granite-embedding-multilingual:latest"


go run main.go
