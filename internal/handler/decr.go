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
		return nil, errors.New("invalid format: no second delimiter")
	}

	i2 += i1 + 1

	command := string(data[:i1])
	keyLenStr := string(data[i1+1 : i2])
	keyLen, err := strconv.Atoi(keyLenStr)
	if err != nil {
		return nil, err
	}

	key := string(data[i2+1 : i2+keyLen+1])
	if len(key) != keyLen {
		return nil, errors.New("key length does not match the provided length")
	}

	return &DecrRequest{
		Command: command,
		KeyLen:  keyLen,
		Key:     key,
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
