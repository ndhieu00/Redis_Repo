package client

import (
	"fmt"
	"io"
	"log"
	"net"
	"redis-repo/internal/core/command"
	"redis-repo/internal/core/executor"
	"redis-repo/internal/core/io_multiplexing"
	"syscall"
)

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
	cmd, err := command.ReadCommand(clientFd)
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
