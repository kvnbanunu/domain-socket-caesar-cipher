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
		log.Fatal(err)
	}

	conn, err := socket.CSetup(&args)
	if err != nil {
		log.Fatal(err)
	}
	
	socket.Request(conn, config, &args)
}
