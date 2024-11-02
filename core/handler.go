package core

import (
	"fmt"
	"net"
	"strings"

	"github.com/melsonic/gnats-server/constants"
	"github.com/melsonic/gnats-server/core/commands"
	"github.com/melsonic/gnats-server/util"
)

func Handler(conn net.Conn) {
	buffer := make([]byte, 1024)
	var responseString string = util.BuildInitialResponseString(strings.Split(conn.RemoteAddr().String(), ":")[0])
	var initialResponse []byte = []byte(responseString)
	conn.Write(initialResponse)
	/// first command should be CONNECT {} from client
	var flag bool = false
	for {
		util.ResetBuffer(buffer)
		_, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from connection")
		}
		tokens := Parse(buffer)
		/// handle initial CONNECT command
		if !flag {
			if tokens[0] == constants.CONNECT {
				commands.HandleConnect(conn)
				flag = true
			} else {
				/// throw the error && close the connection
				conn.Close()
				return
			}
			continue
		}
		/// handle rest of the commands
		switch tokens[0] {
		case constants.PING:
			commands.HandlePing(conn)
		default:
			/// throw error && close the connection
			conn.Close()
			return
		}
	}
}
