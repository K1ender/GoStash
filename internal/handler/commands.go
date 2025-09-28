package handler

// Command always in uppercase string with len 3
type HandlerCommand [3]byte

var (
	GetCommand  HandlerCommand = HandlerCommand{'G', 'E', 'T'}
	SetCommand  HandlerCommand = HandlerCommand{'S', 'E', 'T'}
	IncrCommand HandlerCommand = HandlerCommand{'I', 'N', 'C'}
	DecrCommand HandlerCommand = HandlerCommand{'D', 'E', 'C'}
)
