package main

import (
	"example/io_multiplexing_server/io_multiplexing"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"syscall"
)

const (
	SERVER_ADDRESS = "0.0.0.0:3000"
	BUFFER_SIZE    = 512
)

func main() {
	log.Println("Starting an I/O Multiplexing TCP server on", SERVER_ADDRESS)

	listener, listenerFile, serverFd, err := setupServer()
	if err != nil {
		log.Fatal("Server setup failed:", err)
	}
	defer listener.Close()
	defer listenerFile.Close()

	ioMultiplexer, err := setupIOMultiplexer(serverFd)
	if err != nil {
		log.Fatal("IO Multiplexer setup failed:", err)
	}
	defer ioMultiplexer.Close()

	runEventLoop(ioMultiplexer, serverFd)
}

// setupServer creates and configures the TCP listener, returning the listener,
// file descriptor, and server file descriptor for epoll monitoring
func setupServer() (net.Listener, *os.File, int, error) {
	listener, err := net.Listen("tcp", SERVER_ADDRESS)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to start listener: %w", err)
	}

	// Get the file descriptor from the listener
	tcpListener, ok := listener.(*net.TCPListener)
	if !ok {
		listener.Close()
		return nil, nil, 0, fmt.Errorf("listener is not a TCP Listener")
	}

	listenerFile, err := tcpListener.File()
	if err != nil {
		listener.Close()
		return nil, nil, 0, fmt.Errorf("failed to get listener file: %w", err)
	}

	serverFd := int(listenerFile.Fd())
	return listener, listenerFile, serverFd, nil
}

// setupIOMultiplexer creates and configures the IO multiplexer, monitoring the server file descriptor
func setupIOMultiplexer(serverFd int) (*io_multiplexing.Epoll, error) {
	ioMultiplexer, err := io_multiplexing.CreateIOMultiplexer()
	if err != nil {
		return nil, fmt.Errorf("failed to create IO multiplexer: %w", err)
	}

	if err = ioMultiplexer.Monitor(syscall.EpollEvent{
		Fd:     int32(serverFd),
		Events: syscall.EPOLLIN,
	}); err != nil {
		ioMultiplexer.Close()
		return nil, fmt.Errorf("failed to monitor server file descriptor: %w", err)
	}

	return ioMultiplexer, nil
}

// runEventLoop continuously waits for and processes IO events from the multiplexer
func runEventLoop(ioMultiplexer *io_multiplexing.Epoll, serverFd int) {
	for {
		events, err := ioMultiplexer.Wait()
		if err != nil {
			log.Println("Wait failed:", err)
			continue
		}

		for _, event := range events {
			if event.Fd == int32(serverFd) {
				handleNewConnection(serverFd, ioMultiplexer)
			} else {
				handleClientData(int(event.Fd))
			}
		}
	}
}

// handleNewConnection accepts a new client connection and adds it to the IO multiplexer monitoring
func handleNewConnection(serverFd int, ioMultiplexer *io_multiplexing.Epoll) {
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

// handleClientData reads commands from a client connection and sends responses
func handleClientData(clientFd int) {
	cmd, err := readCommand(clientFd)
	if err != nil {
		if err == io.EOF || err == syscall.ECONNRESET {
			syscall.Close(clientFd)
			return
		}
		log.Println("Read Error:", err)
		return
	}

	if err = respond(cmd, clientFd); err != nil {
		log.Println("Write Error:", err)
	}
}

func readCommand(fd int) (string, error) {
	buf := make([]byte, BUFFER_SIZE)
	n, err := syscall.Read(fd, buf)
	if err != nil {
		return "", err
	}
	if n == 0 {
		return "", io.EOF
	}

	return string(buf[:n]), nil
}

func respond(data string, fd int) error {
	_, err := syscall.Write(fd, []byte(data))
	if err != nil {
		return err
	}
	return nil
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
