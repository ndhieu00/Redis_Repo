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
	case "GET":
		res = executeGet(cmd.Args)
	case "SET":
		res = executeSet(cmd.Args)
	case "TTL":
		res = executeTTL(cmd.Args)
	default:
		res = []byte("-CMD NOT FOUND\r\n")
	}

	_, err := syscall.Write(clientFd, res)
	return err
}
