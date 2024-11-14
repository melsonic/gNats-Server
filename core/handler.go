package core

import (
	"fmt"
	"io"
	"net"

	"github.com/melsonic/gnats-server/util"
)

func Handler(conn net.Conn) {
	defer conn.Close()
	/// send info message
	util.SendInitialInfoMessage(conn)
	/// variables
	config := NewServerConfig(conn)
	go func() {
		for {
			msg := <-config.channel
			conn.Write([]byte(msg))
		}
	}()

	buffer := make([]byte, 4096)
	parser := Parser{}

	for {
		util.ResetBuffer(buffer)
		_, err := config.conn.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Println("Error reading from connection")
		}
		if err == io.EOF {
			break
		}
		err = parser.Parse(&config, buffer)
		if err != nil {
			conn.Write([]byte(err.Error()))
			break
		}
	}
}
