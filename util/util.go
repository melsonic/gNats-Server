package util

import (
	"fmt"

	"github.com/melsonic/gnats-server/constants"
)

func PrintInputData(input []byte) {
	var output string
	for _, b := range input {
		// if b == 0 {
		// 	break
		// }
		fmt.Printf("%b => %c\n", b, b)
		output += string(b)
	}
	fmt.Println(output)
}

func ResetBuffer(input []byte) {
	for i := range input {
		input[i] = 0
	}
}

func BuildInitialResponseString(client_ip string) string {
	return fmt.Sprintf("INFO {\"host\":\"0.0.0.0\",\"port\":%s,\"client_ip\":\"%s\"}\n", constants.PORT, client_ip)
}
