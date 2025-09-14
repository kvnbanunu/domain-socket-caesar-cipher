package socket

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/kvnbanunu/uds-caesar-cipher/internal/caesar"
	"github.com/kvnbanunu/uds-caesar-cipher/internal/options"
)

func Ssetup(a *options.Args) (*net.UnixListener, error) {
	a.Log(true, "socket.Ssetup()\n\tSetting up Server...")
	// remove sockfile if already exists
	if _, err := os.Stat(a.Path); err == nil {
		a.Log(false, "\tSockfile already exists. Removing...")
		os.Remove(a.Path)
		a.Log(false, "\tStale sockfile removed.")
	}

	addr, err := net.ResolveUnixAddr("unix", a.Path)
	if err != nil {
		fmt.Println("Error resolving Unix address:", err)
		return nil, err
	}

	a.Log(false, fmt.Sprintf("\tCreating sockfile: %s", a.Path))
	sock, err := net.ListenUnix("unix", addr)
	if err != nil {
		fmt.Println("Error creating Unix Socket", err)
		return nil, err
	}
	
	a.Log(false, fmt.Sprintf("\tNow listening on sockfile: %s", a.Path))

	return sock, nil
}

func CSetup(a *options.Args) (*net.UnixConn, error) {
	// Check if server running (sockfile exists)
	if _, err := os.Stat(a.Path); err != nil {
		fmt.Println("Server not available:", err)
		return nil, err
	}

	addr, err := net.ResolveUnixAddr("unix", a.Path)
	if err != nil {
		fmt.Println("Error resolving Unix address:", err)
		return nil, err
	}

	conn, err := net.DialUnix("unix", nil, addr)
	if err != nil {
		fmt.Println("Error connecting to Unix Socket:", err)
		return nil, err
	}

	return conn, nil
}

// Closes socket and removes the sockfile
func Cleanup(sock *net.UnixListener, a *options.Args) {
	sock.Close()
	os.Remove(a.Path)
	a.Log(true, "socket.Cleanup()\n\tSockfile removed")
}

// Clean up the socket file on SIGTERM
func HandleSignal(sock *net.UnixListener, a *options.Args) {
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

// encode the message into a string with the following format
// <Shift Value>#<Message Content>
// ex. Shift 6, Message "Hello" = "6#Hello"
func encode(msg options.Message) string {
	return fmt.Sprintf("%d#%s", msg.Shift, msg.Content)
}

func decode(str string) options.Message {
	res := options.Message{}
	nextIndex := 0
	shiftstr := ""
	for _, val := range str {
		nextIndex++
		if val == '#' {
			break
		}
		shiftstr += string(val)
	}
	shift, _ := strconv.Atoi(shiftstr)
	res.Shift = shift
	res.Content = str[nextIndex:] // the remaining string after '#'
	return res
}

func HandleConnection(conn *net.UnixConn, size int, limit int) {
	defer conn.Close()
	buf := make([]byte, size)
	n, err := conn.Read(buf)
	if err != nil {
		return
	}
	str := fmt.Sprintf("%s", buf[:n]) // upto 'n' bytes
	msg := decode(str)
	msg.Content = caesar.Process(msg, "cipher", limit)
	encoded := []byte(encode(msg))
	conn.Write(encoded)
}

func Request(conn *net.UnixConn, msg options.Message, size int, limit int) {
	defer conn.Close()
	buf := make([]byte, size)
	encoded := []byte(encode(msg))
	conn.Write(encoded)
	n, err := conn.Read(buf)
	if err != nil {
		return
	}
	str := fmt.Sprintf("%s", buf[:n])
	msg = decode(str)
	fmt.Println("Encrypted:", msg.Content)
	msg.Content = caesar.Process(msg, "decipher", limit)
	fmt.Println("Decrypted:", msg.Content)
}
