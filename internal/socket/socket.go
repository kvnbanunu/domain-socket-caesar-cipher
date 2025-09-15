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

	a.Log(false, fmt.Sprintf("\tResolving Unix Address: %s", a.Path))
	addr, err := net.ResolveUnixAddr("unix", a.Path)
	if err != nil {
		return nil, err
	}

	a.Log(false, "\tCreating sockfile...")
	sock, err := net.ListenUnix("unix", addr)
	if err != nil {
		return nil, err
	}

	a.Log(false, fmt.Sprintf("\tNow listening on sockfile: %s", a.Path))

	return sock, nil
}

func CSetup(a *options.Args) (*net.UnixConn, error) {
	// Check if server running (sockfile exists)
	a.Log(true, `socket.Csetup()
	Setting up Client...
	Checking server sockfile...`)
	if _, err := os.Stat(a.Path); err != nil {
		return nil, err
	}

	a.Log(false, fmt.Sprintf("\tResolving Unix Address: %s", a.Path))
	addr, err := net.ResolveUnixAddr("unix", a.Path)
	if err != nil {
		return nil, err
	}

	a.Log(false, "\tConnecting to server...")
	conn, err := net.DialUnix("unix", nil, addr)
	if err != nil {
		return nil, err
	}

	a.Log(false, "\tConnection established")
	return conn, nil
}

// Closes socket and removes the sockfile
func Cleanup(sock *net.UnixListener, a *options.Args) {
	a.ExitFlag = true
	a.Log(true, "socket.Cleanup()\n\tClosing socket...")
	sock.Close()

	a.Log(false, "\tUnlinking sockfile...")
	os.Remove(a.Path)
}

// Clean up the socket file on SIGTERM
func HandleSignal(sock *net.UnixListener, a *options.Args) {
	a.Log(true, "socket.HandleSignal()\n\tNow listening for SIGTERM")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("") // print new line after Ctrl-C
		a.Log(true, "socket.HandleSignal()\n\tReceived shutdown signal\n\tInitiate cleanup...")
		Cleanup(sock, a)

		a.Log(true, "socket.HandleSignal()\n\tCleanup success\n\tServer shutting down.")
		os.Exit(0)
	}()
}

func HandleConnection(conn *net.UnixConn, config *options.Config, a *options.Args) error {
	defer conn.Close()

	a.Log(true, `socket.HandleConnection()
	Connection established...
	Reading from socket...`)

	buf := make([]byte, config.BufferSize)
	n, err := conn.Read(buf)
	if err != nil {
		return err
	}

	str := fmt.Sprintf("%s", buf[:n]) // upto 'n' bytes

	a.Log(false, fmt.Sprintf("\tBytes read: %d\n\tMessage Received: %s\n\tDecoding...", n, str))
	msg, err := decode(str)
	if err != nil {
		return err
	}

	a.Log(false, fmt.Sprintf("\tContent: %s\n\tShift Value: %d", msg.Content, msg.Shift))
	msg.Content = caesar.Process(msg, "cipher", config.CipherLimit)

	a.Log(false, fmt.Sprintf("\tApplying Caesar Cipher...\n\tEncrypted Message: %s", msg.Content))
	encoded := []byte(encode(msg))

	a.Log(false, fmt.Sprintf("\tEncoding response...\n\tSending response: %s", encoded))
	conn.Write(encoded)

	a.Log(false, "\tClosing connection...")

	return nil
}

func Request(conn *net.UnixConn, config *options.Config, a *options.Args) error {
	defer conn.Close()

	a.Log(true, "socket.Request()\n\tEncoding message...")
	buf := make([]byte, config.BufferSize)
	encoded := []byte(encode(a.Message))

	a.Log(false, fmt.Sprintf("\tEncoded message: %s\n\tSending to server...", encoded))
	conn.Write(encoded)

	a.Log(false, "\tWaiting for response from server...")
	n, err := conn.Read(buf)
	if err != nil {
		return err
	}

	str := fmt.Sprintf("%s", buf[:n])
	a.Log(false, fmt.Sprintf("\tResponse from server: %s\n\tDecoding...", str))
	a.Message, err = decode(str)
	if err != nil {
		return err
	}

	fmt.Println("\n\tEncrypted response:\t", a.Message.Content)
	a.Message.Content = caesar.Process(a.Message, "decipher", config.CipherLimit)
	fmt.Println("\tDecrypted message:\t", a.Message.Content)

	return nil
}

// encode the message into a string with the following format
// <Shift Value>#<Message Content>
// ex. Shift 6, Message "Hello" = "6#Hello"
func encode(msg options.Message) string {
	return fmt.Sprintf("%d#%s", msg.Shift, msg.Content)
}

func decode(str string) (options.Message, error) {
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
	shift, err := strconv.Atoi(shiftstr)
	if err != nil {
		return res, err
	}

	res.Shift = shift
	res.Content = str[nextIndex:] // the remaining string after '#'
	return res, nil
}
