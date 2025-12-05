# Snip SDK (for Docker Model Runner)

> S.N.I.P. Smart Neural Intelligence Partner

A Go SDK for building AI agents with streaming support and HTTP server capabilities. Built on top of Firebase Genkit for seamless integration with LLM providers.

## Features

- **Simple and Remote Agents**: Create local agents or connect to remote agent services
- **Streaming Support**: Real-time streaming responses with `AskStreamWithMemory`
- **HTTP Server**: Built-in HTTP server with multiple endpoints for chat, streaming, and agent management
- **Message History**: Automatic conversation history management
- **System Context**: Dynamic context injection with `AddSystemMessage`
- **Finish Reasons**: Track completion status (stop, length, content_filter, unknown)
- **Type-Safe**: Strongly typed responses with `ChatResponse` structure

## Quick Start

### Basic Agent

```go
ctx := context.Background()

// Configure the agent
agentConfig := snip.AgentConfig{
    Name:               "my-agent",
    SystemInstructions: "You are a helpful assistant",
    ModelID:            "qwen2.5:0.5b",
    EngineURL:          "http://localhost:11434/v1",
}

// Configure the model
modelConfig := snip.ModelConfig{
    Temperature: 0.7,
    MaxTokens:   2000,
}

// Create the agent
agent, err := snip.NewAgent(ctx, agentConfig, modelConfig,
    snip.WithChatFlow(),
    snip.WithChatStreamFlow(),
)

// Ask a question
response, err := agent.AskWithMemory("Hello!")
fmt.Println(response.Text)

// Stream a response
response, err = agent.AskStreamWithMemory("Tell me a story", func(chunk snip.ChatResponse) error {
    fmt.Print(chunk.Text)
    return nil
})

// Check finish reason
if response.IsFinishReasonStop() {
    fmt.Println("Completed normally")
}
```

### Memory Management

The SDK provides two sets of methods for interacting with agents:

**With Memory (Stateful)** - Maintains conversation history:
```go
// Each call remembers previous questions and answers
response1, _ := agent.AskWithMemory("What is Go?")
response2, _ := agent.AskWithMemory("What are its benefits?") // Knows "its" refers to Go

// Streaming with memory
agent.AskStreamWithMemory("Explain more", func(chunk snip.ChatResponse) error {
    fmt.Print(chunk.Text)
    return nil
})
```

**Without Memory (Stateless)** - Each request is independent:
```go
// Each call is independent, no conversation history
response1, _ := agent.Ask("What is Go?")
response2, _ := agent.Ask("What are its benefits?") // Doesn't know what "its" refers to

// Streaming without memory
agent.AskStream("Tell me a story", func(chunk snip.ChatResponse) error {
    fmt.Print(chunk.Text)
    return nil
})
```

### Remote Agent

Connect to a running agent server:

```go
config := snip.ConfigHTTP{
    Address:            "localhost:8080",
    ChatFlowPath:       "/api/chat",
    ChatStreamFlowPath: "/api/chat-stream",
}

remoteAgent := snip.NewRemoteAgent("remote-assistant", config)

response, err := remoteAgent.AskWithMemory("What's the weather?")
```

### HTTP Server

Enable HTTP endpoints for your agent:

```go
agent, err := snip.NewAgent(ctx, agentConfig, modelConfig,
    snip.WithChatFlow(),
    snip.WithChatStreamFlow(),
    snip.EnableServer(snip.ConfigHTTP{
        Address: ":8080",
    }),
)

// Start the server (blocks until shutdown)
agent.Serve()
```

Default endpoints:
- `POST /api/chat` - Non-streaming chat
- `POST /api/chat-stream` - Streaming chat (SSE)
- `GET /api/info` - Agent information
- `POST /api/add-context` - Add system message
- `GET /api/messages` - Get conversation history
- `GET /api/healthcheck` - Health check
- `POST /api/cancel-stream` - Cancel active stream
- `POST /api/shutdown` - Graceful shutdown

### Adding Context Dynamically

```go
// Add system context during conversation
agent.AddSystemMessage("You are now an expert in astronomy")
response, _ := agent.AskWithMemory("Tell me about stars")
```

### Managing Message History

```go
// Get conversation history size (in characters)
size := agent.GetCurrentContextSize()
fmt.Printf("Context has %d characters\n", size)

// Replace entire conversation history
newMessages := []*ai.Message{
    ai.NewSystemTextMessage("You are a helpful assistant"),
    ai.NewUserTextMessage("Hello"),
}
err := agent.ReplaceMessagesWith(newMessages)

// Clear conversation history
agent.ReplaceMessagesWith([]*ai.Message{})
```

## API Reference

### ChatResponse

