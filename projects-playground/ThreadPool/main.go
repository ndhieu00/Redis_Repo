package main

import (
	"io"
	"log"
	"net"
)

const NUM_WORKERS = 10
const SERVER_ADDRESS = "0.0.0.0:3000"

type Job struct {
	conn net.Conn
}

type Worker struct {
	id         int
	jobChannel chan Job
}

type Pool struct {
	jobQueue chan Job
	workers  []*Worker
}

// -----------------------------------------------------------

func NewWorker(id int, jobChannel chan Job) *Worker {
	return &Worker{
		id:         id,
		jobChannel: jobChannel,
	}
}

func (w *Worker) Start() {
	go (func() {
		for job := range w.jobChannel {
			w.handleConnection(job.conn)
		}
	})()
}

func (w *Worker) handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Println("Accepted Connection from:", conn.RemoteAddr(), "by Worker:", w.id)

	for {
		content, err := readConnection(conn)
		if err != nil {
			if err == io.EOF {
				log.Println("Close connection from client:", conn.RemoteAddr())
			} else {
				log.Println("Error:", err)
			}
			break
		}

		err = writeConnection(conn, "--: "+content)
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}

// -----------------------------------------------------------------

func NewPool(numWorkers int) *Pool {
	pool := Pool{}
	pool.jobQueue = make(chan Job) // Can be buffered to queue jobs without blocking the main thread
	pool.workers = make([]*Worker, numWorkers)

	for i := range numWorkers {
		pool.workers[i] = NewWorker(i, pool.jobQueue)
	}

	return &pool
}

func (p *Pool) Start() {
	for _, worker := range p.workers {
		worker.Start()
	}
}

func (p *Pool) Add(job Job) {
	if p.jobQueue == nil {
		log.Println("Cannot add job. Pool's job queue has not been initialized.")
		return
	}
	p.jobQueue <- job
}

// -----------------------------------------------------------

func main() {
	listener, err := net.Listen("tcp", SERVER_ADDRESS)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Start on", SERVER_ADDRESS)

	pool := NewPool(NUM_WORKERS)
	pool.Start()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Accept error:", err)
			continue
		}
		log.Println("A new connection from:", conn.RemoteAddr())
		pool.Add(Job{conn: conn})
	}
}

// -----------------------------------------------------------

func readConnection(conn net.Conn) (string, error) {
	buf := make([]byte, 1000)
	n, err := conn.Read(buf[:])
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

func writeConnection(c net.Conn, content string) error {
	_, err := c.Write([]byte(content))
	return err
}
