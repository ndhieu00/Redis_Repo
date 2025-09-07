package main

import (
	"example/io_multiplexing_server/io_multiplexing"
	"fmt"
	"io"
	"log"
	"net"
	"syscall"
)

const SERVER_ADDRESS = "0.0.0.0:3000"

func main() {
	log.Println("Starting an I/O Multiplexing TCP server on", SERVER_ADDRESS)
	listener, err := net.Listen("tcp", SERVER_ADDRESS)
	if err != nil {
		log.Fatal("Start failed:", err)
		return
	}
	defer listener.Close()

	// Get the file descriptor from the listener
	tcpListener, ok := listener.(*net.TCPListener)
	if !ok {
		log.Fatal("Listener is not a TCP Listener")
	}
	listenerFile, err := tcpListener.File()
	if err != nil {
		log.Fatal(err)
	}
	defer listenerFile.Close()

	serverFd := int(listenerFile.Fd())

	ioMultiplexer, err := io_multiplexing.CreateIOMultiplexer()
	if err != nil {
		log.Fatal("Create IO Multiplexer failed:", err)
	}
	defer ioMultiplexer.Close()

	if err = ioMultiplexer.Monitor(syscall.EpollEvent{
		Fd:     int32(serverFd),
		Events: syscall.EPOLLIN,
	}); err != nil {
		log.Fatal("Monitor serverFd failed:", err)
	}

	for {
		events, err := ioMultiplexer.Wait()
		if err != nil {
			log.Println("Wait failed:", err)
			continue
		}
		for _, event := range events {
			if event.Fd == int32(serverFd) {
				// This means a new client is trying connect
				connFd, sa, err := syscall.Accept(serverFd)
				formatedAddress := formatSockaddr(sa)
				if err != nil {
					log.Println("Accept connection failed,", err)
					continue
				}

				log.Println("New connection from:", formatedAddress)
				if err = ioMultiplexer.Monitor(syscall.EpollEvent{
					Fd:     int32(connFd),
					Events: syscall.EPOLLIN,
				}); err != nil {
					log.Println("Monitor connection", formatedAddress, "failed:", err)
					syscall.Close(connFd)
				}
			} else {
				// Read commands from connections
				cmd, err := readCommand(int(event.Fd))
				if err != nil {
					if err == io.EOF || err == syscall.ECONNRESET {
						syscall.Close(int(event.Fd))
						continue
					}

					log.Println("Read Error:", err)
					continue
				}

				if err = respond(cmd, int(event.Fd)); err != nil {
					log.Println("Write Error:", err)
				}
			}
		}
	}
}

func readCommand(fd int) (string, error) {
	buf := make([]byte, 512)
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
