# HTTP Server Sample

This sample demonstrates how to expose an AI agent's chat flows via HTTP endpoints using the `EnableServer` option.

## Overview

The agent is configured with:
- Standard chat flow (non-streaming) at `/api/chat`
- Streaming chat flow at `/api/chat-stream`
- Agent information endpoint at `/api/information`
- Shutdown endpoint at `/server/shutdown`
- HTTP server running on `0.0.0.0:9100`

## Running the Sample

```bash
cd samples/04-http-server
go run main.go
```

The server will start and display:
```
Agent 'HTTP Agent' created successfully
Available endpoints:
  POST http://0.0.0.0:9100/api/chat
  POST http://0.0.0.0:9100/api/chat-stream
  GET  http://0.0.0.0:9100/api/information
  POST http://0.0.0.0:9100/shutdown (stop the server)

Registered endpoint: POST /api/chat
Registered endpoint: POST /api/chat-stream
Registered endpoint: GET /api/information
Registered endpoint: POST /server/shutdown
Starting HTTP server on 0.0.0.0:9100 (Press Ctrl+C to stop)
```

## Available Endpoints

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

**Using curl:**
```bash
curl -X POST http://0.0.0.0:9100/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "data": {
      "message": "What is the capital of France?"
    }
  }'
```

**Using the provided script:**
```bash
./01-call-chat.sh
```

**Response:**
```json
{
  "response": "The capital of France is Paris."
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

**Using curl:**
```bash
curl --no-buffer http://0.0.0.0:9100/api/chat-stream \
  -H "Content-Type: application/json" \
  -H "Accept: text/event-stream" \
  -d '{
    "data": {
      "message": "Tell me a short story"
    }
  }'
```

**Using the provided script:**
```bash
./02-call-chat-stream.sh
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

**Using curl:**
```bash
curl -X GET http://0.0.0.0:9100/api/information
```

**Response:**
```json
{
  "name": "HTTP_Agent",
  "model_id": "hf.co/menlo/jan-nano-gguf:q4_k_m",
  "config": {
    "temperature": 0.5,
    "top_p": 0.9
  }
}
```

### 4. POST /server/shutdown
Gracefully shutdown the server

**Using curl:**
```bash
curl -X POST http://0.0.0.0:9100/server/shutdown \
  -H "Content-Type: application/json"
```

**Using the provided script:**
```bash
./03-stop-server.sh
```

**Response:**
```json
{
  "status": "shutting down"
}
```

The server will perform a graceful shutdown with a 5-second timeout.

## Stopping the Server

There are **three ways** to stop the server:

1. **Press Ctrl+C** - Sends SIGINT signal for graceful shutdown
2. **HTTP POST to `/server/shutdown`** - Trigger shutdown via API endpoint
3. **Programmatically** - Call `agent.Stop()` method in code

## Configuration

The server is configured using the `EnableServer` option:

```go
smart.EnableServer(smart.ConfigHTTP{
    Address:            "0.0.0.0:9100",
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

The sample includes three bash scripts for testing:

1. **`01-call-chat.sh`** - Test the standard chat endpoint
2. **`02-call-chat-stream.sh`** - Test the streaming chat endpoint with SSE parsing
3. **`03-stop-server.sh`** - Gracefully shutdown the server via HTTP

All scripts support the `SERVICE_URL` environment variable:

```bash
# Use default URLs
./01-call-chat.sh

# Override the service URL
SERVICE_URL="http://localhost:8080/custom/path" ./01-call-chat.sh
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

## Starting the Server

Call the `Serve()` method on the agent to start the HTTP server:

```go
if err := agent.Serve(); err != nil {
    log.Fatalf("Server error: %v", err)
}
```

**Important notes:**
- `Serve()` is a **blocking call** that runs until the server stops or an error occurs
- The server automatically handles **graceful shutdown** on SIGINT/SIGTERM signals
- A **5-second timeout** is used for graceful shutdown to allow in-flight requests to complete

## Environment Variables

The sample uses the following environment variables:

- `MODEL_RUNNER_BASE_URL` - URL of the model runner (default: `http://localhost:12434/engines/llama.cpp/v1`)
- `CHAT_MODEL` - Model to use for chat (default: `hf.co/menlo/jan-nano-gguf:q4_k_m`)
- `SERVICE_URL` - Used by test scripts to override endpoint URLs

## Example Usage Flow

1. **Start the server:**
   ```bash
   go run main.go
   ```

2. **In another terminal, test the chat endpoint:**
   ```bash
   ./01-call-chat.sh
   ```

3. **Test the streaming endpoint:**
   ```bash
   ./02-call-chat-stream.sh
   ```

4. **Stop the server:**
   ```bash
   ./03-stop-server.sh
   # Or press Ctrl+C in the server terminal
   ```

## Graceful Shutdown

When the server receives a shutdown signal (via Ctrl+C, HTTP endpoint, or `Stop()` method), it:

1. Stops accepting new connections
2. Waits up to 5 seconds for active requests to complete
3. Logs the shutdown process
4. Exits cleanly

Example output:
```
Shutdown requested via HTTP endpoint
Shutting down server gracefully...
Server stopped
```
