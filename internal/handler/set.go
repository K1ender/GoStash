package handler

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/k1ender/go-stash/internal/store"
)

// format
// SET\0<keyLen>\0<key>\0<valueLen>\0<value>\r\n
type SetRequest struct {
	Command  string
	KeyLen   int
	Key      string
	ValueLen int
	Value    string
}

func (r *SetRequest) Serialize() []byte {
	var buf bytes.Buffer
	buf.WriteString(r.Command)
	buf.WriteByte(0)
	buf.WriteString(strconv.Itoa(r.KeyLen))
	buf.WriteByte(0)
	buf.WriteString(r.Key)
	buf.WriteByte(0)
	buf.WriteString(strconv.Itoa(r.ValueLen))
	buf.WriteByte(0)
	buf.WriteString(r.Value)
	buf.WriteString("\r\n")
	return buf.Bytes()
}

func DeserializeSet(data []byte) (*SetRequest, error) {
	// SET\0<keyLen>\0<key>\0<valueLen>\0<value>\r\n
	i1 := 3
	if i1 == -1 {
		return nil, fmt.Errorf("invalid format: no first delimiter")
	}

	i2 := bytes.IndexByte(data[i1+1:], 0)
	if i2 == -1 {
		return nil, fmt.Errorf("invalid format: no second delimiter")
	}

	i2 += i1 + 1

	i3 := bytes.IndexByte(data[i2+1:], 0)
	if i3 == -1 {
		return nil, fmt.Errorf("invalid format: no third delimiter")
	}

	i3 += i2 + 1

	i4 := bytes.IndexByte(data[i3+1:], 0)
	if i4 == -1 {
		return nil, fmt.Errorf("invalid format: no fourth delimiter")
	}

	i4 += i3 + 1

	command := string(data[:i1])
	keyLenStr := string(data[i1+1 : i2])
	keyLen, err := strconv.Atoi(keyLenStr)
	if err != nil {
		return nil, err
	}

	key := string(data[i2+1 : i2+keyLen+1])
	if len(key) != keyLen {
		return nil, fmt.Errorf("key length mismatch")
	}

	valueLenStr := string(data[i3+1 : i4])
	valueLen, err := strconv.Atoi(valueLenStr)
	if err != nil {
		return nil, err
	}

	value := string(data[i4+1 : i4+valueLen+1])
	if len(value) != valueLen {
		return nil, fmt.Errorf("value length mismatch")
	}

	return &SetRequest{
		Command:  command,
		KeyLen:   keyLen,
		Key:      key,
		ValueLen: valueLen,
		Value:    value,
	}, nil
}

type SetResponse struct {
	Value string
}

func (r *SetResponse) Serialize() ([]byte, error) {
	return []byte(r.Value + "\r\n"), nil
}

type SetHandler struct {
	store store.Store
}

func NewSetHandler(store store.Store) *SetHandler {
	return &SetHandler{
		store: store,
	}
}

func (h *SetHandler) Handle(command []byte) (Response, error) {
	cmd, err := DeserializeSet(command)
	if err != nil {
		return nil, err
	}

	err = h.store.Set(cmd.Key, cmd.Value)
	if err != nil {
		return nil, err
	}

	return &SetResponse{Value: "OK"}, nil
}
