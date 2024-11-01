package main

import (
	"log"
	"net"

	constants "github.com/melsonic/gnats-server/constants"
  core "github.com/melsonic/gnats-server/core"
)

func main() {
	listener, err := net.Listen("tcp", constants.ADDRESS)
	defer listener.Close()
	if err != nil {
		log.Fatalln("Error creating a listener")
	}
	for {
		conn, err := listener.Accept()
		defer conn.Close()
		if err != nil {
			log.Fatalln("Error establishing a connection")
		}
    core.Handler(conn)
	}
}
