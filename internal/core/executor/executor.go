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
	case "SADD":
		res = cmdSADD(cmd.Args)
	case "SREM":
		res = cmdSREM(cmd.Args)
	case "SMISMEMBER":
		res = cmdSMISMEMBER(cmd.Args)
	case "SMEMBERS":
		res = cmdSMEMBERS(cmd.Args)
	case "SCARD":
		res = cmdSCARD(cmd.Args)
	case "SINTER":
		res = cmdSINTER(cmd.Args)
	default:
		res = []byte("-CMD NOT FOUND\r\n")
	}

	_, err := syscall.Write(clientFd, res)
	return err
}
