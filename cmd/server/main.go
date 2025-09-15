package main

import (
	"fmt"
	"log"
	"time"

	"github.com/kvnbanunu/uds-caesar-cipher/internal/options"
	"github.com/kvnbanunu/uds-caesar-cipher/internal/socket"
)

func main() {
	args := options.Args{}
	args.ParseArgs(false)

	config, err := args.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	sock, err := socket.Ssetup(&args)
	if err != nil {
		log.Fatal("Error setting up server:", err)
	}

	defer socket.Cleanup(sock, &args)

	socket.HandleSignal(sock, &args)

	for {
		conn, err := sock.AcceptUnix()
		if err != nil {
			if args.ExitFlag {
				time.Sleep(time.Millisecond)
				continue // server shutting down, wait for cleanup to finish
			}

			fmt.Println("Error accepting connection:", err)
			break
		}
		err = socket.HandleConnection(conn, config, &args)
		if err != nil {
			fmt.Println("Error handling connection:", err)
		}
	}
}
