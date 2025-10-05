# Architecture Overview

High-level system design of the Redis server.

## System Flow

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Server    │ --> │   Handler   │ --> │  Executor   │
│             │     │             │     │             │
│ • Accept    │     │ • Parse     │     │ • Execute   │
│ • Monitor   │     │ • Validate  │     │ • Respond   │
│ • Route     │     │ • Route     │     │ • Write     │
│ • System    │     │ • Client    │     │ • Cleanup   │
│   Events    │     │ • System    │     │             │
└─────────────┘     └─────────────┘     └─────────────┘
```

## Layer Responsibilities

### Server Layer
- **Network I/O**: TCP socket management and connection handling
- **Event Loop**: epoll-based I/O multiplexing for efficient event handling
- **Connection Management**: Accept new connections and monitor existing ones
- **System Events**: Triggers system-level operations (cleanup)

### Handler Layer
- **Client Handler**: Command parsing, client connection management, RESP protocol
- **Server Handler**: System-level operations (cleanup)

### Executor Layer
- **Command Execution**: Business logic for each Redis command
- **Response Generation**: RESP protocol encoding
- **System Operations**: Expired key cleanup
## Core Components

### I/O Multiplexing
Uses Linux epoll for efficient event-driven I/O, handling thousands of concurrent connections.

### RESP Protocol
Implements the Redis Serialization Protocol for client-server communication.

### Data Structures
Custom dictionary implementation with TTL support for key-value storage.

## Project Structure

```
internal/
├── core/
│   ├── command/          # Command type definitions
│   ├── executor/         # Command execution logic
│   ├── resp/            # RESP protocol encoding/decoding
│   └── io_multiplexing/ # epoll-based I/O multiplexing
├── data_structure/      # Custom data structures
├── handler/
│   ├── client/          # Client connection handling
│   └── server/          # System-level operations
├── server/              # Main server implementation
├── config/              # Configuration management
└── constant/            # Application constants
```