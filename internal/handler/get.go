package handler

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

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
	return buf.Bytes()
}

func DeserializeGet(data []byte) (*GetRequest, error) {
	i1 := bytes.IndexByte(data, 0)
	if i1 == -1 {
		return nil, errors.New("invalid format: no first delimiter")
	}

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

	key := string(data[i2+1 : i2+keyLen+1]) // Exclude \r\n
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
	return []byte(r.Value), nil
}

type GetHandler struct {
	// todo: add fields if necessary
}

func (h *GetHandler) Handle(command string) (Response, error) {
	cmd, err := DeserializeGet([]byte(command))
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize GET command: %w", err)
	}

	return &GetResponse{Value: cmd.Key}, nil
}

func NewGetHandler() *GetHandler {
	return &GetHandler{}
}
