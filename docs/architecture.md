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
└─────────────┘     └─────────────┘     └─────────────┘
```

## Layer Responsibilities

### Server Layer
- **Network I/O**: TCP socket management and connection handling
- **Event Loop**: epoll-based I/O multiplexing for efficient event handling
- **Connection Management**: Accept new connections and monitor existing ones

### Handler Layer
- **Command Parsing**: RESP (REdis Serialization Protocol) decoding
- **Client Management**: Per-client connection state and lifecycle

### Executor Layer
- **Command Execution**: Business logic for each Redis command
- **Response Generation**: RESP protocol encoding
- **TTL Management**: Automatic expiry handling

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
├── handler/             # Client connection handling
├── server/              # Main server implementation
├── config/              # Configuration management
└── constant/            # Application constants
```