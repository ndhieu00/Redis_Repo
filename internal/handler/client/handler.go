package client

import (
	"fmt"
	"io"
	"log"
	"net"
	"redis-repo/internal/core/command"
	"redis-repo/internal/core/executor"
	"redis-repo/internal/core/io_multiplexing"
	"redis-repo/internal/core/resp"
	"strings"
	"syscall"
)

// parseCmd parses RESP data into a Command struct
func parseCmd(data []byte) (*command.Command, error) {
	value, err := resp.Decode(data)
	if err != nil {
		return nil, err
	}

	array := value.([]any)
	tokens := make([]string, len(array))
	for i := range tokens {
		tokens[i] = array[i].(string)
	}

	res := &command.Command{
		Cmd:  strings.ToUpper(tokens[0]),
		Args: tokens[1:],
	}
	return res, nil
}

// readCommand reads a command from a file descriptor
func readCommand(fd int) (*command.Command, error) {
	buf := make([]byte, 512)
	n, err := syscall.Read(fd, buf)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, io.EOF
	}

	cmd, err := parseCmd(buf[:n])
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

// HandleNewConnection accepts a new client connection and adds it to the IO multiplexer monitoring
func HandleNewConnection(serverFd int, ioMultiplexer *io_multiplexing.Epoll) {
	connFd, sa, err := syscall.Accept(serverFd)
	formattedAddress := formatSockaddr(sa)
	if err != nil {
		log.Println("Accept connection failed:", err)
		return
	}

	log.Println("New connection from:", formattedAddress)
	if err = ioMultiplexer.Monitor(syscall.EpollEvent{
		Fd:     int32(connFd),
		Events: syscall.EPOLLIN,
	}); err != nil {
		log.Println("Monitor connection", formattedAddress, "failed:", err)
		syscall.Close(connFd)
	}
}

// HandleClientData reads commands from a client connection and sends responses
func HandleClientData(clientFd int) {
	cmd, err := readCommand(clientFd)
	if err != nil {
		if err == io.EOF || err == syscall.ECONNRESET {
			syscall.Close(clientFd)
			return
		}
		log.Println("Read Error:", err)
		return
	}

	if err = executor.ExecuteAndRespond(cmd, clientFd); err != nil {
		log.Println("Execute and respond failed:", err)
	}
}

func formatSockaddr(sa syscall.Sockaddr) string {
	switch a := sa.(type) {
	case *syscall.SockaddrInet4:
		ip := net.IPv4(a.Addr[0], a.Addr[1], a.Addr[2], a.Addr[3])
		return fmt.Sprintf("%s:%d", ip, a.Port)
	case *syscall.SockaddrInet6:
		ip := net.IP(a.Addr[:])
		return fmt.Sprintf("[%s]:%d", ip, a.Port)
	default:
		return fmt.Sprintf("%v", sa)
	}
}
