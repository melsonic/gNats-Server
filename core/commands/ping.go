package commands

import (
	"github.com/melsonic/gnats-server/constants"
)

func PingHandler(channel chan string) {
	channel <- constants.RESPONSE_PONG
}
