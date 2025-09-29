package handler

import (
	"bytes"
	"errors"
	"strconv"

	"github.com/k1ender/go-stash/internal/constants"
	"github.com/k1ender/go-stash/internal/store"
)

// DEL\0<keyLen>\0<key>\r\n
// Example: DEL\0\x03\0foo\r\n
// Command: "DEL"
// KeyLen: length of the key (e.g., 3 for "foo")
// Key: the actual key string (e.g., "foo")
type DelRequest struct {
	Command string
	KeyLen  int
	Key     string
}

func (r *DelRequest) Serialize() []byte {
	var buf bytes.Buffer
	buf.WriteString(r.Command)
	buf.WriteByte(0)
	buf.WriteString(strconv.Itoa(r.KeyLen))
	buf.WriteByte(0)
	buf.WriteString(r.Key)
	buf.WriteString("\r\n")
	return buf.Bytes()
}

func DeserializeDel(data []byte) (*DelRequest, error) {
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
		return nil, errors.New("key length mismatch")
	}

	return &DelRequest{
		Command: command,
		KeyLen:  keyLen,
		Key:     key,
	}, nil
}

type DelResponse struct {
	Value string
}

func (r *DelResponse) Serialize() ([]byte, error) {
	return []byte(r.Value + "\r\n"), nil
}

type DelHandler struct {
	store store.Store
}

func NewDelHandler(store store.Store) *DelHandler {
	return &DelHandler{store: store}
}

func (h *DelHandler) Handle(command []byte) (Response, error) {
	req, err := DeserializeDel(command)
	if err != nil {
		return nil, err
	}

	err = h.store.Del(req.Key)
	if err != nil {
		return nil, err
	}

	return &DelResponse{Value: "OK"}, nil
}
