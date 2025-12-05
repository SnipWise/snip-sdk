# Sample 11: Compressor Agent - Message Compression

This example demonstrates how to use the `CompressorAgent` to compress conversation messages while preserving essential information.

## Features Demonstrated

- **Message Compression**: Compress `[]*ai.Message` objects (conversation history)
- **Non-streaming Compression**: Get the complete compressed result at once using `CompressMessages()`
- **Streaming Compression**: Watch the compression happen in real-time using `CompressMessagesStream()`
- **Context Size Reduction**: Reduce the size of conversation history while maintaining key information

## Use Cases

The CompressorAgent is useful for:

1. **Memory Management**: Reduce context size when conversation history grows too large
2. **Token Optimization**: Save tokens when working with LLMs that have token limits
3. **Conversation Summarization**: Create concise summaries of long conversations
4. **Context Window Management**: Keep conversations within model context limits

## How It Works

The CompressorAgent:
1. Takes a list of `ai.Message` objects (conversation history)
2. Converts them to a formatted text representation
3. Applies intelligent compression using the compression prompt
4. Removes redundancy while preserving:
   - Important facts and decisions
   - Code snippets and technical details
   - File paths and function names
   - Key takeaways and outcomes

## Running the Example

```bash
cd samples/11-compressor-agent
go run main.go
```

## Expected Output

The example will:
1. Display the original conversation messages
2. Show non-streaming compression with statistics
3. Demonstrate streaming compression in real-time
4. Display compression ratios and size comparisons

## Comparison with Sample 10

- **Sample 10**: Compresses raw text strings using `CompressText()` and `CompressTextStream()`
- **Sample 11**: Compresses `ai.Message` objects using `CompressMessages()` and `CompressMessagesStream()`

Both samples use the same underlying compression logic but handle different input types.
