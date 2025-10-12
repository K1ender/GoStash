package handler

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"

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
	cmd, err := io.ReadAll(client)
	if err != nil {
		h.fail(client)
		return true, fmt.Errorf("failed to read command from client: %w", err)
	}

	var response Response

	switch Command(cmd[:3]) {
	case GetCommand:
		slog.Debug("Received GET command")
		handler := h.handlers[GetCommand]
		response, err = handler.Handle(cmd[:])
		if err != nil {
			h.fail(client)
			return false, fmt.Errorf("failed to handle GET command: %w", err)
		}
	case SetCommand:
		slog.Debug("Received SET command")
		handler := h.handlers[SetCommand]
		response, err = handler.Handle(cmd[:])
		if err != nil {
			h.fail(client)
			return false, fmt.Errorf("failed to handle SET command: %w", err)
		}
	case IncrCommand:
		slog.Debug("Received INC command")
		handler := h.handlers[IncrCommand]
		response, err = handler.Handle(cmd[:])
		if err != nil {
			h.fail(client)
			return false, fmt.Errorf("failed to handle INC command: %w", err)
		}
	case DecrCommand:
		slog.Debug("Received DEC command")
		handler := h.handlers[DecrCommand]
		response, err = handler.Handle(cmd[:])
		if err != nil {
			h.fail(client)
			return false, fmt.Errorf("failed to handle DEC command: %w", err)
		}
	case DelCommand:
		slog.Debug("Received DEL command")
		handler := h.handlers[DelCommand]
		response, err = handler.Handle(cmd[:])
		if err != nil {
			h.fail(client)
			return false, fmt.Errorf("failed to handle DEL command: %w", err)
		}
	default:
		h.fail(client)
		return false, errors.New("unknown command")
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
