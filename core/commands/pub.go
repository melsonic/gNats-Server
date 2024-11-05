package commands

import (
	"net"
	"strconv"

	"github.com/melsonic/gnats-server/data"
)

func HandlePub(conn net.Conn, args []string) bool {
	if len(args) < 2 {
		return false
	}
	subject := args[0]
	bytesLen, err := strconv.Atoi(args[1])
	if err != nil {
		return false
	}
	inputMsg := make([]byte, bytesLen+2) // adding 2 for CRLF bytes
	_, err = conn.Read(inputMsg)
	if err != nil {
		return false
	}
	// to remove 2 bytes for CRLF bytes
	go data.GSubjectSIDs.Publish(subject, inputMsg[:len(inputMsg)-2])
	/// TODO: connection getting closed after publishing
	return true
}
