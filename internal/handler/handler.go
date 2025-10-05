package handler

import (
	"errors"
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
	handlers map[HandlerCommand]CommandHandler
}

func NewHandler(store store.Store) *Handler {
	handlers := make(map[HandlerCommand]CommandHandler)

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

// Handle processes an incoming client connection by reading a command from the client,
// dispatching it to the appropriate handler based on the command type, and writing the
// response back to the client. If an error occurs at any stage, an error response is sent
// and the connection is closed. Currently, only the GetCommand is supported; all other
// commands result in an error response.
func (h *Handler) Handle(client net.Conn) error {
	cmd := [1024]byte{}

	_, err := client.Read(cmd[:])
	if errors.Is(err, io.EOF) {
		return nil
	}
	if err != nil {
		h.fail(client)
		return err
	}

	var response Response

	switch HandlerCommand(cmd[:3]) {
	case GetCommand:
		slog.Debug("Received GET command")
		handler := h.handlers[GetCommand]
		response, err = handler.Handle(cmd[:])
		if err != nil {
			h.fail(client)
			return nil
		}
	case SetCommand:
		slog.Debug("Received SET command")
		handler := h.handlers[SetCommand]
		response, err = handler.Handle(cmd[:])
		if err != nil {
			h.fail(client)
			return nil
		}
	case IncrCommand:
		slog.Debug("Received INC command")
		handler := h.handlers[IncrCommand]
		response, err = handler.Handle(cmd[:])
		if err != nil {
			h.fail(client)
			return nil
		}
	case DecrCommand:
		slog.Debug("Received DEC command")
		handler := h.handlers[DecrCommand]
		response, err = handler.Handle(cmd[:])
		if err != nil {
			h.fail(client)
			return nil
		}
	case DelCommand:
		slog.Debug("Received DEL command")
		handler := h.handlers[DelCommand]
		response, err = handler.Handle(cmd[:])
		if err != nil {
			h.fail(client)
			return nil
		}
	default:
		h.fail(client)
		return nil
	}

	data, err := response.Serialize()
	if err != nil {
		h.fail(client)
		return nil
	}

	client.Write(data)
	return nil
}

func (h *Handler) fail(c net.Conn) {
	c.Write(ErrResponse)
}
