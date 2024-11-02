package core

import (
	"fmt"
	"net"
	"strings"

	"github.com/melsonic/gnats-server/util"
)

func Handler(conn net.Conn) {
	buffer := make([]byte, 1024)
	var responseString string = util.BuildInitialResponseString(strings.Split(conn.RemoteAddr().String(), ":")[0])
	var initialResponse []byte = []byte(responseString)
	conn.Write(initialResponse)
	for {
		util.ResetBuffer(buffer)
		_, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from connection")
		}
		tokens := Parse(buffer)
	}
}
