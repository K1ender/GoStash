package handler

// Command always in uppercase string with len 3
type Command [3]byte

var (
	GetCommand  Command = Command{'G', 'E', 'T'}
	SetCommand  Command = Command{'S', 'E', 'T'}
	IncrCommand Command = Command{'I', 'N', 'C'}
	DecrCommand Command = Command{'D', 'E', 'C'}
	DelCommand  Command = Command{'D', 'E', 'L'}
)
