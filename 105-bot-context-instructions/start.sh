#!/bin/bash
export MODEL_RUNNER_BASE_URL="http://localhost:12434/engines/llama.cpp/v1"

#export COOK_MODEL="ai/qwen2.5:3B-F16"
export COOK_MODEL="ai/qwen2.5:1.5B-F16"
#export COOK_MODEL="ai/qwen2.5:0.5B-F16"

#export TEMPERATURE=0.5
export TEMPERATURE=0.3
export TOP_P=0.8

export AGENT_NAME="Bibendum"

export SYSTEM_INSTRUCTIONS_PATH="./bibendum_system_instructions.md"

go run main.go
