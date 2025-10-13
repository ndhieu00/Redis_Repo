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

## Set Commands

### SADD
Add one or more members to a set.

```bash
127.0.0.1:3000> SADD myset "member1" "member2" "member3"
(integer) 3
127.0.0.1:3000> SADD myset "member1" "member4"
(integer) 1
```

### SMEMBERS
Get all members of a set.

```bash
127.0.0.1:3000> SMEMBERS myset
1) "member1"
2) "member2"
3) "member3"
4) "member4"
127.0.0.1:3000> SMEMBERS empty
(empty array)
```

### SMISMEMBER
Check if one or more members exist in a set.

```bash
127.0.0.1:3000> SMISMEMBER myset "member1" "member5"
1) (integer) 1
2) (integer) 0
127.0.0.1:3000> SMISMEMBER nonexistent "member1"
1) (integer) 0
```

### SREM
Remove one or more members from a set.

```bash
127.0.0.1:3000> SREM myset "member1" "member2"
(integer) 2
127.0.0.1:3000> SREM myset "nonexistent"
(integer) 0
```

### SCARD
Get the cardinality (number of elements) of a set.

```bash
127.0.0.1:3000> SCARD myset
(integer) 2
127.0.0.1:3000> SCARD nonexistent
(integer) 0
```

### SINTER
Get the intersection of multiple sets.

```bash
127.0.0.1:3000> SADD set1 "a" "b" "c"
(integer) 3
127.0.0.1:3000> SADD set2 "b" "c" "d"
(integer) 3
127.0.0.1:3000> SINTER set1 set2
1) "b"
2) "c"
127.0.0.1:3000> SINTER set1 set2 set3
(empty array)
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
