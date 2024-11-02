package commands

import "net"

func HandlePing(conn net.Conn) {
	var response string = "PONG\r\n"
	conn.Write([]byte(response))
}
