package commands

import (
	"github.com/melsonic/gnats-server/constants"
)

func ConnectHandler(verbose bool, channel chan string) {
	if verbose {
		channel <- constants.RESPONSE_OK
	}
}
