# Two HTTP Servers Sample

This sample demonstrates how to run **two AI agents simultaneously** with their own HTTP endpoints using goroutines and the `EnableServer` option.

## Overview

This sample runs two independent agents in parallel:

### Agent One (HTTP_Agent_1)
- HTTP server running on `0.0.0.0:9100`
- Standard chat flow (non-streaming) at `/api/chat`
- Streaming chat flow at `/api/chat-stream`
- Agent information endpoint at `/api/information`
- Shutdown endpoint at `/server/shutdown`

### Agent Two (HTTP_Agent_2)
- HTTP server running on `0.0.0.0:9200`
- Standard chat flow (non-streaming) at `/api/chat`
- Streaming chat flow at `/api/chat-stream`
- Agent information endpoint at `/api/information`
- Shutdown endpoint at `/server/shutdown`

Both agents run concurrently using Go goroutines, allowing them to handle requests independently.

## Running the Sample

```bash
cd samples/05-two-http-servers
go run main.go
```

The servers will start and display:
```
Agent 'HTTP_Agent_1' created successfully
Available endpoints for Agent One:
  POST http://0.0.0.0:9100/api/chat
  POST http://0.0.0.0:9100/api/chat-stream
  GET  http://0.0.0.0:9100/api/information
  POST http://0.0.0.0:9100/server/shutdown

Agent 'HTTP_Agent_2' created successfully
Available endpoints for Agent Two:
  POST http://0.0.0.0:9200/api/chat
  POST http://0.0.0.0:9200/api/chat-stream
  GET  http://0.0.0.0:9200/api/information
  POST http://0.0.0.0:9200/server/shutdown

Starting Agent One server on port 9100...
Starting Agent Two server on port 9200...
```

## Available Endpoints

Both agents expose the same endpoints on different ports.

### 1. POST /api/chat
Standard chat completion (non-streaming)

**Request format:**
```json
{
  "data": {
    "message": "Your question here"
  }
}
```

**Using curl for Agent One (port 9100):**
```bash
curl -X POST http://0.0.0.0:9100/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "data": {
      "message": "Hello, who are you?"
    }
  }'
```

**Using curl for Agent Two (port 9200):**
```bash
curl -X POST http://0.0.0.0:9200/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "data": {
      "message": "Hello, who are you?"
    }
  }'
```

**Using the provided script (defaults to Agent One):**
```bash
./01-call-chat.sh
```

**To test Agent Two, override the SERVICE_URL:**
```bash
SERVICE_URL="http://0.0.0.0:9200/api/chat" ./01-call-chat.sh
```

**Response:**
```json
{
  "response": "I am a helpful assistant..."
}
```

### 2. POST /api/chat-stream
Streaming chat completion (Server-Sent Events)

**Request format:**
```json
{
  "data": {
    "message": "Your question here"
  }
}
```

**Using curl for Agent One (port 9100):**
```bash
curl --no-buffer http://0.0.0.0:9100/api/chat-stream \
  -H "Content-Type: application/json" \
  -H "Accept: text/event-stream" \
  -d '{
    "data": {
      "message": "Tell me a story about your adventures"
    }
  }'
```

**Using curl for Agent Two (port 9200):**
```bash
curl --no-buffer http://0.0.0.0:9200/api/chat-stream \
  -H "Content-Type: application/json" \
  -H "Accept: text/event-stream" \
  -d '{
    "data": {
      "message": "Tell me a story about your adventures"
    }
  }'
```

**Using the provided script (defaults to Agent One):**
```bash
./02-call-chat-stream.sh
```

**To test Agent Two, override the SERVICE_URL:**
```bash
SERVICE_URL="http://0.0.0.0:9200/api/chat-stream" ./02-call-chat-stream.sh
```

**Response:**
The response will be streamed token by token using Server-Sent Events (SSE) format:
```
data: {"message":"Once"}
data: {"message":" upon"}
data: {"message":" a"}
data: {"message":" time"}
...
```

### 3. GET /api/information
Get agent information (name, model, configuration)

**Using curl for Agent One (port 9100):**
```bash
curl -X GET http://0.0.0.0:9100/api/information
```

**Using curl for Agent Two (port 9200):**
```bash
curl -X GET http://0.0.0.0:9200/api/information
```

**Response:**
```json
{
  "name": "HTTP_Agent_1",
  "model_id": "hf.co/menlo/jan-nano-gguf:q4_k_m",
  "config": {
    "temperature": 0.5,
    "top_p": 0.9
  }
}
```

### 4. POST /server/shutdown
Gracefully shutdown a specific server

