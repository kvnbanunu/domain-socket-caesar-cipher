package main

import (
	"log"

	"github.com/kvnbanunu/uds-caesar-cipher/internal/options"
	"github.com/kvnbanunu/uds-caesar-cipher/internal/socket"
)

func main() {
	args := options.Args{}
	args.ParseArgs(false)

	config, err := args.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	sock, err := socket.Ssetup(&args)
	if err != nil {
		log.Fatal(err)
	}

	socket.HandleSignal(sock, &args)

	for {
		conn, err := sock.AcceptUnix(); if err == nil {
			socket.HandleConnection(conn, config, &args)
		}
	}
}
