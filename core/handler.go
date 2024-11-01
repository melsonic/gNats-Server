package core

import (
	"fmt"
	"net"
  util "github.com/melsonic/gnats-server/util"
)

func Handler(conn net.Conn) {
	buffer := make([]byte, 1024)
	for {
    util.ResetBuffer(buffer)
		_, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from connection")
		}
    // fmt.Println(buffer)
		util.PrintInputData(buffer)
	}
}
