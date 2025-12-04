# Smart Package

The `smart` package provides AI agent functionality for building conversational applications with LLM models. It simplifies the creation and management of AI agents with persistent conversation history.

## Overview

This package offers a simple way to create AI agents that can:
- Have conversations with memory (conversation history)
- Stream responses in real-time
- Be exposed as HTTP services
- Work with OpenAI-compatible APIs (OpenAI, Ollama, etc.)

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    "github.com/snipwise/snip-sdk/smart"
)

func main() {
    ctx := context.Background()

    // Create a simple agent
    agent := smart.NewAgent(
        ctx,
        "my-assistant",                                        // Agent name
        "You are a helpful AI assistant",                      // System instructions
        "hf.co/menlo/jan-nano-gguf:q4_k_m",                    // Model ID
        "http://localhost:12434/engines/llama.cpp/v1",         // Engine URL (Docker Model Runner)
        smart.Config{
            Temperature: 0.7,
            MaxTokens:   500,
        },
        smart.EnableChatFlowWithMemory(),                      // Enable basic chat with memory
    )

    // Ask a question
    answer, err := agent.Ask("What is the capital of France?")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Answer:", answer)
}
```

### Streaming Responses

```go
package main

import (
    "context"
    "fmt"
    "log"
    "github.com/snipwise/snip-sdk/smart"
)

func main() {
    ctx := context.Background()

    agent := smart.NewAgent(
        ctx,
        "streaming-assistant",
        "You are a helpful AI assistant",
        "hf.co/menlo/jan-nano-gguf:q4_k_m",
        "http://localhost:12434/engines/llama.cpp/v1",  // Docker Model Runner
        smart.Config{Temperature: 0.7},
        smart.EnableChatStreamFlowWithMemory(),  // Enable streaming with memory
    )

    // Stream the response
    _, err := agent.AskStream("Tell me a short story", func(chunk string) error {
        fmt.Print(chunk)  // Print each chunk as it arrives
        return nil
    })

    if err != nil {
        log.Fatal(err)
    }
}
```

### Agent as HTTP Server

```go
package main

import (
    "context"
    "log"
    "github.com/snipwise/snip-sdk/smart"
)

func main() {
    ctx := context.Background()

    agent := smart.NewAgent(
        ctx,
        "http-agent",
        "You are a helpful AI assistant",
        "hf.co/menlo/jan-nano-gguf:q4_k_m",
        "http://localhost:12434/engines/llama.cpp/v1",  // Docker Model Runner
        smart.Config{
            Temperature: 0.7,
            TopP:        0.9,
        },
        smart.EnableChatFlowWithMemory(),
        smart.EnableChatStreamFlowWithMemory(),
        smart.EnableServer(smart.ConfigHTTP{
            Address:            "0.0.0.0:8080",
            ChatFlowPath:       "/api/chat",
            ChatStreamFlowPath: "/api/chat-stream",
            ShutdownPath:       "/server/shutdown",
        }),
    )

    // Start the server (blocks until shutdown)
    if err := agent.Serve(); err != nil {
        log.Fatal(err)
    }
}
```

Once running, you can interact with the agent via HTTP:

```bash
# Ask a question
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello, how are you?"}'

# Stream a response
curl -X POST http://localhost:8080/api/chat-stream \
  -H "Content-Type: application/json" \
  -d '{"message": "Tell me a story"}' \
  --no-buffer

# Get agent information
curl http://localhost:8080/information

# Health check
curl http://localhost:8080/healthcheck

# Shutdown server
curl -X POST http://localhost:8080/server/shutdown
```

## Core Types

### Agent

The main structure representing an AI agent.

```go
type Agent struct {
    Name               string
    SystemInstructions string
    ModelID            string
    Config             Config
    Messages           []*ai.Message
}
```

### Config

Configuration for the LLM model behavior.

```go
type Config struct {
    Temperature      float64   // Randomness (0.0-2.0)
    TopP             float64   // Nucleus sampling (0.0-1.0)
    MaxTokens        int64     // Maximum tokens to generate
    FrequencyPenalty float64   // Reduce repetition (-2.0 to 2.0)
    PresencePenalty  float64   // Encourage new topics (-2.0 to 2.0)
    Stop             []string  // Stop sequences
    Seed             *int64    // For deterministic sampling
}
```

### AIAgent Interface

```go
type AIAgent interface {
    Ask(question string) (string, error)
    AskStream(question string, callback func(string) error) (string, error)
    GetName() string
    GetMessages() []*ai.Message
    GetInfo() (AgentInfo, error)
    Kind() AgentKind
    AddSystemMessage(context string) error
}
```

## Agent Options

Configure your agent using functional options:

```go
// Enable basic chat with conversation memory (recommended)
smart.EnableChatFlowWithMemory()

// Enable streaming chat with conversation memory (recommended)
smart.EnableChatStreamFlowWithMemory()

