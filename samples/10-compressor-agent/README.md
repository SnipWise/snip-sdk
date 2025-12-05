# Compressor Agent Example

This example demonstrates how to use the `CompressorAgent` to compress long conversation histories while preserving essential information.

## What is CompressorAgent?

The `CompressorAgent` is a specialized agent that uses AI to intelligently compress text by:
- Removing redundancy and conversational fluff
- Preserving critical information (facts, code, decisions)
- Maintaining chronological flow
- Summarizing long discussions into concise bullet points

## Use Cases

- **Reduce context window size**: When your conversation history is getting too long
- **Optimize token usage**: Compress before continuing a conversation to save tokens
- **Summarize discussions**: Extract key points from verbose exchanges
- **Maintain continuity**: Keep conversation context manageable in long-running sessions

## Features Demonstrated

1. **Non-streaming compression**: Get the complete compressed result at once
2. **Streaming compression**: Watch the compression happen in real-time
3. **Compression metrics**: See the reduction in text length

## Running the Example

```bash
# Make sure your model server is running
# For example, with Docker Model Runner on port 12434

# Run the example
go run main.go
```

## Expected Output

The example will:
1. Compress a long conversation about Go authentication (non-streaming)
2. Show compression statistics (original length, compressed length, ratio)
3. Display the compressed text
4. Compress another conversation about database optimization (streaming)
5. Show real-time compression as it happens

## Configuration

The example uses these environment variables (with defaults):

- `MODEL_RUNNER_BASE_URL`: Model server URL (default: `http://localhost:12434/engines/llama.cpp/v1`)
- `CHAT_MODEL`: Model to use (default: `ai/qwen2.5:latest`)

## How It Works

The CompressorAgent uses a specialized system prompt that instructs the AI to:

1. **Preserve Critical Information**:
   - Facts, decisions, code snippets
   - File paths, function names, technical details

2. **Remove Redundancy**:
   - Greetings and acknowledgments
   - Repetitive discussions
   - Failed attempts and verbose explanations

3. **Maintain Structure**:
   - Logical flow and chronology
   - Context for continuing the conversation

4. **Output Format**:
   - Clear, concise language
   - Grouped related topics
   - Highlighted key decisions and outcomes

## Example Compression

**Before (verbose conversation):**
```
User: Hello, I'm working on a Go project...
Agent: That's great! For user authentication in Go, you have several options...
User: It's a REST API for a mobile app.
Agent: Perfect! For a REST API serving a mobile app, I'd recommend...
[... many more exchanges ...]
```

**After (compressed):**
```
Discussion about implementing JWT authentication for a Go REST API mobile backend.
Key points:
- JWT recommended for stateless mobile API authentication
- Use jwt-go package (github.com/golang-jwt/jwt/v5)
- Implement middleware for token validation
- Use bcrypt for password hashing (crypto/bcrypt)
- Methods: GenerateFromPassword() and CompareHashAndPassword()
```

## Tips

- Use **lower temperature** (0.3-0.5) for more consistent compression
- The compression ratio typically ranges from 40-70% depending on content
- Streaming is useful for very long texts to see progress
- The compressed output is designed to maintain conversation continuity
