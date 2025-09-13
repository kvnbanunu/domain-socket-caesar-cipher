package options

// Package options implements utility functions for 
// parsing and storing command line arguments
// logging debug statements

import (
	"flag"
	"fmt"
)

// Args is used to store command-line arguments
type Args struct {
	Path string
	Debug bool
	Message Message
}

// Store a message along with the cipher shift value
type Message struct {
	Content string
	Shift int
}

// Parse command line arguments for either client or server
func (a *Args) ParseArgs(isClient bool) {
	path := flag.String("p", "bin/domain.sock", "Path to domain socket")
	debug := flag.Bool("d", false, "Run program with debug log statements")
	var content *string
	var shift *int

	if (isClient) {
		content = flag.String("i", "Hello, World", "String message to be encrypted")
		shift = flag.Int("s", 3, "Shift value for Caesar Cipher")
	}

	flag.Parse()

	a.Path = *path
	a.Debug = *debug
	
	if (isClient) {
		a.Message.Content = *content
		a.Message.Shift = *shift
	}
}

// This method prints debug statements to stdout if debug flag is enabled
func (a Args) Log(msg string) {
	if (a.Debug) {
		fmt.Println(msg)
	}
}
