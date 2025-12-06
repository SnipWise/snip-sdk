package chat

import "net/http"

const (
	// DefaultChatFlowPath is the default endpoint path for the standard chat flow
	DefaultChatFlowPath = "/api/chat"

	// DefaultChatStreamFlowPath is the default endpoint path for the streaming chat flow
	DefaultChatStreamFlowPath = "/api/chat-stream"

	// DefaultInformationPath is the default endpoint path for agent information
	DefaultInformationPath = "/api/information"

	// DefaultShutdownPath is the default endpoint path for shutting down the server
	// Set to "-" by default to disable the shutdown endpoint
	DefaultShutdownPath = "-"

	// DefaultCancelStreamPath is the default endpoint path for canceling a streaming completion
	DefaultCancelStreamPath = "/api/cancel-stream-completion"

	// DefaultAddSystemMessagePath is the default endpoint path for adding context to messages
	DefaultAddSystemMessagePath = "/api/add-system-message"

	// DefaultHealthcheckPath is the default endpoint path for healthcheck
	DefaultHealthcheckPath = "/healthcheck"

	// DefaultGetMessagesPath is the default endpoint path for retrieving conversation messages
	DefaultGetMessagesPath = "/api/messages"
)

// ConfigHTTP holds the HTTP server configuration for exposing agent flows
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

	// CancelStreamPath is the endpoint path for canceling a streaming completion
	// If empty, defaults to DefaultCancelStreamPath ("/api/cancel-stream-completion")
	CancelStreamPath string

	// AddContextPath is the endpoint path for adding context to messages
	// If empty, defaults to DefaultAddContextPath ("/api/add-context")
	AddContextPath string

	// HealthcheckPath is the endpoint path for healthcheck
	// If empty, defaults to DefaultHealthcheckPath ("/healthcheck")
	HealthcheckPath string

	// GetMessagesPath is the endpoint path for retrieving conversation messages
	// If empty, defaults to DefaultGetMessagesPath ("/api/messages")
	GetMessagesPath string

	// ChatFlowHandler is the HTTP handler for the standard chat flow endpoint
	// If nil, will be auto-configured from the agent's chatFlow
	ChatFlowHandler http.HandlerFunc

	// ChatStreamFlowHandler is the HTTP handler for the streaming chat flow endpoint
	// If nil, will be auto-configured from the agent's chatStreamFlow
	ChatStreamFlowHandler http.HandlerFunc
}
