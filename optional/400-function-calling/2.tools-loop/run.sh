#!/bin/bash
export MODEL_RUNNER_BASE_URL="http://localhost:12434/engines/llama.cpp/v1"
export MODEL_RUNNER_LLM_TOOLS="hf.co/menlo/jan-nano-gguf:q4_k_m"

go run main.go