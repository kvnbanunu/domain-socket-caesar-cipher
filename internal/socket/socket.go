package socket

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/kvnbanunu/uds-caesar-cipher/internal/options"
)

func Ssetup(a *options.Args) (net.Listener, error) {
	// remove sockfile if already exists
	if _, err := os.Stat(a.Path); err == nil {
		os.Remove(a.Path)
	}

	sock, err := net.Listen("unix", a.Path)
	if err != nil {
		return _, err
	}

	return sock, nil
}

func Cleanup(sock net.Listener, a *options.Args) {
	sock.Close()
	os.Remove(a.Path)
	a.Log(true, "socket.Cleanup()\n\tSockfile removed")
}

// Clean up the socket file on SIGTERM
func HandleSignal(a *options.Args) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		a.Log(true, "socket.HandleSignal()\n\tReceived shutdown signal")
		os.Exit(1)
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
