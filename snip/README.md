# Snip Package

The `snip` package provides AI agent functionality for building conversational applications with LLM models. It simplifies the creation and management of AI agents with persistent conversation history.

## Overview

This package offers a simple way to create AI agents that can:
- Have conversations with memory (conversation history)
- Stream responses in real-time
- Be exposed as HTTP services
- Work with OpenAI-compatible APIs (OpenAI, Ollama, etc.)
- Perform semantic search with RAG (Retrieval-Augmented Generation)
- Verify model availability before agent creation

## Important Changes

### Breaking Changes in Recent Updates

**NewAgent now returns an error:**

The `NewAgent` function signature has changed to include error handling and automatic model availability verification:

```go
// Old (deprecated)
agent := snip.NewAgent(ctx, agentConfig, modelConfig, opts...)

// New (current)
agent, err := snip.NewAgent(ctx, agentConfig, modelConfig, opts...)
if err != nil {
    log.Fatal(err)  // Handle model unavailability or other errors
}
```

**What this means:**
- The agent creation now verifies that the specified model is available at the engine URL
- If the model is not available, an error is returned immediately
- This prevents runtime errors and provides clearer feedback during initialization

**Migration:**
Update all `NewAgent` calls to handle the returned error:

```go
// Before
agent := snip.NewAgent(ctx, agentConfig, modelConfig, snip.EnableChatFlowWithMemory())

// After
agent, err := snip.NewAgent(ctx, agentConfig, modelConfig, snip.EnableChatFlowWithMemory())
if err != nil {
    log.Fatalf("Failed to create agent: %v", err)
}
```

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    "github.com/snipwise/snip-sdk/snip"
)

func main() {
    ctx := context.Background()

    // Create a simple agent
    agent, err := snip.NewAgent(
        ctx,
        snip.AgentConfig{
            Name:               "my-assistant",
            SystemInstructions: "You are a helpful AI assistant",
            ModelID:            "hf.co/menlo/jan-nano-gguf:q4_k_m",
            EngineURL:          "http://localhost:12434/engines/llama.cpp/v1",
        },
        snip.ModelConfig{
            Temperature: 0.7,
            MaxTokens:   500,
        },
        snip.EnableChatFlowWithMemory(),                      // Enable basic chat with memory
    )
    if err != nil {
        log.Fatal(err)
    }

    // Ask a question
    response, err := agent.Ask("What is the capital of France?")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Answer:", response.Text)
}
```

### Streaming Responses

```go
package main

import (
    "context"
    "fmt"
    "log"
    "github.com/snipwise/snip-sdk/snip"
)

