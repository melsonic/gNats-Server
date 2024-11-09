package commands

import (
	"net"

	"github.com/melsonic/gnats-server/constants"
)

func HandleConnect(conn net.Conn) {
	conn.Write([]byte(constants.RESPONSE_PONG))
}
