package core

// CONNECT
// PING
// PONG
// SUB FOO

// PUB CC 7
// ronaldo
import (
	"fmt"

	"github.com/melsonic/gnats-server/constants"
)

func Parse(input []byte) []string {
	var tokens []string
	var current string
	for i := 0; i < len(input)-1; i++ {
		if input[i] == constants.CR && input[i+1] == constants.LF {
			tokens = append(tokens, current)
      current = ""
		}
		if input[i] == constants.SPACE {
			tokens = append(tokens, current)
			current = ""
		} else {
			current += string(input[i])
		}
	}
	return tokens
}