func main() {
    ctx := context.Background()

    agent, err := snip.NewAgent(
        ctx,
        snip.AgentConfig{
            Name:               "streaming-assistant",
            SystemInstructions: "You are a helpful AI assistant",
            ModelID:            "hf.co/menlo/jan-nano-gguf:q4_k_m",
            EngineURL:          "http://localhost:12434/engines/llama.cpp/v1",
        },
        snip.ModelConfig{Temperature: 0.7},
        snip.EnableChatStreamFlowWithMemory(),  // Enable streaming with memory
    )
    if err != nil {
        log.Fatal(err)
    }

    // Stream the response
    _, err = agent.AskStream("Tell me a short story", func(chunk ChatResponse) error {
        fmt.Print(chunk.Text)  // Print each chunk as it arrives
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
    "github.com/snipwise/snip-sdk/snip"
)

func main() {
    ctx := context.Background()

    agent, err := snip.NewAgent(
        ctx,
        snip.AgentConfig{
            Name:               "http-agent",
            SystemInstructions: "You are a helpful AI assistant",
            ModelID:            "hf.co/menlo/jan-nano-gguf:q4_k_m",
            EngineURL:          "http://localhost:12434/engines/llama.cpp/v1",
        },
        snip.ModelConfig{
            Temperature: 0.7,
            TopP:        0.9,
        },
        snip.EnableChatFlowWithMemory(),
        snip.EnableChatStreamFlowWithMemory(),
        snip.EnableServer(snip.ConfigHTTP{
            Address:            "0.0.0.0:8080",
            ChatFlowPath:       "/api/chat",
            ChatStreamFlowPath: "/api/chat-stream",
            ShutdownPath:       "/server/shutdown",
        }),
    )
    if err != nil {
        log.Fatal(err)
    }

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

### RAG Agent (Semantic Search)

Create an agent for Retrieval-Augmented Generation (RAG) with semantic search capabilities.

```go
package main

import (
    "context"
    "fmt"
    "log"
    "github.com/snipwise/snip-sdk/snip"
)

func main() {
    ctx := context.Background()

    // Create a RAG agent with embedding model
    ragAgent, err := snip.NewRagAgent(
        ctx,
        snip.RagAgentConfig{
            Name:      "rag-assistant",
            ModelID:   "ai/mxbai-embed-large",
            EngineURL: "http://localhost:12434/engines/llama.cpp/v1",
        },
        snip.StoreConfig{
            StoreName: "my-documents",
            StorePath: "./data",
        },
    )
    if err != nil {
        log.Fatal(err)
    }

    // Add documents to the store
    chunks := []snip.TextChunk{
        {
            Content:  "Squirrels run in the forest",
            Metadata: map[string]any{"source": "nature.txt"},
        },
        {
            Content:  "Dolphins leap out of the ocean",
            Metadata: map[string]any{"source": "marine.txt"},
        },
    }

    count, err := ragAgent.AddTextChunksToStore(chunks)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Indexed %d documents\n", count)

    // Search for similar documents
    similarities, err := ragAgent.SearchSimilarities("Which animals swim?")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Similar documents:")
    for _, doc := range similarities {
        fmt.Println("  -", doc)
    }
}
```

### Remote Agent (Client)

Connect to a remote agent server and interact with it programmatically.

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    "github.com/snipwise/snip-sdk/snip"
)

func main() {
    ctx := context.Background()
    engineURL := "http://localhost:12434/engines/llama.cpp/v1"
    chatModelId := "hf.co/menlo/jan-nano-gguf:q4_k_m"

    // Create a local agent with HTTP server
    agent, err := snip.NewAgent(
        ctx,
        snip.AgentConfig{
            Name:               "Server Agent",
            SystemInstructions: "You are a helpful assistant.",
            ModelID:            chatModelId,
            EngineURL:          engineURL,
        },
        snip.ModelConfig{
            Temperature: 0.7,
            TopP:        0.9,
        },
        snip.EnableChatFlowWithMemory(),
        snip.EnableChatStreamFlowWithMemory(),
        snip.EnableServer(snip.ConfigHTTP{
            Address:            "0.0.0.0:9100",
            ChatFlowPath:       "/api/chat",
            ChatStreamFlowPath: "/api/chat-stream",
            InformationPath:    "/api/information",
        }),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Start server in background
    go func() {
        if err := agent.Serve(); err != nil {
            log.Fatalf("Server error: %v", err)
        }
    }()

    // Wait for server to start
    time.Sleep(2 * time.Second)

    // Create a remote agent client
    remoteAgent := snip.NewRemoteAgent(
        "Remote Client",
        snip.ConfigHTTP{
            Address:            "0.0.0.0:9100",
            ChatFlowPath:       "/api/chat",
            ChatStreamFlowPath: "/api/chat-stream",
            InformationPath:    "/api/information",
        },
    )

    // Get agent information
    info, err := remoteAgent.GetInfo()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Connected to: %s (Model: %s)\n", info.Name, info.ModelID)

    // Ask a question
    response, err := remoteAgent.Ask("What is Go?")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Answer:", response.Text)

    // Stream a response
    _, err = remoteAgent.AskStream("Explain goroutines", func(chunk ChatResponse) error {
        fmt.Print(chunk.Text)
        return nil
    })
    if err != nil {
        log.Fatal(err)
    }

    // Add context
    err = remoteAgent.AddSystemMessage("The user is a beginner programmer.")
    if err != nil {
        log.Fatal(err)
    }

    // Get conversation history
    messages := remoteAgent.GetMessages()
    fmt.Printf("\nTotal messages: %d\n", len(messages))
}
```

## Model Availability

The package includes helper functions to check model availability before creating agents.

### GetModelsList

Retrieve a list of available models from the inference engine.

```go
ctx := context.Background()
models, err := snip.GetModelsList(ctx, "http://localhost:12434/engines/llama.cpp/v1")
if err != nil {
    log.Fatal(err)
}

fmt.Println("Available models:")
for _, model := range models {
    fmt.Println("  -", model)
}
```

### IsModelAvailable

Check if a specific model is available.

```go
ctx := context.Background()
engineURL := "http://localhost:12434/engines/llama.cpp/v1"
modelID := "hf.co/menlo/jan-nano-gguf:q4_k_m"

if snip.IsModelAvailable(ctx, engineURL, modelID) {
    fmt.Printf("Model %s is available\n", modelID)
} else {
    fmt.Printf("Model %s is not available\n", modelID)
}
```

**Note:** `NewAgent` automatically verifies model availability when creating an agent. If the model is not available, it returns an error.

## Core Types

### Agent

The main structure representing an AI agent.

```go
type Agent struct {
    Name               string
    SystemInstructions string
    ModelID            string
    Config             ModelConfig
    Messages           []*ai.Message
}
```

### AgentConfig

Configuration for creating an agent.

```go
type AgentConfig struct {
    Name               string  // Agent identifier
    SystemInstructions string  // Agent's behavior and role
    ModelID            string  // Language model to use
    EngineURL          string  // Model inference engine base URL
}
```

### ModelConfig

Configuration for the LLM model behavior.

```go
type ModelConfig struct {
    Temperature      float64   // Randomness (0.0-2.0)
    TopP             float64   // Nucleus sampling (0.0-1.0)
    MaxTokens        int64     // Maximum tokens to generate
    FrequencyPenalty float64   // Reduce repetition (-2.0 to 2.0)
    PresencePenalty  float64   // Encourage new topics (-2.0 to 2.0)
    Stop             []string  // Stop sequences
    Seed             *int64    // For deterministic sampling
}
```

### RagAgent

Agent for Retrieval-Augmented Generation with semantic search.

```go
type RagAgent struct {
    ctx    context.Context
    Name   string
    ModelID string

    storeName string
    storePath string

    genKitInstance    *genkit.Genkit
    embedder          ai.Embedder
    docStore          *localvec.DocStore
    documentRetriever ai.Retriever

    embeddingDimension int
}
```

**Configuration Types:**

```go
type RagAgentConfig struct {
    Name      string  // Agent identifier
    ModelID   string  // Embedding model to use
    EngineURL string  // Model inference engine base URL
}

type StoreConfig struct {
    StoreName string  // Name of the document store
    StorePath string  // Path to store documents
}

type TextChunk struct {
    Content  string         // Text content
    Metadata map[string]any // Optional metadata
}
```

**Constructor:**

```go
func NewRagAgent(ctx context.Context, ragAgentConfig RagAgentConfig, storeConfig StoreConfig, opts ...RagAgentOption) (*RagAgent, error)
```

**Example:**

```go
ragAgent, err := snip.NewRagAgent(
    ctx,
    snip.RagAgentConfig{
        Name:      "my-rag-agent",
        ModelID:   "ai/mxbai-embed-large",
        EngineURL: "http://localhost:12434/engines/llama.cpp/v1",
    },
    snip.StoreConfig{
        StoreName: "documents",
        StorePath: "./data",
    },
)
if err != nil {
    log.Fatal(err)
}
```

### RemoteAgent

Client for connecting to a remote agent server.

```go
type RemoteAgent struct {
    Name                string  // Client identifier
    ChatStreamEndpoint  string  // Full URL for streaming chat
    ChatEndPoint        string  // Full URL for standard chat
    InformationEndpoint string  // Full URL for agent information
    AddContextEndpoint  string  // Full URL for adding context
    GetMessagesEndpoint string  // Full URL for getting messages
}
```

**Constructor:**

```go
func NewRemoteAgent(name string, config ConfigHTTP) *RemoteAgent
```

**Example:**

```go
remoteAgent := snip.NewRemoteAgent(
    "My Remote Client",
    snip.ConfigHTTP{
        Address:            "0.0.0.0:9100",
        ChatFlowPath:       "/api/chat",
        ChatStreamFlowPath: "/api/chat-stream",
        InformationPath:    "/api/information",
    },
)
```

### ChatResponse

The response structure returned by `Ask` and `AskStream` methods.

```go
type ChatResponse struct {
    Text          string // The response text
    FinishReason  string // Completion reason: stop, length, content_filter, unknown
    FinishMessage string // Optional message about finish reason
}
```

**Helper Methods:**

```go
func (chatResponse *ChatResponse) IsFinishReasonStop() bool
func (chatResponse *ChatResponse) IsFinishReasonLength() bool
func (chatResponse *ChatResponse) IsFinishReasonContentFilter() bool
func (chatResponse *ChatResponse) IsFinishReasonUnknown() bool
```

**Example:**

```go
response, err := agent.Ask("What is Go?")
if err != nil {
    log.Fatal(err)
}

fmt.Println(response.Text)

// Check finish reason
if response.IsFinishReasonStop() {
    fmt.Println("Completed normally")
} else if response.IsFinishReasonLength() {
    fmt.Println("Hit token limit")
}
```

### AIAgent Interface

Both `Agent` and `RemoteAgent` implement the `AIAgent` interface.

```go
type AIAgent interface {
    Ask(question string) (ChatResponse, error)
    AskStream(question string, callback func(ChatResponse) error) (ChatResponse, error)
    GetName() string
    GetMessages() []*ai.Message
    GetCurrentContextSize() int
    ReplaceMessagesWith(messages []*ai.Message) error
    GetInfo() (AgentInfo, error)
    Kind() AgentKind
    AddSystemMessage(context string) error
}
```

### AIRagAgent Interface

`RagAgent` implements the `AIRagAgent` interface for semantic search operations.

```go
type AIRagAgent interface {
    GetName() string
    GetInfo() (RagAgentInfo, error)
    Kind() AgentKind
    AddTextChunksToStore(chunks []TextChunk) (int, error)
    SearchSimilarities(query string) ([]string, error)
}
```

### AgentKind

Agent types are identified by the `AgentKind` constant.

```go
type AgentKind string

const (
    Basic  AgentKind = "Basic"
    Remote AgentKind = "Remote"
    Tool   AgentKind = "Tool"    // Reserved for future use
    Intent AgentKind = "Intent"  // Reserved for future use
    Rag    AgentKind = "Rag"
)
```

## Agent Options

Configure your agent using functional options:

```go
// Enable basic chat with conversation memory (recommended)
snip.EnableChatFlowWithMemory()

// Enable streaming chat with conversation memory (recommended)
snip.EnableChatStreamFlowWithMemory()

// Enable HTTP server with ConfigHTTP struct
snip.EnableServer(snip.ConfigHTTP{
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
snip.EnableServer(snip.ConfigHTTP{
    Address: "0.0.0.0:8080",
})
```

**Note:**
- The "WithMemory" variants maintain conversation history, which is essential for multi-turn conversations. Use these for most applications.
- Currently, all "Chat Flow" implementations store conversation history in memory (RAM). Persistence to database or other storage is not yet implemented.

## Agent Methods

### Standard Agent Methods

These methods are available for both `Agent` and `RemoteAgent` (they implement the `AIAgent` interface).

### Ask

Send a message and get a complete response.

```go
// Local agent
response, err := agent.Ask("What is 2+2?")
if err != nil {
    log.Fatal(err)
}
fmt.Println(response.Text)

// Check finish reason
if response.IsFinishReasonStop() {
    fmt.Println("Completed normally")
}

// Remote agent
response, err := remoteAgent.Ask("What is 2+2?")
```

### AskStream

Stream the response chunk by chunk.

```go
// Local agent
finalResponse, err := agent.AskStream("Explain AI", func(chunk ChatResponse) error {
    fmt.Print(chunk.Text)
    return nil
})
if err != nil {
    log.Fatal(err)
}

// Check final response
if finalResponse.IsFinishReasonStop() {
    fmt.Println("\nCompleted normally")
}

// Remote agent
finalResponse, err := remoteAgent.AskStream("Explain AI", func(chunk ChatResponse) error {
    fmt.Print(chunk.Text)
    // You can also check chunk.FinishReason during streaming
    if chunk.FinishReason != "" {
        fmt.Printf("\nFinish reason: %s\n", chunk.FinishReason)
    }
    return nil
})
```

### AddSystemMessage

Add context to the conversation.

```go
// Local agent
err := agent.AddSystemMessage("The user prefers concise answers.")

// Remote agent (sends HTTP request to the server)
err := remoteAgent.AddSystemMessage("The user prefers concise answers.")
```

### GetCurrentContextSize

Get the total size of the conversation context in characters (includes system instructions and all message content).

```go
// Local agent (returns immediately from memory)
size := agent.GetCurrentContextSize()
fmt.Printf("Context has %d characters\n", size)

// Remote agent (fetches from server via HTTP)
size := remoteAgent.GetCurrentContextSize()
fmt.Printf("Remote context has %d characters\n", size)
```

**Notes:**
- Returns total character count (system instructions + all message content)
- Returns 0 if no messages exist
- For `RemoteAgent`, returns 0 on HTTP errors (doesn't include system instructions)
- Very efficient for local agents

### ReplaceMessagesWith

Replace the entire conversation history with new messages.

```go
// Local agent - replace messages
newMessages := []*ai.Message{
    ai.NewSystemTextMessage("You are a helpful assistant"),
    ai.NewUserTextMessage("Hello"),
    ai.NewModelTextMessage("Hi! How can I help you?"),
}
err := agent.ReplaceMessagesWith(newMessages)
if err != nil {
    log.Fatal(err)
}

// Clear all messages
err = agent.ReplaceMessagesWith([]*ai.Message{})

// Remote agent - NOT SUPPORTED
err := remoteAgent.ReplaceMessagesWith(newMessages)
// Returns error: "ReplaceMessagesWith is not supported for remote agents"
```

**Notes:**
- `ReplaceMessagesWith` is **NOT supported** for `RemoteAgent`
- Remote agents manage history on the server side
- Returns error if `messages` is `nil` (use empty slice to clear)
- Use this to restore saved conversations or implement context pruning

### GetMessages

Get the conversation history.

```go
// Local agent (returns from memory)
messages := agent.GetMessages()
for _, msg := range messages {
    fmt.Printf("%s: %s\n", msg.Role, msg.Content[0].Text)
}

// Remote agent (fetches from server via HTTP)
messages := remoteAgent.GetMessages()
for _, msg := range messages {
    fmt.Printf("%s: %s\n", msg.Role, msg.Content[0].Text)
}
```

### GetInfo

Get agent information.

```go
// Local agent
info, err := agent.GetInfo()
fmt.Printf("Agent: %s, Model: %s\n", info.Name, info.ModelID)
fmt.Printf("Temperature: %.2f, TopP: %.2f\n", info.Config.Temperature, info.Config.TopP)

// Remote agent (fetches from server)
info, err := remoteAgent.GetInfo()
fmt.Printf("Agent: %s, Model: %s\n", info.Name, info.ModelID)
fmt.Printf("Temperature: %.2f, TopP: %.2f\n", info.Config.Temperature, info.Config.TopP)
```

### RAG Agent Methods

These methods are specific to `RagAgent` for semantic search operations.

#### AddTextChunksToStore

Add text chunks (documents) to the RAG store for indexing.

```go
chunks := []snip.TextChunk{
    {
        Content:  "Paris is the capital of France",
        Metadata: map[string]any{"source": "geography.txt", "category": "cities"},
    },
    {
        Content:  "Go is a programming language created by Google",
        Metadata: map[string]any{"source": "tech.txt", "category": "programming"},
    },
}

count, err := ragAgent.AddTextChunksToStore(chunks)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Indexed %d documents\n", count)
```

#### SearchSimilarities

Search for documents similar to a query using semantic search.

```go
similarities, err := ragAgent.SearchSimilarities("What is the capital of France?")
if err != nil {
    log.Fatal(err)
}

fmt.Println("Similar documents:")
for _, doc := range similarities {
    fmt.Println("  -", doc)
}
```

#### GetInfo (RAG Agent)

Get information about the RAG agent, including store details.

```go
info, err := ragAgent.GetInfo()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Agent: %s\n", info.Name)
fmt.Printf("Model: %s\n", info.ModelID)
fmt.Printf("Embedding Dimension: %d\n", info.EmbeddingDimension)
fmt.Printf("Store: %s at %s\n", info.StoreName, info.StorePath)
fmt.Printf("Documents: %d\n", info.NumberOfDocuments)
```

#### GetNumberOfDocuments

Get the count of documents in the store.

```go
count := ragAgent.GetNumberOfDocuments()
fmt.Printf("Store contains %d documents\n", count)
```

#### IsStoreInitialized

Check if the document store is initialized.

```go
if ragAgent.IsStoreInitialized() {
    fmt.Println("Store is ready")
} else {
    fmt.Println("Store not initialized")
}
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
    "github.com/snipwise/snip-sdk/snip"
)

func main() {
    ctx := context.Background()

    // Create an agent with custom configuration
    agent, err := snip.NewAgent(
        ctx,
        snip.AgentConfig{
            Name:               "tutor-bot",
            SystemInstructions: "You are a patient and knowledgeable tutor who explains concepts clearly.",
            ModelID:            "hf.co/menlo/jan-nano-gguf:q4_k_m",
            EngineURL:          "http://localhost:12434/engines/llama.cpp/v1",
        },
        snip.ModelConfig{
            Temperature:      0.8,
            MaxTokens:        1000,
            FrequencyPenalty: 0.5,
        },
        snip.EnableChatFlowWithMemory(),
        snip.EnableChatStreamFlowWithMemory(),
    )
    if err != nil {
        log.Fatal(err)
    }

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

        _, err := agent.AskStream(question, func(chunk ChatResponse) error {
            fmt.Print(chunk.Text)
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

## Context Management Best Practices

### When to use GetCurrentContextSize()

- **Before asking questions**: Check if context is getting too large
- **For monitoring**: Track conversation size in logs/metrics
- **Conditional logic**: Implement different behaviors based on context size

```go
size := agent.GetCurrentContextSize()
if size > 5000 {
    log.Printf("Warning: Large context (%d characters)", size)
}
```

### When to use ReplaceMessagesWith()

- **Context pruning**: Keep only recent messages when hitting limits
- **Conversation restoration**: Load saved conversations from disk/database
- **Reset conversations**: Start fresh while keeping the same agent
- **Testing**: Set up specific conversation states for tests

```go
// Prune old messages when context is too large
if agent.GetCurrentContextSize() > 10000 {
    messages := agent.GetMessages()
    agent.ReplaceMessagesWith(messages[len(messages)-20:]) // Keep last 20
}
```

**Important:** `ReplaceMessagesWith()` is only available for local `Agent`, not `RemoteAgent`.

## Remote Agent Use Cases

The `RemoteAgent` is useful for:

1. **Microservices Architecture**: Connect to AI agents running as separate services
2. **Load Distribution**: Multiple clients can connect to the same agent server
3. **Language Agnostic Clients**: Any HTTP client can interact with the agent server
4. **Testing**: Test your agent's HTTP API programmatically
5. **Agent Orchestration**: Build systems where multiple agents communicate

**Key Benefits:**
- Same interface as local `Agent` (implements `AIAgent`)
- Automatic endpoint URL construction
- Built-in error handling for HTTP communication
- Supports all agent operations (Ask, AskStream, AddSystemMessage, GetMessages, GetInfo)

**Example: Multiple Clients to One Server**

```go
// Start one server
agent, err := snip.NewAgent(ctx, agentConfig, modelConfig,
    snip.EnableChatFlowWithMemory(),
    snip.EnableServer(snip.ConfigHTTP{Address: "0.0.0.0:9100"}),
)
if err != nil {
    log.Fatal(err)
}
go agent.Serve()

// Connect multiple clients
client1 := snip.NewRemoteAgent("Client 1", snip.ConfigHTTP{Address: "0.0.0.0:9100"})
client2 := snip.NewRemoteAgent("Client 2", snip.ConfigHTTP{Address: "0.0.0.0:9100"})

// Both clients share the same conversation history on the server
response1, _ := client1.Ask("My name is Alice")
response2, _ := client2.Ask("What is my name?")  // Will know it's Alice
fmt.Println(response2.Text)
```

## RAG Agent Use Cases

The `RagAgent` is useful for:

1. **Document Search**: Build semantic search engines for your documents
2. **Knowledge Bases**: Create question-answering systems over your own data
3. **Context Retrieval**: Find relevant information to augment LLM prompts
4. **Similarity Detection**: Identify similar content across documents
5. **Content Recommendation**: Suggest related content based on semantic similarity

**Key Features:**
- Automatic embedding generation using specified model
- Persistent document storage (localvec)
- Metadata support for document organization
- Semantic search (not just keyword matching)
- Configurable storage path and name

**Common Embedding Models:**
- `ai/mxbai-embed-large` - Good general-purpose embedding model
- `ai/granite-embedding-multilingual` - Multilingual support
- `ai/embeddinggemma` - Google's Gemma embedding model

**Example: Building a Knowledge Base**

```go
// Create RAG agent with your preferred embedding model
ragAgent, err := snip.NewRagAgent(
    ctx,
    snip.RagAgentConfig{
        Name:      "knowledge-base",
        ModelID:   "ai/mxbai-embed-large",
        EngineURL: "http://localhost:12434/engines/llama.cpp/v1",
    },
    snip.StoreConfig{
        StoreName: "company-docs",
        StorePath: "./knowledge-base",
    },
)
if err != nil {
    log.Fatal(err)
}

// Index your documents with metadata
docs := []snip.TextChunk{
    {
        Content:  "Our customer support is available 24/7 via email and chat",
        Metadata: map[string]any{"category": "support", "doc_id": "cs-001"},
    },
    {
        Content:  "Standard shipping takes 5-7 business days",
        Metadata: map[string]any{"category": "shipping", "doc_id": "ship-001"},
    },
}

_, err = ragAgent.AddTextChunksToStore(docs)
if err != nil {
    log.Fatal(err)
}

// Search for relevant information
results, err := ragAgent.SearchSimilarities("How can I contact support?")
if err != nil {
    log.Fatal(err)
}

// Use results to augment your LLM prompt
for _, result := range results {
    fmt.Println("Relevant info:", result)
}
```

**Best Practices:**
- Choose an embedding model that matches your content language and domain
- Include meaningful metadata to help organize and filter results
- Use consistent chunking strategies for better search results
- Store documents persistently by using a dedicated storage path
- Check `GetNumberOfDocuments()` to avoid re-indexing existing documents

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
response, err := agent.Ask("Hello")
if err != nil {
    log.Printf("Error: %v", err)
    // Handle error appropriately
    return
}

// Check completion status
if !response.IsFinishReasonStop() {
    log.Printf("Warning: Response did not complete normally. Reason: %s", response.FinishReason)
    if response.FinishMessage != "" {
        log.Printf("Message: %s", response.FinishMessage)
    }
}

fmt.Println(response.Text)
```

## Installation

```bash
go get github.com/snipwise/snip-sdk/snip
```

## Helper Functions

The package provides utility functions for working with models:

### GetModelsList

```go
func GetModelsList(ctx context.Context, modelRunnerEndpoint string) ([]string, error)
```

Retrieves a list of all available models from the inference engine.

### IsModelAvailable

```go
func IsModelAvailable(ctx context.Context, modelRunnerEndpoint, modelID string) bool
```

Checks if a specific model is available at the inference engine endpoint.

**Note:** This function is automatically called by `NewAgent` and `NewRagAgent` during initialization.

## Dependencies

This package uses:
- [Firebase Genkit](https://github.com/firebase/genkit) - AI framework and orchestration
- [OpenAI Go SDK](https://github.com/openai/openai-go) - OpenAI API client for model interaction
- [localvec](https://github.com/firebase/genkit/go/plugins/localvec) - Local vector store for RAG operations

## License

See the main repository for license information.
