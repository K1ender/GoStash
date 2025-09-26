package handler

import (
	"net"
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

func NewHandler() *Handler {
	handlers := make(map[HandlerCommand]CommandHandler)
	getHandler := NewGetHandler()
	handlers[GetCommand] = getHandler

	return &Handler{
		handlers: handlers,
	}
}

// Handle processes an incoming client connection by reading a command from the client,
// dispatching it to the appropriate handler based on the command type, and writing the
// response back to the client. If an error occurs at any stage, an error response is sent
// and the connection is closed. Currently, only the GetCommand is supported; all other
// commands result in an error response.
func (h *Handler) Handle(client net.Conn) {
	cmd := [1024]byte{}

	_, err := client.Read(cmd[:])
	if err != nil {
		h.fail(client)
		return
	}

	var response Response

	switch HandlerCommand(cmd[:3]) {
	case GetCommand:
		handler := h.handlers[GetCommand]
		response, err = handler.Handle(cmd[:])
		if err != nil {
			h.fail(client)
			return
		}
	default:
		h.fail(client)
		return
	}

	data, err := response.Serialize()
	if err != nil {
		h.fail(client)
		return
	}

	client.Write(data)
}

func (h *Handler) fail(c net.Conn) {
	c.Write(ErrResponse)
}
