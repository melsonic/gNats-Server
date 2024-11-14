package core

import (
	"errors"
	"strconv"

	"github.com/melsonic/gnats-server/constants"
	"github.com/melsonic/gnats-server/core/commands"
)

type parserState int

const (
	OP_START = iota
	OP_C
	OP_CO
	OP_CON
	OP_CONN
	OP_CONNE
	OP_CONNEC
	OP_CONNECT
	OP_CONNECT_ARG
	OP_P
	OP_PI
	OP_PIN
	OP_PING
	OP_PO
	OP_PON
	OP_PONG
	OP_S
	OP_SU
	OP_SUB
	OP_SUB_ARG_SUB
	OP_SUB_SUB_ID
	OP_PU
	OP_PUB
	OP_PUB_SUB
	OP_PUB_MSG_LEN
	OP_PUB_MSG
)

type pubConfig struct {
	subject string
	msgLen  int
	msg     []byte
}

type subConfig struct {
	subject string
	sid     string
}

type Parser struct {
	state    parserState
	subState subConfig
	pubState pubConfig
	arg      []byte
}

func (p *Parser) Reset() {
	p.state = OP_START
	p.subState = subConfig{}
	p.pubState = pubConfig{}
	p.arg = nil
}

func (p *Parser) Parse(config *serverConfig, buffer []byte) error {
	for _, b := range buffer {
		if b == constants.ZERO_BYTE {
			continue
		}
		switch p.state {
		case OP_START:
			p.Reset()
			switch b {
			case 'c', 'C':
				p.state = OP_C
			case 'p', 'P':
				p.state = OP_P
			case 's', 'S':
				p.state = OP_S
			default:
				// error
				goto parseError
			}
		case OP_C:
			switch b {
			case 'o', 'O':
				p.state = OP_CO
			default:
				// error
				goto parseError
			}
		case OP_CO:
			switch b {
			case 'n', 'N':
				p.state = OP_CON
			default:
				// error
				goto parseError
			}
		case OP_CON:
			switch b {
			case 'n', 'N':
				p.state = OP_CONN
			default:
				// error
				goto parseError
			}
		case OP_CONN:
			switch b {
			case 'e', 'E':
				p.state = OP_CONNE
			default:
				// error
				goto parseError
			}
		case OP_CONNE:
			switch b {
			case 'c', 'C':
				p.state = OP_CONNEC
			default:
				// error
				goto parseError
			}
		case OP_CONNEC:
			switch b {
			case 't', 'T':
				p.state = OP_CONNECT
			default:
				// error
				goto parseError
			}
		case OP_CONNECT:
			switch b {
			case ' ':
				p.state = OP_CONNECT_ARG
			default:
				// error
				goto parseError
			}
		case OP_CONNECT_ARG:
			switch b {
			case '\r':
				// skip
			case '\n':
				p.state = OP_START
				err := config.connectConf.SetUpConnect(p.arg)
				if err != nil {
					// error
					goto parseError
				}
				commands.ConnectHandler(config.connectConf.Verbose, config.channel)
			default:
				p.arg = append(p.arg, b)
			}
		case OP_P:
			switch b {
			case 'i', 'I':
				p.state = OP_PI
			case 'u', 'U':
				p.state = OP_PU
			default:
				// error
				goto parseError
			}
		case OP_S:
			switch b {
			case 'u', 'U':
				p.state = OP_SU
			default:
				// error
				goto parseError
			}
		case OP_PI:
			switch b {
			case 'n', 'N':
				p.state = OP_PIN
			default:
				// error
				goto parseError
			}
		case OP_PU:
			switch b {
			case 'b', 'B':
				p.state = OP_PUB
			default:
				// error
				goto parseError
			}
		case OP_SU:
			switch b {
			case 'b', 'B':
				p.state = OP_SUB
			default:
				// error
				goto parseError
			}
		case OP_PIN:
			switch b {
			case 'g', 'G':
				p.state = OP_PING
			default:
				// error
				goto parseError
			}
		case OP_PUB:
			switch b {
			case '\r':
				// skip
			case ' ', '\n':
				p.state = OP_PUB_SUB
			default:
				// error
				goto parseError
			}
		case OP_SUB:
			switch b {
			case '\r':
				// skip
			case ' ', '\n':
				p.state = OP_SUB_ARG_SUB
				p.arg = nil
			default:
				// error
				goto parseError
			}
		case OP_PING:
			switch b {
			case '\r':
				// skip
			case '\n':
				p.state = OP_START
				commands.PingHandler(config.channel)
			default:
				// error
				goto parseError
			}
		case OP_PUB_SUB:
			switch b {
			// PUB CC 6\R\NHI MOM\R\N
			case '\r':
				// skip
			case ' ', '\n':
				p.state = OP_PUB_MSG_LEN
				p.pubState.subject = string(p.arg)
				p.arg = nil
			default:
				p.arg = append(p.arg, b)
			}
		case OP_SUB_ARG_SUB:
			switch b {
			case '\r':
				// skip
			case ' ', '\n':
				p.state = OP_SUB_SUB_ID
				p.subState.subject = string(p.arg)
				p.arg = nil
			default:
				p.arg = append(p.arg, b)
			}
		case OP_PUB_MSG_LEN:
			switch b {
			case '\r':
				// skip
			case ' ', '\n':
				p.state = OP_PUB_MSG
				msgLen, err := strconv.Atoi(string(p.arg))
				if err != nil {
					// error
					goto parseError
				}
				p.pubState.msgLen = msgLen
				p.arg = nil
			default:
				p.arg = append(p.arg, b)
			}
		case OP_SUB_SUB_ID:
			switch b {
			case '\r':
				// skip
			case ' ', '\n':
				p.state = OP_START
				p.subState.sid = string(p.arg)
				commands.SubHandler(config.connectConf.Verbose, p.subState.subject, p.subState.sid, config.channel)
			default:
				p.arg = append(p.arg, b)
			}
		case OP_PUB_MSG:
			switch b {
			case '\r':
				// skip
			case ' ', '\n':
				if b == ' ' && len(p.arg) < p.pubState.msgLen {
					p.arg = append(p.arg, b)
					break
				}
				p.state = OP_START
				if len(p.arg) > p.pubState.msgLen {
					p.arg = p.arg[:p.pubState.msgLen]
				}
				p.pubState.msg = p.arg
				commands.PubHandler(config.connectConf.Verbose, p.pubState.subject, p.pubState.msgLen, p.pubState.msg, config.channel)
			default:
				p.arg = append(p.arg, b)
			}
		}
	}

	return nil
	// raise parse error label
parseError:
	return errors.New("Error Parsing Input\r\n")
}
