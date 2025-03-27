package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/metoro-io/mcp-golang/transport"
)

type SimpleHandlerFunc func(res http.ResponseWriter, req *http.Request)

type SimpleHandlerWithContextFunc func(ctx context.Context, res http.ResponseWriter, req *http.Request)

// SimpleTransport implements a stateless HTTP transport for MCP using Simple
type SimpleTransport struct {
	*baseTransport
}

// NewSimpleTransport creates a new Simple transport
func NewSimpleTransport() *SimpleTransport {
	return &SimpleTransport{
		baseTransport: newBaseTransport(),
	}
}

// Start implements Transport.Start - no-op for Simple transport as it's handled by Simple
func (t *SimpleTransport) Start(ctx context.Context) error {
	return nil
}

// Send implements Transport.Send
func (t *SimpleTransport) Send(ctx context.Context, message *transport.BaseJsonRpcMessage) error {
	key := message.JsonRpcResponse.Id
	responseChannel := t.responseMap[int64(key)]
	if responseChannel == nil {
		return fmt.Errorf("no response channel found for key: %d", key)
	}
	responseChannel <- message
	return nil
}

// Close implements Transport.Close
func (t *SimpleTransport) Close() error {
	if t.closeHandler != nil {
		t.closeHandler()
	}
	return nil
}

// SetCloseHandler implements Transport.SetCloseHandler
func (t *SimpleTransport) SetCloseHandler(handler func()) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.closeHandler = handler
}

// SetErrorHandler implements Transport.SetErrorHandler
func (t *SimpleTransport) SetErrorHandler(handler func(error)) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.errorHandler = handler
}

// SetMessageHandler implements Transport.SetMessageHandler
func (t *SimpleTransport) SetMessageHandler(handler func(ctx context.Context, message *transport.BaseJsonRpcMessage)) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.messageHandler = handler
}

// Handler returns a Simple handler function that can be used with Simple's router
func (t *SimpleTransport) HandlerWithContext() SimpleHandlerWithContextFunc {
	return func(ctx context.Context, res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			res.WriteHeader(http.StatusMethodNotAllowed)
			res.Write([]byte("Only POST method is supported"))
			return
		}

		body, err := t.readBody(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte(err.Error()))
			return
		}

		response, err := t.handleMessage(ctx, body)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(err.Error()))
			return
		}

		jsonData, err := json.Marshal(response)
		if err != nil {
			if t.errorHandler != nil {
				t.errorHandler(fmt.Errorf("failed to marshal response: %w", err))
			}

			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte("Failed to marshal response"))
			return
		}

		res.WriteHeader(http.StatusOK)
		res.Write(jsonData)
	}
}

func (t *SimpleTransport) Handler() SimpleHandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx := context.Background()
		if req.Method != http.MethodPost {
			res.WriteHeader(http.StatusMethodNotAllowed)
			res.Write([]byte("Only POST method is supported"))
			return
		}

		body, err := t.readBody(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte(err.Error()))
			return
		}

		response, err := t.handleMessage(ctx, body)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(err.Error()))
			return
		}

		jsonData, err := json.Marshal(response)
		if err != nil {
			if t.errorHandler != nil {
				t.errorHandler(fmt.Errorf("failed to marshal response: %w", err))
			}

			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte("Failed to marshal response"))
			return
		}

		res.WriteHeader(http.StatusOK)
		res.Write(jsonData)
	}
}
