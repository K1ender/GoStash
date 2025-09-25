package handler

import (
	"net"
)

type Response interface {
	Serialize() ([]byte, error)
}

type CommandHandler interface {
	Handle(command string) (Response, error)
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

func (h *Handler) Handle(client net.Conn) {
	cmd := [1024]byte{}

	_, err := client.Read(cmd[:])
	if err != nil {
		client.Write([]byte(ErrResponse))
		client.Close()
		return
	}

	switch HandlerCommand(cmd[:3]) {
	case GetCommand:
		handler := h.handlers[GetCommand]
		response, err := handler.Handle(string(cmd[:]))
		if err != nil {
			client.Write([]byte(ErrResponse))
			client.Close()
			return
		}

		data, err := response.Serialize()
		if err != nil {
			client.Write([]byte(ErrResponse))
			client.Close()
			return
		}

		client.Write(data)
	default:
		client.Write([]byte(ErrResponse))
	}
}
