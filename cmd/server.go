package main

import (
	"fmt"
	"log"

	"github.com/kvnbanunu/uds-caesar-cipher/internal/options"
)

func main() {
	args := options.Args{}

	args.ParseArgs(false)

	cfg, err := args.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
}
