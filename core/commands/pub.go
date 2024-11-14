package commands

import (
	"github.com/melsonic/gnats-server/constants"
	"github.com/melsonic/gnats-server/data"
)

func PubHandler(verbose bool, subject string, msgLen int, msg []byte, channel chan string) {
	if verbose {
		channel <- constants.RESPONSE_OK
	}
	go data.GSubjectSIDs.Publish(subject, msgLen, msg)
}
