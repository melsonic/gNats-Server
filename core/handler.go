package core

import (
	"fmt"
	"io"
	"net"

	"github.com/melsonic/gnats-server/constants"
	"github.com/melsonic/gnats-server/core/commands"
  "github.com/melsonic/gnats-server/util"
)

func Handler(conn net.Conn) {
	defer conn.Close()
  /// send info message
  util.SendInitialInfoMessage(conn)
  /// variables
	buffer := make([]byte, 100000)
	channel := make(chan string)
	success := true
	go func() {
		for {
			msg := <-channel
			conn.Write([]byte(msg))
		}
	}()
	for success {
		buffer = make([]byte, 100000)
		_, err := conn.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Println(err.Error())
			fmt.Println("Error reading from connection")
		}
		fmt.Printf("input -> ~%s~\n", string(buffer))
		tokens := Parse(buffer)
		if len(tokens) == 0 {
			conn.Close()
			return
		}
		switch tokens[0] {
		case constants.CONNECT:
			// commands.HandleConnect(conn)
		case constants.PING:
			commands.HandlePing(conn)
		case constants.SUB:
			success = commands.HandleSub(conn, tokens[1:], channel)
		case constants.PUB:
			success = commands.HandlePub(conn, tokens[1:])
		default:
			/// throw error && close the connection
			// conn.Close()
			// return
		}
	}
}
