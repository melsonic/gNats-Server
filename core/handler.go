package core

import (
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/melsonic/gnats-server/constants"
	"github.com/melsonic/gnats-server/core/commands"
	"github.com/melsonic/gnats-server/util"
)

func Handler(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 4096)
	var responseString string = util.BuildInitialResponseString(strings.Split(conn.RemoteAddr().String(), ":")[0])
	var initialResponse []byte = []byte(responseString)
	conn.Write(initialResponse)
	channel := make(chan string)
	success := true
	go func() {
		for {
			msg := <-channel
			conn.Write([]byte(msg))
		}
	}()
	for {
		util.ResetBuffer(buffer)
		_, err := conn.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Println(err.Error())
			fmt.Println("Error reading from connection")
		}
		tokens := Parse(buffer)
		if len(tokens) == 0 {
			conn.Close()
			return
		}
		switch tokens[0] {
		case constants.CONNECT:
			fmt.Println("connect")
			commands.HandleConnect(conn)
		case constants.PING:
			fmt.Println("ping")
			commands.HandlePing(conn)
		case constants.SUB:
			fmt.Println("sub")
			success = commands.HandleSub(conn, tokens[1:], channel)
		case constants.PUB:
			fmt.Println("pub")
			success = commands.HandlePub(conn, tokens[1:])
		default:
			/// throw error && close the connection
			conn.Close()
			return
		}

		if !success {
			conn.Close()
			return
		}
	}
}
