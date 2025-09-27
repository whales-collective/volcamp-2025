#!/bin/bash
#export MODEL_RUNNER_BASE_URL=http://model-runner.docker.internal/engines/llama.cpp/v1
export MODEL_RUNNER_BASE_URL=http://localhost:12434/engines/llama.cpp/v1
EMBEDDING_MODEL=ai/granite-embedding-multilingual:latest
export MCP_HTTP_PORT=9090
export LIMIT=0.6
export MAX_RESULTS=2
#export JSON_STORE_FILE_PATH=store/rag-memory-store.json
export ALIGOT_AGENT_KNOWLEDGE_BASE_PATH=./documents/kb_aligot-enriched.md
export VECTOR_STORES_PATH=./data

go run main.go
