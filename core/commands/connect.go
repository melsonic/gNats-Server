package commands

import "net"

func HandleConnect(conn net.Conn) {
	var response string = "+OK\r\n"
	conn.Write([]byte(response))
}