**Using curl to shutdown Agent One:**
```bash
curl -X POST http://0.0.0.0:9100/server/shutdown \
  -H "Content-Type: application/json"
```

**Using curl to shutdown Agent Two:**
```bash
curl -X POST http://0.0.0.0:9200/server/shutdown \
  -H "Content-Type: application/json"
```

**Using the provided scripts:**
```bash
# Shutdown Agent One (port 9100)
./03-stop-server-1.sh

# Shutdown Agent Two (port 9200)
./04-stop-server-2.sh
```

**Response:**
```json
{
  "status": "shutting down"
}
```

Each server will perform a graceful shutdown with a 5-second timeout independently.

## Stopping the Servers

There are **three ways** to stop each server:

1. **Press Ctrl+C** - Sends SIGINT signal for graceful shutdown of **both servers**
2. **HTTP POST to `/server/shutdown`** - Trigger shutdown of a **specific agent** via API endpoint:
   - `./03-stop-server-1.sh` - Stop Agent One only
   - `./04-stop-server-2.sh` - Stop Agent Two only
3. **Programmatically** - Call `agent.Stop()` method in code for a specific agent

Note: When using Ctrl+C, both servers will shut down gracefully. To stop only one server, use the HTTP shutdown endpoint for that specific agent.

## Configuration

Each server is configured independently using the `EnableServer` option with different ports:

**Agent One:**
```go
smart.EnableServer(smart.ConfigHTTP{
    Address:            "0.0.0.0:9100",
    ChatFlowPath:       "/api/chat",
    ChatStreamFlowPath: "/api/chat-stream",
    InformationPath:    "/api/information",
    ShutdownPath:       "/server/shutdown",
})
```

**Agent Two:**
```go
smart.EnableServer(smart.ConfigHTTP{
    Address:            "0.0.0.0:9200",
    ChatFlowPath:       "/api/chat",
    ChatStreamFlowPath: "/api/chat-stream",
    InformationPath:    "/api/information",
    ShutdownPath:       "/server/shutdown",
})
```

### ConfigHTTP Structure

```go
type ConfigHTTP struct {
    // Address is the network address to bind to (e.g., "0.0.0.0:9100", ":8080")
    Address string

    // ChatFlowPath is the endpoint path for the standard chat flow
    // If empty, defaults to DefaultChatFlowPath ("/api/chat")
    ChatFlowPath string

    // ChatStreamFlowPath is the endpoint path for the streaming chat flow
    // If empty, defaults to DefaultChatStreamFlowPath ("/api/chat-stream")
    ChatStreamFlowPath string

    // InformationPath is the endpoint path for agent information
    // If empty, defaults to DefaultInformationPath ("/api/information")
    InformationPath string

    // ShutdownPath is the endpoint path for shutting down the server
    // If empty, defaults to DefaultShutdownPath ("-" - disabled)
    // Set to "-" to disable the shutdown endpoint, or provide a custom path like "/server/shutdown"
    ShutdownPath string

    // ChatFlowHandler is the HTTP handler for the standard chat flow endpoint
    // If nil, will be auto-configured from the agent's chatFlow
    ChatFlowHandler http.HandlerFunc

    // ChatStreamFlowHandler is the HTTP handler for the streaming chat flow endpoint
    // If nil, will be auto-configured from the agent's chatStreamFlow
    ChatStreamFlowHandler http.HandlerFunc
}
```

### Configuration Examples

**Using default paths:**
```go
smart.EnableServer(smart.ConfigHTTP{
    Address: "0.0.0.0:9100",
    // All paths will use defaults:
    // ChatFlowPath: "/api/chat"
    // ChatStreamFlowPath: "/api/chat-stream"
    // InformationPath: "/api/information"
    // ShutdownPath: "-" (disabled)
})
```

**Custom paths:**
```go
smart.EnableServer(smart.ConfigHTTP{
    Address:            "0.0.0.0:9100",
    ChatFlowPath:       "/api/chat",
    ChatStreamFlowPath: "/api/chat-stream",
    InformationPath:    "/api/info",
    ShutdownPath:       "/api/stop",
})
```

**Disable shutdown endpoint:**
```go
smart.EnableServer(smart.ConfigHTTP{
    Address:      "0.0.0.0:9100",
    ShutdownPath: "-", // Disables the shutdown endpoint
})
```

## Test Scripts

The sample includes four bash scripts for testing both agents:

1. **`01-call-chat.sh`** - Test the standard chat endpoint (defaults to Agent One on port 9100)
2. **`02-call-chat-stream.sh`** - Test the streaming chat endpoint with SSE parsing (defaults to Agent One on port 9100)
3. **`03-stop-server-1.sh`** - Gracefully shutdown Agent One via HTTP
4. **`04-stop-server-2.sh`** - Gracefully shutdown Agent Two via HTTP

