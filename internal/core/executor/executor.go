package executor

import (
	"redis-repo/internal/core/command"
	"syscall"
)

func ExecuteAndRespond(cmd *command.Command, clientFd int) error {
	var res []byte

	switch cmd.Cmd {
	case "PING":
		res = cmdPING(cmd.Args)
	case "GET":
		res = cmdGET(cmd.Args)
	case "SET":
		res = cmdSET(cmd.Args)
	case "TTL":
		res = cmdTTL(cmd.Args)
	case "DEL":
		res = cmdDEL(cmd.Args)

	default:
		res = []byte("-CMD NOT FOUND\r\n")
	}

	_, err := syscall.Write(clientFd, res)
	return err
}
