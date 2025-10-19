package handler

import (
	"fmt"
	"log/slog"
	"net"
	"sync"

	"github.com/k1ender/go-stash/internal/store"
)

type Response interface {
	Serialize() ([]byte, error)
}

type CommandHandler interface {
	Handle(cmd []byte) (Response, error)
}

type Handler struct {
	handlers map[Command]CommandHandler
}

func NewHandler(store store.Store) *Handler {
	handlers := make(map[Command]CommandHandler)

	getHandler := NewGetHandler(store)
	handlers[GetCommand] = getHandler

	setHandler := NewSetHandler(store)
	handlers[SetCommand] = setHandler

	incrHandler := NewIncrHandler(store)
	handlers[IncrCommand] = incrHandler

	decrHandler := NewDecrHandler(store)
	handlers[DecrCommand] = decrHandler

	delHandler := NewDelHandler(store)
	handlers[DelCommand] = delHandler

	return &Handler{
		handlers: handlers,
	}
}

var bufPool = sync.Pool{
	New: func() any {
		return make([]byte, 1024)
	},
}

// Handle processes a client connection by reading a command, dispatching it to the appropriate handler,
// and sending the response back to the client.
//
// Parameters:
//   - client: The client connection to handle.
//
// Returns:
//   - A boolean indicating whether the connection should be closed.
//   - An error if any occurred during processing.
func (h *Handler) Handle(client net.Conn) (bool, error) {
	buf := bufPool.Get().([]byte)
	defer bufPool.Put(buf)

	n, err := client.Read(buf)
	cmd := buf[:n]
	if err != nil {
		slog.Debug("failed to read command from client", slog.Any("error", err))
		h.fail(client)
		return true, fmt.Errorf("failed to read command from client: %w", err)
	}

	command := Command(cmd[:3])
	handler, ok := h.handlers[command]
	if !ok {
		slog.Debug("Received unknown command", slog.Any("command", command))
		h.fail(client)
		return false, fmt.Errorf("unknown command: %s", command)
	}
	response, err := handler.Handle(cmd)
	if err != nil {
		h.fail(client)
		return false, fmt.Errorf("failed to handle GET command: %w", err)
	}

	data, err := response.Serialize()
	if err != nil {
		h.fail(client)
		return false, fmt.Errorf("failed to serialize response: %w", err)
	}

	_, err = client.Write(data)
	if err != nil {
		return true, fmt.Errorf("failed to write response to client: %w", err)
	}

	return false, nil
}

func (h *Handler) fail(c net.Conn) {
	_, err := c.Write(ErrResponse)
	if err != nil {
		return
	}
}
