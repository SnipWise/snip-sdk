# Sample 52: Talk to Lucy - Decrease Context Size with Compressor

This sample demonstrates how to use a **CompressorAgent** to intelligently reduce the context size of a conversation while preserving important information.

## Overview

When having long conversations with an AI agent, the context can grow very large, consuming more tokens and potentially hitting model limits. This example shows how to:

1. Build up a conversation with knowledge base context
2. Use a **CompressorAgent** to compress the conversation history
3. Replace the original messages with the compressed version
4. Continue the conversation with reduced context size

## Key Features

- **Two agents working together**: A main chat agent (Bob) and a compressor agent
- **Smart compression**: The compressor agent intelligently summarizes conversations while keeping important facts
- **Context management**: Demonstrates how to replace messages to reduce context size
- **Streaming compression**: Shows real-time compression output

## How It Works

1. **Initial conversation**: Bob answers questions about pizza, including adding a large knowledge base
2. **Context grows**: The conversation history and knowledge base increase the context size
3. **Compression**: The CompressorAgent analyzes all messages and creates a concise summary
4. **Replacement**: The original messages are replaced with:
   - The original system instructions
   - The compressed conversation summary
5. **Continue**: The agent can continue answering questions with much less context

## Usage

```bash
# Run the example
go run main.go

# Or with custom models
MODEL_RUNNER_BASE_URL=http://localhost:12434/engines/llama.cpp/v1 \
CHAT_MODEL=hf.co/menlo/lucy-gguf:q4_k_m \
COMPRESSOR_MODEL=hf.co/menlo/jan-nano-128k-gguf:q4_k_m \
go run main.go
```

## Environment Variables

- `MODEL_RUNNER_BASE_URL`: Base URL for the model server (default: `http://localhost:12434/engines/llama.cpp/v1`)
- `CHAT_MODEL`: Model ID for the main chat agent (default: `hf.co/menlo/lucy-gguf:q4_k_m`)
- `COMPRESSOR_MODEL`: Model ID for the compressor agent (default: `hf.co/menlo/jan-nano-128k-gguf:q4_k_m`)
- `SYSTEM_INSTRUCTION_PATH`: Path to system instructions file (default: `./system-instructions.md`)
- `KNOWLEDGE_BASE_PATH`: Path to knowledge base file (default: `./knowledge-base.md`)

## Benefits of Context Compression

- **Reduced token usage**: Fewer tokens means lower costs and faster responses
- **Avoid context limits**: Prevents hitting model context window limits
- **Preserved information**: Important facts and decisions are retained
- **Better performance**: Smaller context can lead to more focused responses

## Comparison with Sample 51

Sample 51 shows how to manually manage context by replacing messages with new system messages. This sample (52) uses an AI-powered compressor agent to intelligently summarize the conversation, which:

- Automatically identifies important information
- Removes redundancy and fluff
- Preserves technical details and facts
- Creates a coherent summary that maintains conversation flow

## Example Output

The program will show:
1. Original conversation with increasing context size
2. Real-time streaming compression of the conversation
3. Comparison of before/after context sizes
4. Test question to verify the compressed context still works
