package util

import "fmt"

func PrintInputData(input []byte) {
	var output string
	for _, b := range input {
		if b == 0 {
			break
		}
		output += string(b)
	}
	fmt.Println(output)
}

func ResetBuffer(input []byte) {
	for i := range input {
    input[i] = 0
	}
}
