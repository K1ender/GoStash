package handler

import (
	"bytes"
	"errors"
	"strconv"

	"github.com/k1ender/go-stash/internal/constants"
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
	i1 := constants.CommandKeyLen
	i2 := bytes.IndexByte(data[i1+1:], 0)
	if i2 == -1 {
		return nil, errors.New("invalid format")
	}
	i2 += i1 + 1

	i3 := bytes.IndexByte(data[i2+1:], 0)
	if i3 == -1 {
		return nil, errors.New("invalid format")
	}
	i3 += i2 + 1

	i4 := bytes.IndexByte(data[i3+1:], 0)
	if i4 == -1 {
		return nil, errors.New("invalid format")
	}
	i4 += i3 + 1

	keyLen, err := strconv.Atoi(string(data[i1+1 : i2]))
	if err != nil {
		return nil, err
	}
	key := data[i2+1 : i2+1+keyLen]

	valueLen, err := strconv.Atoi(string(data[i3+1 : i4]))
	if err != nil {
		return nil, err
	}
	value := data[i4+1 : i4+1+valueLen]

	return &SetRequest{Command: string(data[:i1]), Key: string(key), Value: string(value)}, nil
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
