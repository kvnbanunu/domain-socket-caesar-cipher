package main

import (
	"log"

	"github.com/kvnbanunu/uds-caesar-cipher/internal/options"
	"github.com/kvnbanunu/uds-caesar-cipher/internal/socket"
)

func main() {
	args := options.Args{}
	args.ParseArgs(true)
	
	config, err := args.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	conn, err := socket.Csetup(&args)
	if err != nil {
		log.Fatal("Error setting up client:", err)
	}
	
	socket.Request(conn, config, &args)
}
