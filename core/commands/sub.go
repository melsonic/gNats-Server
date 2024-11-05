package commands

import (
	"net"
	"strconv"

	data "github.com/melsonic/gnats-server/data"
)

func HandleSub(conn net.Conn, args []string, channel chan string) bool {
	if len(args) < 2 {
		return false
	}
	subject := args[0]
	sid, err := strconv.Atoi(args[1])
	if err != nil {
		return false
	}
	success := data.GSubjectSIDs.Add(subject, sid, channel)
	var response string
	if success {
		response = "+OK\r\n"
	} else {
		response = "\r\n"
	}
	conn.Write([]byte(response))
	return true
}
