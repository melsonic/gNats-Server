package commands

import (
	"github.com/melsonic/gnats-server/constants"
	"github.com/melsonic/gnats-server/data"
)

func UnsubHandler(verbose bool, sid int, channel chan string) {
	if verbose {
		channel <- constants.RESPONSE_OK
	}
	go data.GSubjectSIDs.Unsub(sid, channel)
}
