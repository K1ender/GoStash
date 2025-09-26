package handler

// Command always in uppercase string with len 3
type HandlerCommand [3]byte

var (
	GetCommand HandlerCommand = [3]byte{'G', 'E', 'T'}
	SetCommand HandlerCommand = [3]byte{'S', 'E', 'T'}
)