All scripts support the `SERVICE_URL` environment variable:

```bash
# Test Agent One (default)
./01-call-chat.sh
./02-call-chat-stream.sh

# Test Agent Two by overriding SERVICE_URL
SERVICE_URL="http://0.0.0.0:9200/api/chat" ./01-call-chat.sh
SERVICE_URL="http://0.0.0.0:9200/api/chat-stream" ./02-call-chat-stream.sh

# Stop servers individually
./03-stop-server-1.sh  # Stop Agent One
./04-stop-server-2.sh  # Stop Agent Two
```

## Custom Handlers (Optional)

You can provide custom HTTP handlers if needed:

```go
smart.EnableServer(smart.ConfigHTTP{
    Address: "0.0.0.0:9100",
    ChatFlowHandler: func(w http.ResponseWriter, r *http.Request) {
        // Custom handler logic
        // Your implementation here
    },
})
```

If handlers are not provided, they will be automatically configured using `genkit.Handler()`.

## Starting Multiple Servers

To run multiple servers concurrently, use Go goroutines. The first server runs in a goroutine, while the second runs in the main thread:

```go
// Start agentOne in a goroutine (non-blocking)
go func() {
    log.Println("Starting Agent One server on port 9100...")
    if err := agentOne.Serve(); err != nil {
        log.Fatalf("Server One error: %v", err)
    }
}()

// Start agentTwo in the main thread (blocking, handles Ctrl+C)
log.Println("Starting Agent Two server on port 9200...")
if err := agentTwo.Serve(); err != nil {
    log.Fatalf("Server Two error: %v", err)
}
```

**Important notes:**
- `Serve()` is a **blocking call** that runs until the server stops or an error occurs
- The first server runs in a **goroutine** for concurrent execution
- The second server runs in the **main thread** to handle graceful shutdown signals
- Both servers automatically handle **graceful shutdown** on SIGINT/SIGTERM signals (Ctrl+C)
- A **5-second timeout** is used for graceful shutdown to allow in-flight requests to complete
- Each server can be stopped independently via its HTTP shutdown endpoint

## Environment Variables

The sample uses the following environment variables:

- `MODEL_RUNNER_BASE_URL` - URL of the model runner (default: `http://localhost:12434/engines/llama.cpp/v1`)
- `CHAT_MODEL` - Model to use for chat (default: `hf.co/menlo/jan-nano-gguf:q4_k_m`)
- `SERVICE_URL` - Used by test scripts to override endpoint URLs

## Example Usage Flow

1. **Start both servers:**
   ```bash
   go run main.go
   ```
   Both Agent One (port 9100) and Agent Two (port 9200) will start simultaneously.

2. **In another terminal, test Agent One:**
   ```bash
   # Test standard chat
   ./01-call-chat.sh

   # Test streaming chat
   ./02-call-chat-stream.sh
   ```

3. **Test Agent Two:**
   ```bash
   # Test standard chat on Agent Two
   SERVICE_URL="http://0.0.0.0:9200/api/chat" ./01-call-chat.sh

   # Test streaming chat on Agent Two
   SERVICE_URL="http://0.0.0.0:9200/api/chat-stream" ./02-call-chat-stream.sh
   ```

4. **Stop servers:**
   ```bash
   # Option 1: Stop both servers at once
   # Press Ctrl+C in the server terminal

   # Option 2: Stop servers individually
   ./03-stop-server-1.sh  # Stops Agent One only
   ./04-stop-server-2.sh  # Stops Agent Two only
   ```

## Graceful Shutdown

When a server receives a shutdown signal (via Ctrl+C, HTTP endpoint, or `Stop()` method), it:

1. Stops accepting new connections
2. Waits up to 5 seconds for active requests to complete
3. Logs the shutdown process
4. Exits cleanly

**Shutdown Behavior:**
- **Ctrl+C**: Both servers receive the SIGINT signal and shut down gracefully
- **HTTP endpoint**: Only the targeted server shuts down (Agent One or Agent Two)
- **Independent operation**: Each server manages its own lifecycle

Example output when stopping via HTTP endpoint:
```
Shutdown requested via HTTP endpoint
Shutting down server gracefully...
Server stopped
```

## Key Differences from Single Server Sample

This sample demonstrates:
1. **Concurrent execution** using goroutines to run multiple servers
2. **Independent agents** with separate ports and configurations
3. **Isolated shutdown** - each server can be stopped independently
4. **Parallel request handling** - both agents can serve requests simultaneously