```go
type ChatResponse struct {
    Text          string // The response text
    FinishReason  string // stop, length, content_filter, unknown
    FinishMessage string // Optional message about finish reason
}
```

Helper methods:
- `IsFinishReasonStop()` - Normal completion
- `IsFinishReasonLength()` - Hit token limit
- `IsFinishReasonContentFilter()` - Content filtered
- `IsFinishReasonUnknown()` - Unknown reason

### AIAgent Interface

```go
type AIAgent interface {
    AskWithMemory(question string) (ChatResponse, error)
    AskStreamWithMemory(question string, callback func(ChatResponse) error) (ChatResponse, error)
    GetName() string
    GetMessages() []*ai.Message
    GetCurrentContextSize() int
    ReplaceMessagesWith(messages []*ai.Message) error
    ReplaceMessagesWithSystemMessages(systemMessages []string) error
    GetInfo() (AgentInfo, error)
    Kind() AgentKind
    AddSystemMessage(context string) error
}
```

## Advanced Usage

### Context Management

The SDK provides powerful tools for managing conversation history:

**Monitor context size:**
```go
// Check context size (in characters) before asking
if agent.GetCurrentContextSize() > 5000 {
    log.Println("Warning: Large context, consider pruning")
}

response, _ := agent.AskWithMemory("What's the weather?")
```

**Implement context window limits:**
```go
const maxContextSize = 10000 // characters

if agent.GetCurrentContextSize() > maxContextSize {
    // Keep only the last 10 messages
    messages := agent.GetMessages()
    recentMessages := messages[len(messages)-10:]
    agent.ReplaceMessagesWith(recentMessages)
}
```

**Save and restore conversations:**
```go
// Save conversation to file
messages := agent.GetMessages()
data, _ := json.Marshal(messages)
os.WriteFile("conversation.json", data, 0644)

// Later, restore the conversation
data, _ := os.ReadFile("conversation.json")
var messages []*ai.Message
json.Unmarshal(data, &messages)
agent.ReplaceMessagesWith(messages)
```

**Start fresh conversation:**
```go
// Clear all history
agent.ReplaceMessagesWith([]*ai.Message{})

// Or start with new system instructions (using ReplaceMessagesWith)
agent.ReplaceMessagesWith([]*ai.Message{
    ai.NewSystemTextMessage("You are a code review expert"),
})

// Easier way: use ReplaceMessagesWithSystemMessages
agent.ReplaceMessagesWithSystemMessages([]string{
    "You are a helpful assistant specialized in French cuisine.",
    "You should always emphasize the importance of using fresh, local ingredients.",
    "You are passionate about traditional cooking techniques.",
})
```

### Remote Agent Limitations

When using `RemoteAgent`, be aware that:
- `ReplaceMessagesWith()` and `ReplaceMessagesWithSystemMessages()` are not supported (message history is managed server-side)
- `GetCurrentContextSize()` makes an HTTP call and may return 0 on errors
- `GetMessages()` fetches from the server, so cache results if needed

## Testing

The project includes comprehensive unit tests for all packages.

### Run All Tests

```bash
# Run all tests
go test ./...

# Run all tests with verbose output
go test ./... -v

# Run tests without cache
go test ./... -count=1
```

### Run Tests for Specific Package

```bash
# Test the snip package
go test ./snip -v

# Test the conversion package
go test ./conversion -v

# Test the files package
go test ./files -v

# Test the env package
go test ./env -v
```

### Run Tests with Coverage

```bash
# Generate coverage report for all packages
go test ./... -cover

# Generate detailed coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Run Benchmarks

```bash
# Run all benchmarks
go test ./... -bench=.

# Run benchmarks for specific package
go test ./env -bench=.
go test ./files -bench=.
go test ./conversion -bench=.
```

### Test Statistics

The project currently includes:
- **215+ unit tests** covering all packages
- **~3500 lines of test code**
- Tests for edge cases, error handling, and integration scenarios
- Comprehensive tests for message management (GetCurrentContextSize, ReplaceMessagesWith)
- Benchmarks for performance-critical functions

## Examples

See the [samples](samples/) directory for example usage:

- [02-stream-completion](samples/02-stream-completion/) - Streaming responses
- [03-both-completions](samples/03-both-completions/) - Regular and streaming
- [04-http-server](samples/04-http-server/) - HTTP server with endpoints
- [05-two-http-servers](samples/05-two-http-servers/) - Multiple agents
- [06-talk-to-agent](samples/06-talk-to-agent/) - Interactive conversation
- [07-add-context](samples/07-add-context/) - Dynamic context injection
- [08-add-context-to-remote](samples/08-add-context-to-remote/) - Remote agent context
- [50-talk-to-lucy](samples/50-talk-to-lucy/) - Advanced example

## License

See LICENSE file for details.

