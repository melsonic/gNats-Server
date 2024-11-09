package util

import (
	"fmt"
	"net"
	"strings"

	"github.com/melsonic/gnats-server/constants"
)

func SendInitialInfoMessage(conn net.Conn) {
	var client_ip string = strings.Split(conn.RemoteAddr().String(), ":")[0]
	var initialResponse string = fmt.Sprintf("INFO {\"host\":\"0.0.0.0\",\"port\":%s,\"max_payload\":1048576,\"client_ip\":\"%s\"}\n", constants.PORT, client_ip)
	conn.Write([]byte(initialResponse))
}
