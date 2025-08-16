package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:3000")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Start on localhost:3000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	fmt.Println("Connection Infor:", conn.RemoteAddr())

	buffer := make([]byte, 1000)
	_, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second * 2)

	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\nYou sent: \r\n" + string(buffer)))

	conn.Close()
}
