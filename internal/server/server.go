package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"redis-repo/internal/config"
	"redis-repo/internal/core/io_multiplexing"
	"redis-repo/internal/handler/client"
	"syscall"
)

func RunRedisServer() {
	log.Println("Starting an I/O Multiplexing TCP server on", config.Port)

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
	listener, err := net.Listen(config.Protocol, config.Port)
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
			if err != syscall.EINTR {
				// EINTR is expected when the system call is interrupted by a signal
				// and should not log it as an error
				log.Println("Wait failed:", err)
			}
			continue
		}

		for _, event := range events {
			if event.Fd == int32(serverFd) {
				client.HandleNewConnection(serverFd, ioMultiplexer)
			} else {
				client.HandleClientData(int(event.Fd))
			}
		}
	}
}
