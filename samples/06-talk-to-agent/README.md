# 06 - Talk to Agent (Remote Agent Communication)

This example demonstrates how to create an agent that runs both as an HTTP server and as a client that connects to itself using `RemoteAgent` to communicate.

## Overview

This is an all-in-one example that shows:
1. **Server Agent**: An HTTP server agent that exposes chat endpoints
2. **Remote Agent Client**: A client that connects to the server and asks questions
3. **All in one process**: Both server and client run together in the same program

## Features Demonstrated

- Creating an agent with HTTP server enabled
- Using `RemoteAgent` to connect to a remote agent
- Getting remote agent information with `GetInfo()`
- Non-streaming questions with `Ask()`
- Streaming questions with `AskStream()`
- Memory persistence across multiple questions
- Running server and client in the same process using goroutines

## How to Run

### Prerequisites

Make sure you have a model runner running (e.g., Jan, LM Studio, or Ollama).

### Run the Example

```bash
cd samples/06-talk-to-agent
go run main.go
```

The program will:
1. Start an HTTP server on port 9100 in a goroutine
2. Create a remote agent client that connects to this server
3. Execute four examples:
   - **Example 0**: Get remote agent information
   - **Example 1**: Non-streaming question about Go programming language
   - **Example 2**: Streaming question about goroutines
   - **Example 3**: Follow-up question to test conversation memory

## Environment Variables

You can customize the model and runner:

```bash
export MODEL_RUNNER_BASE_URL="http://localhost:12434/engines/llama.cpp/v1"
export CHAT_MODEL="hf.co/menlo/jan-nano-gguf:q4_k_m"
```

Default values:
- Model Runner: `http://localhost:12434/engines/llama.cpp/v1`
- Chat Model: `hf.co/menlo/jan-nano-gguf:q4_k_m`

## Server Endpoints

The server exposes the following endpoints on `http://0.0.0.0:9100`:

- `GET /api/information` - Get agent information
- `POST /api/chat` - Non-streaming chat
- `POST /api/chat-stream` - Streaming chat
- `POST /server/shutdown` - Shutdown the server

### Manual API Testing

**Get agent information:**
```bash
curl -X GET http://0.0.0.0:9100/api/information
```

**Non-streaming chat:**
```bash
curl -X POST http://0.0.0.0:9100/api/chat \
  -H "Content-Type: application/json" \
  -d '{"data": {"message": "What is Go?"}}'
```

**Streaming chat:**
```bash
curl -X POST http://0.0.0.0:9100/api/chat-stream \
  -H "Content-Type: application/json" \
  -H "Accept: text/event-stream" \
  -d '{"data": {"message": "Explain goroutines"}}'
```

## Code Highlights

### Creating the Server Agent

```go
agentOne, err := snip.NewAgent(
    ctx,
    snip.AgentConfig{
        Name:               "HTTP_Agent_1",
        SystemInstructions: "You are a helpful assistant.",
        ModelID:            chatModelId,
        EngineURL:          engineURL,
    },
    snip.ModelConfig{
        Temperature: 0.5,
        TopP:        0.9,
    },
    snip.EnableChatFlowWithMemory(),
    snip.EnableChatStreamFlowWithMemory(),
    snip.EnableServer(snip.ConfigHTTP{
        Address:            "0.0.0.0:9100",
        ChatFlowPath:       "/api/chat",
        ChatStreamFlowPath: "/api/chat-stream",
        InformationPath:    "/api/information",
        ShutdownPath:       "/server/shutdown",
    }),
)
if err != nil {
    log.Fatalf("Failed to create agent: %v", err)
}
```

### Starting the Server in a Goroutine

```go
go func() {
    log.Println("Starting Agent One server on port 9100...")
    if err := agentOne.Serve(); err != nil {
        log.Fatalf("Server One error: %v", err)
    }
}()
```

### Creating the Remote Agent Client

