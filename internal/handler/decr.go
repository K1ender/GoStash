package handler

import (
	"bytes"
	"errors"
	"strconv"

	"github.com/k1ender/go-stash/internal/constants"
	"github.com/k1ender/go-stash/internal/store"
)

type DecrRequest struct {
	Command string
	KeyLen  int
	Key     string
}

func (r *DecrRequest) Serialize() []byte {
	var buf bytes.Buffer
	buf.WriteString(r.Command)
	buf.WriteByte(0)
	buf.WriteString(strconv.Itoa(r.KeyLen))
	buf.WriteByte(0)
	buf.WriteString(r.Key)
	buf.WriteString("\r\n")
	return buf.Bytes()
}

func DeserializeDecr(data []byte) (*DecrRequest, error) {
	i1 := constants.CommandKeyLen
	i2 := bytes.IndexByte(data[i1+1:], 0)
	if i2 == -1 {
		return nil, errors.New("invalid format")
	}
	i2 += i1 + 1

	keyLen, err := strconv.Atoi(string(data[i1+1 : i2]))
	if err != nil {
		return nil, err
	}

	key := data[i2+1 : i2+1+keyLen]
	return &DecrRequest{
		Command: string(data[:i1]),
		KeyLen:  keyLen,
		Key:     string(key),
	}, nil
}

type DecrResponse struct {
	Value int
}

func (r *DecrResponse) Serialize() ([]byte, error) {
	return []byte(strconv.Itoa(r.Value) + "\r\n"), nil
}

type DecrHandler struct {
	store store.Store
}

func NewDecrHandler(store store.Store) *DecrHandler {
	return &DecrHandler{store: store}
}

func (h *DecrHandler) Handle(command []byte) (Response, error) {
	req, err := DeserializeDecr(command)
	if err != nil {
		return nil, err
	}

	value, err := h.store.Decr(req.Key)
	if err != nil {
		return nil, err
	}

	return &DecrResponse{Value: value}, nil
}
