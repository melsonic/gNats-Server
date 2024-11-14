package core

import (
	"encoding/json"
	"net"
)

type serverConfig struct {
	conn        net.Conn
	channel     chan string
	connectConf connectConfig
}

func NewServerConfig(conn net.Conn) serverConfig {
	return serverConfig{
		conn:        conn,
		channel:     make(chan string),
		connectConf: connectConfig{},
	}
}

type connectConfig struct {
	Verbose bool   `json:"verbose"`
	Name    string `json:"name"`
	Lang    string `json:"lang"`
	Version string `json:"version"`
}

func (cc *connectConfig) SetUpConnect(data []byte) error {
	err := json.Unmarshal(data, cc)
	return err
}
