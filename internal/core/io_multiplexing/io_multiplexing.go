package io_multiplexing

import (
	"redis-repo/internal/config"
	"syscall"
)

type Epoll struct {
	fd          int
	epollEvents []syscall.EpollEvent
}

func CreateIOMultiplexer() (*Epoll, error) {
	// syscall.EPOLL_CLOEXEC: This means if your process calls execve, the epoll FD will be automatically closed.
	// Flag 0 (no flags) is fine in Go, because the runtime sets CLOEXEC afterward anyway.
	epollFD, err := syscall.EpollCreate1(0)
	if err != nil {
		return nil, err
	}

	return &Epoll{
		fd:          epollFD,
		epollEvents: make([]syscall.EpollEvent, config.MaxConnection),
	}, nil
}

func (ep *Epoll) Monitor(epEvent syscall.EpollEvent) error {
	return syscall.EpollCtl(ep.fd, syscall.EPOLL_CTL_ADD, int(epEvent.Fd), &epEvent)
}

func (ep *Epoll) Wait() ([]syscall.EpollEvent, error) {
	n, err := syscall.EpollWait(ep.fd, ep.epollEvents, -1)
	if err != nil {
		return nil, err
	}
	return ep.epollEvents[:n], nil
}

func (ep *Epoll) Close() error {
	return syscall.Close(ep.fd)
}
