package executor

import (
	"redis-repo/internal/core/command"
	"syscall"
)

func ExecuteAndRespond(cmd *command.Command, clientFd int) error {
	var res []byte
	switch cmd.Cmd {
	case "PING":
		res = executePing(cmd.Args)
	default:
		res = []byte("-CMD NOT FOUND\r\n")
	}

	return Respond(string(res), clientFd)
}

func Respond(data string, fd int) error {
	_, err := syscall.Write(fd, []byte(data))
	if err != nil {
		return err
	}
	return nil
}
