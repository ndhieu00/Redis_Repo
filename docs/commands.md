# Redis Commands

Supported Redis commands in this server.

### GET
Retrieve the value of a key.

```bash
127.0.0.1:3000> GET mykey
"Hello World"
127.0.0.1:3000> GET nonexistent
(nil)
```

### SET
Set a key to hold a string value with optional expiry.

```bash
# Basic SET
127.0.0.1:3000> SET mykey "Hello World"
OK

# SET with expiry in seconds
127.0.0.1:3000> SET mykey "Hello" EX 60
OK

# SET with expiry in milliseconds
127.0.0.1:3000> SET mykey "Hello" PX 60000
OK

# SET with expiry at specific timestamp
127.0.0.1:3000> SET mykey "Hello" EXAT 9999999999
OK
```

### DEL
Delete one or more keys.

```bash
127.0.0.1:3000> DEL key1
(integer) 1
127.0.0.1:3000> DEL key1 key2 key3
(integer) 2
```

### TTL
Get the time to live for a key in seconds.

```bash
127.0.0.1:3000> SET mykey "Hello" EX 60
OK
127.0.0.1:3000> TTL mykey
(integer) 60
127.0.0.1:3000> TTL nonexistent
(integer) -2
```

## Utility Commands

### PING
Test the connection to the server.

```bash
127.0.0.1:3000> PING
PONG
127.0.0.1:3000> PING "Hello Redis"
"Hello Redis"
```
