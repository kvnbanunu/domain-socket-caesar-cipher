# Unix Domain Socket - Caesar Cipher

This is a client-server application that uses domain sockets for communication.

The server sets creates a sockfile at the given path to accept client connections.

The client can send a message along with a shift amount to the server where each letter will then be shifted by that amount before responding.

Case is kept and special characters are ignored.

---

## Setup
1. Clone repo
```sh
git clone https://github.com/kvnbanunu/uds-caesar-cipher
```
2. Build using make
```sh
make build-all
```

or

Build with Go

```sh
go build cmd/server/main.go -o bin/server
go build cmd/client/main.go -o bin/client
cp config.json bin/
```

---

## Run
1. Start Server
```sh
./bin/server -p <path to sockfile>
```
2. Send request with Client
```sh
./bin/client -p <path to sockfile> -i <message content> -s <shift value>
```
Both programs can be run with a -d flag to display debug statements

---

## Config
config.json includes two fields that can be changed (You do not need to rebuild)

- BufferSize sets the size of the buffer for read/write

- CipherLimit sets the loop threshold for the caesar cipher (Set to 26 for the full alphabet)
