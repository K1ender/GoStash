package handler

type StatusCode []byte

var (
	ErrResponse StatusCode = []byte("ERR")
)
