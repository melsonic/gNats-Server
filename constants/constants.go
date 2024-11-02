package constants

import "fmt"

const (
	PORT  string = "4222"
	CR    byte   = byte('\r')
	LF    byte   = byte('\n')
	SPACE byte   = byte(' ')

	/// commands
	CONNECT string = "CONNECT"
	PING    string = "PING"
	PONG    string = "PONG"
	PUB     string = "PUB"
	SUB     string = "SUB"
)

var (
	ADDRESS string = fmt.Sprintf(":%s", PORT)
	CRLF    []byte = []byte("\r\n")
)
