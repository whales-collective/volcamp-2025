#!/bin/bash
export MODEL_RUNNER_BASE_URL="http://localhost:12434/engines/llama.cpp/v1"

#export COOK_MODEL="hf.co/menlo/jan-nano-gguf:q4_k_m"
#export COOK_MODEL="hf.co/menlo/lucy-gguf:q8_0"
#export COOK_MODEL="ai/qwen2.5:3B-F16"
#export COOK_MODEL="ai/qwen2.5:1.5B-F16"
export COOK_MODEL="ai/qwen2.5:0.5B-F16"

#export TEMPERATURE=0.5
export TEMPERATURE=0.3
export TOP_P=0.8

export AGENT_NAME="Bibendum"

read -r -d '' SYSTEM_INSTRUCTIONS <<- EOM
Tu es un chef cusinier, expert en Aligot et Truffade. 
Ton nom est ${AGENT_NAME}.
Tu réponds de manière concise et précise.
EOM
export SYSTEM_INSTRUCTIONS="${SYSTEM_INSTRUCTIONS}"

go run main.go
