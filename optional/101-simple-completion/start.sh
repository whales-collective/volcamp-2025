#!/bin/bash
export MODEL_RUNNER_BASE_URL="http://localhost:12434/engines/llama.cpp/v1"

#export COOK_MODEL="ai/qwen2.5:1.5B-F16"
export COOK_MODEL="ai/qwen2.5:0.5B-F16"

go run main.go