```go
remoteAgent := snip.NewRemoteAgent(
    "Remote Knowledge Agent",
    snip.ConfigHTTP{
        Address:            "0.0.0.0:9100",
        ChatFlowPath:       "/api/chat",
        ChatStreamFlowPath: "/api/chat-stream",
        InformationPath:    "/api/information",
    },
)
```

### Getting Agent Information

```go
info, err := remoteAgent.GetInfo()
if err != nil {
    log.Fatalf("❌ Error getting agent info: %v", err)
}

log.Printf("Agent Name: %s", info.Name)
log.Printf("Model ID: %s", info.ModelID)
log.Printf("Temperature: %.2f", info.Config.Temperature)
log.Printf("TopP: %.2f", info.Config.TopP)
```

### Non-streaming Communication

```go
response, err := remoteAgent.AskWithMemory("What is Go programming language?")
if err != nil {
    log.Fatalf("❌ Error asking question: %v", err)
}
fmt.Println(response)
```

### Streaming Communication

```go
fullResponse, err := remoteAgent.AskStreamWithMemory(
    "Explain what are goroutines in 2 sentences",
    func(chunk snip.ChatResponse) error {
        fmt.Print(chunk.Text)
        return nil
    },
)
if err != nil {
    log.Fatalf("❌ Error asking streaming question: %v", err)
}
```

### Testing Memory with Follow-up Questions

```go
response3, err := remoteAgent.AskWithMemory("What are the benefits of using it?")
if err != nil {
    log.Fatalf("❌ Error asking follow-up question: %v", err)
}
fmt.Println(response3)
```

## Expected Output

When you run the example, you should see output similar to:

```
2025/12/02 08:34:34 Starting Agent One server on port 9100...
2025/12/02 08:34:34 ⏳ Waiting for the server to start and the model to load...
2025/12/02 08:34:34 📝 Example 0: Get remote agent information
2025/12/02 08:34:34 ---
2025/12/02 08:34:34 [AGENT_ONE] Registered endpoint: GET /api/information
2025/12/02 08:34:34 [AGENT_ONE] Registered endpoint: POST /api/chat
2025/12/02 08:34:34 [AGENT_ONE] Registered endpoint: POST /api/chat-stream
2025/12/02 08:34:34 [AGENT_ONE] Registered endpoint: POST /server/shutdown
2025/12/02 08:34:34 [AGENT_ONE] Starting HTTP server on 0.0.0.0:9100
2025/12/02 08:34:34 Agent Name: AGENT_ONE
2025/12/02 08:34:34 Model ID: ai/qwen2.5:0.5B-F16
2025/12/02 08:34:34 Temperature: 0.50
2025/12/02 08:34:34 TopP: 0.90
2025/12/02 08:34:34
2025/12/02 08:34:34 ✅ Agent information retrieved
2025/12/02 08:34:34
2025/12/02 08:34:34 📝 Example 1: Non-streaming question
2025/12/02 08:34:34 Question: What is Go programming language?
2025/12/02 08:34:34 ---

[... Go programming language explanation ...]

2025/12/02 08:34:50 ✅ Non-streaming request completed

2025/12/02 08:34:50 📝 Example 2: Streaming question
2025/12/02 08:34:50 Question: Explain what are goroutines in 2 sentences
2025/12/02 08:34:50 ---
Goroutines are lightweight threads...

2025/12/02 08:34:53 ✅ Streaming request completed
2025/12/02 08:34:53 📊 Total response length: 259 characters

2025/12/02 08:34:53 📝 Example 3: Follow-up question (testing memory)
2025/12/02 08:34:53 Question: What are the benefits of using it?
2025/12/02 08:34:53 ---
The benefits of using Go include...

2025/12/02 08:34:55 ✅ Follow-up question completed
```

## Notes

- The server and client run in the same process for demonstration purposes
- The server keeps conversation history in memory
- Each remote agent client maintains its own session with the server
- The server can be stopped with Ctrl+C
- In production, you would typically run the server and client as separate processes
- The third question ("What are the benefits of using it?") tests memory by referring back to previous context
