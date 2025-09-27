package command

import (
	"io"
	"redis-repo/internal/core/resp"
	"strings"
	"syscall"
)

// Command represents a Redis command with its name and arguments
type Command struct {
	Cmd  string
	Args []string
}

func ParseCmd(data []byte) (*Command, error) {
	value, err := resp.Decode(data)
	if err != nil {
		return nil, err
	}

	array := value.([]any)
	tokens := make([]string, len(array))
	for i := range tokens {
		tokens[i] = array[i].(string)
	}

	res := &Command{
		Cmd:  strings.ToUpper(tokens[0]),
		Args: tokens[1:],
	}
	return res, nil
}

func ReadCommand(fd int) (*Command, error) {
	buf := make([]byte, 512)
	n, err := syscall.Read(fd, buf)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, io.EOF
	}

	cmd, err := ParseCmd(buf[:n])
	if err != nil {
		return nil, err
	}

	return cmd, nil
}
