package handler

import (
	"bytes"
	"errors"
	"strconv"

	"github.com/k1ender/go-stash/internal/constants"
	"github.com/k1ender/go-stash/internal/store"
)

type IncrRequest struct {
	Command string
	KeyLen  int
	Key     string
}

func (r *IncrRequest) Serialize() []byte {
	var buf bytes.Buffer
	buf.WriteString(r.Command)
	buf.WriteByte(0)
	buf.WriteString(strconv.Itoa(r.KeyLen))
	buf.WriteByte(0)
	buf.WriteString(r.Key)
	buf.WriteString("\r\n")
	return buf.Bytes()
}

func DeserializeIncr(data []byte) (*IncrRequest, error) {
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

	return &IncrRequest{
		Command: command,
		KeyLen:  keyLen,
		Key:     key,
	}, nil
}

type IncrResponse struct {
	Value int
}

func (r *IncrResponse) Serialize() ([]byte, error) {
	return []byte(strconv.Itoa(r.Value) + "\r\n"), nil
}

type IncrHandler struct {
	store store.Store
}

func NewIncrHandler(store store.Store) *IncrHandler {
	return &IncrHandler{store: store}
}

func (h *IncrHandler) Handle(command []byte) (Response, error) {
	request, err := DeserializeIncr(command)
	if err != nil {
		return nil, err
	}

	val, err := h.store.Incr(request.Key)
	if err != nil {
		return nil, err
	}

	return &IncrResponse{Value: val}, nil
}
