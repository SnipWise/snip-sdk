#!/bin/bash

# Default configuration
export MODEL_RUNNER_BASE_URL="${MODEL_RUNNER_BASE_URL:-http://localhost:12434/engines/llama.cpp/v1}"
export CHAT_MODEL="${CHAT_MODEL:-hf.co/menlo/lucy-gguf:q4_k_m}"
export COMPRESSOR_MODEL="${COMPRESSOR_MODEL:-hf.co/menlo/jan-nano-128k-gguf:q4_k_m}"
export SYSTEM_INSTRUCTION_PATH="${SYSTEM_INSTRUCTION_PATH:-./system-instructions.md}"
export KNOWLEDGE_BASE_PATH="${KNOWLEDGE_BASE_PATH:-./knowledge-base.md}"

echo "🚀 Starting Sample 52: Context Compression with CompressorAgent"
echo "=================================================="
echo "Configuration:"
echo "  Model Runner URL: $MODEL_RUNNER_BASE_URL"
echo "  Chat Model: $CHAT_MODEL"
echo "  Compressor Model: $COMPRESSOR_MODEL"
echo "  System Instructions: $SYSTEM_INSTRUCTION_PATH"
echo "  Knowledge Base: $KNOWLEDGE_BASE_PATH"
echo "=================================================="
echo ""

go run main.go
