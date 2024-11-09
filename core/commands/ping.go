package commands

import (
	"net"

	"github.com/melsonic/gnats-server/constants"
)

func HandlePing(conn net.Conn) {
	conn.Write([]byte(constants.RESPONSE_PONG))
}
