package handler

import (
	"bytes"
	"errors"
	"strconv"

	"github.com/k1ender/go-stash/internal/constants"
	"github.com/k1ender/go-stash/internal/store"
)

// GetRequest
// GET\03\0key\r\n
// Format explanation:
// - Command: "GET"
// - KeyLen: length of the key (e.g., 3 for "key")
// - Key: the actual key string (e.g., "key")
type GetRequest struct {
	Command string
	KeyLen  int
	Key     string
}

func (r *GetRequest) Serialize() []byte {
	var buf bytes.Buffer
	buf.WriteString(r.Command)
	buf.WriteByte(0)
	buf.WriteString(strconv.Itoa(r.KeyLen))
	buf.WriteByte(0)
	buf.WriteString(r.Key)
	buf.WriteString("\r\n")
	return buf.Bytes()
}

// DeserializeGet parses a byte slice into a GetRequest struct.
// The expected format of the input data is:
//
//	<command>\x00<keyLen>\x00<key>
//
// where <command> is a string, <keyLen> is the length of the key as a string,
// and <key> is the key itself. Each section is separated by a null byte (0).
// Returns an error if the format is invalid or if the key length does not match
// the provided value.
func DeserializeGet(data []byte) (*GetRequest, error) {
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

	return &GetRequest{
		Command: command,
		KeyLen:  keyLen,
		Key:     key,
	}, nil
}

type GetResponse struct {
	Value string
}

func (r *GetResponse) Serialize() ([]byte, error) {
	return []byte(r.Value + "\r\n"), nil
}

type GetHandler struct {
	store store.Store
}

func NewGetHandler(store store.Store) *GetHandler {
	return &GetHandler{
		store: store,
	}
}

func (h *GetHandler) Handle(command []byte) (Response, error) {
	cmd, err := DeserializeGet(command)
	if err != nil {
		return nil, err
	}

	value, err := h.store.Get(cmd.Key)
	if err != nil {
		return nil, err
	}

	return &GetResponse{Value: value}, nil
}