// Enable HTTP server with ConfigHTTP struct
smart.EnableServer(smart.ConfigHTTP{
    Address:            "0.0.0.0:8080",        // Required: server address
    ChatFlowPath:       "/api/chat",           // Optional: defaults to "/api/chat"
    ChatStreamFlowPath: "/api/chat-stream",    // Optional: defaults to "/api/chat-stream"
    ShutdownPath:       "/server/shutdown",    // Optional: defaults to "-" (disabled)
    InformationPath:    "/api/information",    // Optional: defaults to "/api/information"
    HealthcheckPath:    "/healthcheck",        // Optional: defaults to "/healthcheck"
    CancelStreamPath:   "/api/cancel-stream",  // Optional: defaults to "/api/cancel-stream-completion"
    AddContextPath:     "/api/add-context",    // Optional: defaults to "/api/add-system-message"
    GetMessagesPath:    "/api/messages",       // Optional: defaults to "/api/messages"
})
```

### Default Values

All endpoint paths have default values. You can omit any field to use the default:

| Field | Default Value | Description |
|-------|--------------|-------------|
| `ChatFlowPath` | `/api/chat` | Standard chat endpoint |
| `ChatStreamFlowPath` | `/api/chat-stream` | Streaming chat endpoint |
| `InformationPath` | `/api/information` | Agent information endpoint |
| `HealthcheckPath` | `/healthcheck` | Health check endpoint |
| `ShutdownPath` | `-` (disabled) | Server shutdown endpoint. Set to a path like `/server/shutdown` to enable |
| `CancelStreamPath` | `/api/cancel-stream-completion` | Cancel streaming endpoint |
| `AddContextPath` | `/api/add-system-message` | Add context/system message endpoint |
| `GetMessagesPath` | `/api/messages` | Get conversation history endpoint |

**Minimal server configuration:**
```go
// Use all defaults except address
smart.EnableServer(smart.ConfigHTTP{
    Address: "0.0.0.0:8080",
})
```

**Note:**
- The "WithMemory" variants maintain conversation history, which is essential for multi-turn conversations. Use these for most applications.
- Currently, all "Chat Flow" implementations store conversation history in memory (RAM). Persistence to database or other storage is not yet implemented.

## Agent Methods

### Ask

Send a message and get a complete response.

```go
answer, err := agent.Ask("What is 2+2?")
```

### AskStream

Stream the response chunk by chunk.

```go
fullAnswer, err := agent.AskStream("Explain AI", func(chunk string) error {
    fmt.Print(chunk)
    return nil
})
```

### AddSystemMessage

Add context to the conversation.

```go
err := agent.AddSystemMessage("The user prefers concise answers.")
```

### GetMessages

Get the conversation history.

```go
messages := agent.GetMessages()
for _, msg := range messages {
    fmt.Printf("%s: %s\n", msg.Role, msg.Content[0].Text)
}
```

### GetInfo

Get agent information.

```go
info, err := agent.GetInfo()
fmt.Printf("Agent: %s, Model: %s\n", info.Name, info.ModelID)
```

## HTTP Server Endpoints

When `EnableServer()` is used, the following endpoints are available (with default paths):

| Method | Default Path | Description |
|--------|--------------|-------------|
| `POST` | `/api/chat` | Send a message, get complete response |
| `POST` | `/api/chat-stream` | Send a message, stream response |
| `GET` | `/api/information` | Get agent information |
| `GET` | `/healthcheck` | Health check endpoint |
| `GET` | `/api/messages` | Get conversation history |
| `POST` | `/api/add-system-message` | Add system message context |
| `POST` | `/api/cancel-stream-completion` | Cancel ongoing stream |
| `POST` | Disabled by default | Gracefully shutdown the server (set `ShutdownPath` to enable) |

**Note:** All paths are configurable via `ConfigHTTP`. See the [Default Values](#default-values) section above for details.

## Complete Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "github.com/snipwise/snip-sdk/smart"
)

func main() {
    ctx := context.Background()

    // Create an agent with custom configuration
    agent := smart.NewAgent(
        ctx,
        "tutor-bot",
        "You are a patient and knowledgeable tutor who explains concepts clearly.",
        "hf.co/menlo/jan-nano-gguf:q4_k_m",
        "http://localhost:12434/engines/llama.cpp/v1",  // Docker Model Runner
        smart.Config{
            Temperature:      0.8,
            MaxTokens:        1000,
            FrequencyPenalty: 0.5,
        },
        smart.EnableChatFlowWithMemory(),
        smart.EnableChatStreamFlowWithMemory(),
    )

    // Add context
    agent.AddSystemMessage("The student is learning Go programming.")

    // Have a conversation
    questions := []string{
        "What is a goroutine?",
        "Can you give me a simple example?",
    }

    for _, question := range questions {
        fmt.Printf("\n🧑 Question: %s\n", question)
        fmt.Print("🤖 Answer: ")

        _, err := agent.AskStream(question, func(chunk string) error {
            fmt.Print(chunk)
            return nil
        })

        if err != nil {
            log.Fatal(err)
        }
        fmt.Println()
    }

    // View conversation history
    messages := agent.GetMessages()
    fmt.Printf("\n📚 Total messages in history: %d\n", len(messages))
}
```

## Configuration Tips

### Temperature
- `0.0-0.3`: Focused, deterministic responses
- `0.4-0.7`: Balanced creativity and consistency
- `0.8-1.0`: More creative and varied responses

### MaxTokens
- Short answers: `100-300`
- Medium responses: `500-1000`
- Long-form content: `1000-4000`

## Error Handling

```go
answer, err := agent.Ask("Hello")
if err != nil {
    log.Printf("Error: %v", err)
    // Handle error appropriately
}
```

## Installation

```bash
go get github.com/snipwise/snip-sdk/smart
```

## Dependencies

This package uses:
- [Firebase Genkit](https://github.com/firebase/genkit) - AI framework
- [OpenAI Go SDK](https://github.com/openai/openai-go) - OpenAI API client

## License

See the main repository for license information.
