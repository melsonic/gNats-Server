package commands

import (
	"net"
	"strconv"

	data "github.com/melsonic/gnats-server/data"
)

func HandleSub(conn net.Conn, args []string) bool {
	if len(args) < 2 {
		return false
	}
	subject := args[0]
	sid, err := strconv.Atoi(args[1])
	if err != nil {
		return false
	}
	data.GSubjectSIDs.Add(subject, sid)
	var response string = "OK+\r\n"
	conn.Write([]byte(response))
	return true
}
