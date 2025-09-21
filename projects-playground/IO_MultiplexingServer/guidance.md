# I/O Multiplexing and Epoll

**I/O multiplexing** is a programming model that allows a single process or thread to monitor multiple I/O sources (file descriptors) — like sockets, pipes, or files — and react only when one or more of them become ready for reading or writing.

This model is particularly useful for building high-performance, non-blocking network servers that handle many simultaneous connections efficiently, without using one thread per connection.

In Linux, I/O multiplexing is implemented through several system calls:

- `select` — legacy and limited in scalability
- `poll` — more flexible but still inefficient for large numbers of FDs
- ✅ **`epoll`** — modern, efficient, and scalable

## Epoll: Linux's I/O Multiplexing Mechanism

**`epoll`** was introduced in Linux as a high-performance mechanism for I/O multiplexing. It allows a program to monitor thousands (or even millions) of file descriptors with minimal overhead, making it the preferred choice for scalable server applications.

### How epoll Works

- Create an **epoll instance** using `epoll_create1`.
- **Register file descriptors** (e.g., sockets) with `epoll_ctl`.
- **Wait for events** using `epoll_wait`.

### epoll Operations

- `EPOLL_CTL_ADD`: Add a file descriptor to monitor.
- `EPOLL_CTL_MOD`: Modify an existing monitored descriptor.
- `EPOLL_CTL_DEL`: Remove a monitored descriptor.

### Common Event Flags

- `EPOLLIN`: Ready for read
- `EPOLLOUT`: Ready for write
- `EPOLLERR`, `EPOLLHUP`: Error or hangup

### Note on `CLOEXEC` flag for file descriptor
`CLOEXEC` closes file descriptors on `execve` (which replaces the process with a new program), preventing leaks into child processes.