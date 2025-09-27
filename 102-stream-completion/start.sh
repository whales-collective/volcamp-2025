#!/bin/bash
export MODEL_RUNNER_BASE_URL="http://localhost:12434/engines/llama.cpp/v1"

#export COOK_MODEL="ai/qwen2.5:3B-F16"
#export COOK_MODEL="ai/qwen2.5:1.5B-F16"
export COOK_MODEL="ai/qwen2.5:0.5B-F16"

#export TEMPERATURE=0.5
export TEMPERATURE=0.2
export TOP_P=0.8

SYSTEM_INSTRUCTION=$(cat <<- EOM
Tu es un chef cuisinier, expert en Aligot et Truffade. 
Tu réponds de manière concise et précise.
EOM
)

export SYSTEM_INSTRUCTION="$SYSTEM_INSTRUCTION"
export USER_PROMPT="Qu'est ce que la truffade ?"
#export USER_PROMPT="Qu'est ce que l'aligot ? ?"

go run main.go
