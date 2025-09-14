package options

// Package options implements utility functions for
// parsing and storing command line arguments
// logging debug statements

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

// Config to be read during runtime, holds limit for buffer size
type Config struct {
	BufferSize  int `json:"BufferSize"`
	CipherLimit int `json:"CipherLimit"`
}

// Args is used to store command-line arguments
type Args struct {
	Path    string
	Debug   bool
	Message Message
}

// Store a message along with the cipher shift value
type Message struct {
	Content string
	Shift   int
}

// This method prints debug statements to stdout if debug flag is enabled
func (a Args) Log(timestamp bool, msg string) {
	if a.Debug {
		if timestamp {
			// log pkg Println includes date + time
			log.Println(msg)
		} else {
			fmt.Println(msg)
		}
	}
}

// Load config into memory
func (a *Args) LoadConfig() (*Config, error) {
	a.Log(true, "options.LoadConfig()\n\tLoading config...")

	file, err := os.ReadFile("config.json")
	if err != nil {
		return nil, err
	}

	// Convert json to struct
	var cfg Config
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		return nil, err
	}

	a.Log(false, "\tConfig Loaded")

	return &cfg, nil
}

// Parse command line arguments for either client or server
func (a *Args) ParseArgs(isClient bool) {
	// Set universal flags
	path := flag.String("p", "domain.sock", "Path to domain socket")
	debug := flag.Bool("d", false, "Run program with debug log statements")
	var content *string
	var shift *int

	// Client only flags
	if isClient {
		content = flag.String("i", "Hello, World", "String message to be encrypted")
		shift = flag.Int("s", 3, "Shift value for Caesar Cipher")
	}

	flag.Parse()

	a.Path = *path
	a.Debug = *debug

	// Format debug log
	logState := fmt.Sprintf(
		`options.ParseArgs()
	Parsing Command Line Arguments:
	Path: %s
	Debug: %t`,
		a.Path, a.Debug)

	// Only prints if debug = true
	a.Log(true, logState)

	if isClient {
		a.Message.Content = *content
		a.Message.Shift = *shift

		logState = fmt.Sprintf(
			"\tMessage: %s\n\tShift Value: %d",
			a.Message.Content, a.Message.Shift)
		a.Log(false, logState)
	}
}
