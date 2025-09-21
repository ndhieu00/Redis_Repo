package main

import (
	"io"
	"log"
	"net"
)

const SERVER_ADDRESS = "0.0.0.0:3000"

func main() {
	listener, err := net.Listen("tcp", SERVER_ADDRESS)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Start on", SERVER_ADDRESS)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Accept error:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Println("Client from:", conn.RemoteAddr())

	for {
		reqStr, err := readConnection(conn)
		if err != nil {
			if err == io.EOF {
				log.Println("Client disconnected:", conn.RemoteAddr())
			} else {
				log.Println("Read error:", err)
			}
			break
		}

		err = writeConnection(reqStr, conn)
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}

func readConnection(c net.Conn) (string, error) {
	buf := make([]byte, 1000)
	n, err := c.Read(buf[:])
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

func writeConnection(content string, c net.Conn) error {
	_, err := c.Write([]byte(content))
	return err
}
