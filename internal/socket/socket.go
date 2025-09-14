package socket

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/kvnbanunu/uds-caesar-cipher/internal/options"
)

func Ssetup(a *options.Args) (net.Listener, error) {
	a.Log(true, "socket.Ssetup()\n\tSetting up Server...")
	// remove sockfile if already exists
	if _, err := os.Stat(a.Path); err == nil {
		a.Log(false, "\tSockfile already exists. Removing...")
		os.Remove(a.Path)
		a.Log(false, "\tStale sockfile removed.")
	}

	a.Log(false, fmt.Sprintf("\tCreating sockfile: %s", a.Path))
	sock, err := net.Listen("unix", a.Path)
	if err != nil {
		return nil, err
	}

	a.Log(false, fmt.Sprintf("\tNow listening on sockfile: %s", a.Path))

	return sock, nil
}

// Closes socket and removes the sockfile
func Cleanup(sock net.Listener, a *options.Args) {
	sock.Close()
	os.Remove(a.Path)
	a.Log(true, "socket.Cleanup()\n\tSockfile removed")
}

// Clean up the socket file on SIGTERM
func HandleSignal(sock net.Listener, a *options.Args) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("") // print new line after Ctrl-C
		a.Log(true, "socket.HandleSignal()\n\tReceived shutdown signal")
		Cleanup(sock, a)
		os.Exit(0)
	}()
}

// func HandleConnection(conn net.Conn, size int) {
// 	defer conn.Close()
// 	buf := make([]byte, size)
//
// 	for {
// 		nBytes, err := conn.Read(buf)
// 		if err != nil {
// 			return
// 		}
// 	}
// }
