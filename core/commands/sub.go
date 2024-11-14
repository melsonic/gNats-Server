package commands

import (
	"strconv"

	"github.com/melsonic/gnats-server/constants"
	data "github.com/melsonic/gnats-server/data"
)

func SubHandler(verbose bool, subject string, subjectid string, channel chan string) bool {
	sid, err := strconv.Atoi(subjectid)
	if err != nil {
		return false
	}
	go data.GSubjectSIDs.Add(subject, sid, channel)
	if verbose {
		channel <- constants.RESPONSE_OK
	}
	return true
}
