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
	OP_SUB_SUBJECT
	OP_SUB_SUBJECT_ID
	OP_PU
	OP_PUB
	OP_PUB_SUB
	OP_PUB_MSG
	OP_PUB_MSG_LEN
	OP_U
	OP_UN
	OP_UNS
	OP_UNSU
	OP_UNSUB
	OP_UNSUB_SID
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
			case 'u', 'U':
				p.state = OP_U
			default:
				goto parseError
			}
		case OP_C:
			switch b {
			case 'o', 'O':
				p.state = OP_CO
			default:
				goto parseError
			}
		case OP_CO:
			switch b {
			case 'n', 'N':
				p.state = OP_CON
			default:
				goto parseError
			}
		case OP_CON:
			switch b {
			case 'n', 'N':
				p.state = OP_CONN
			default:
				goto parseError
			}
		case OP_CONN:
			switch b {
			case 'e', 'E':
				p.state = OP_CONNE
			default:
				goto parseError
			}
		case OP_CONNE:
			switch b {
			case 'c', 'C':
				p.state = OP_CONNEC
			default:
				goto parseError
			}
		case OP_CONNEC:
			switch b {
			case 't', 'T':
				p.state = OP_CONNECT
			default:
				goto parseError
			}
		case OP_CONNECT:
			switch b {
			case ' ':
				p.state = OP_CONNECT_ARG
			default:
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
				goto parseError
			}
		case OP_S:
			switch b {
			case 'u', 'U':
				p.state = OP_SU
			default:
				goto parseError
			}
		case OP_PI:
			switch b {
			case 'n', 'N':
				p.state = OP_PIN
			default:
				goto parseError
			}
		case OP_PU:
			switch b {
			case 'b', 'B':
				p.state = OP_PUB
			default:
				goto parseError
			}
		case OP_SU:
			switch b {
			case 'b', 'B':
				p.state = OP_SUB
			default:
				goto parseError
			}
		case OP_PIN:
			switch b {
			case 'g', 'G':
				p.state = OP_PING
			default:
				goto parseError
			}
		case OP_PUB:
			switch b {
			case '\r':
				// skip
			case ' ', '\n':
				p.state = OP_PUB_SUB
			default:
				goto parseError
			}
		case OP_SUB:
			switch b {
			case '\r':
				// skip
			case ' ', '\n':
				p.state = OP_SUB_SUBJECT
				p.arg = nil
			default:
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
		case OP_SUB_SUBJECT:
			switch b {
			case '\r':
				// skip
			case ' ', '\n':
				p.state = OP_SUB_SUBJECT_ID
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
					goto parseError
				}
				p.pubState.msgLen = msgLen
				p.arg = nil
			default:
				p.arg = append(p.arg, b)
			}
		case OP_SUB_SUBJECT_ID:
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
		case OP_U:
			switch b {
			case 'n', 'N':
				p.state = OP_UN
			default:
				goto parseError
			}
		case OP_UN:
			switch b {
			case 's', 'S':
				p.state = OP_UNS
			default:
				goto parseError
			}
		case OP_UNS:
			switch b {
			case 'u', 'U':
				p.state = OP_UNSU
			default:
				goto parseError
			}
		case OP_UNSU:
			switch b {
			case 'b', 'B':
				p.state = OP_UNSUB
			default:
				goto parseError
			}
		case OP_UNSUB:
			switch b {
			case '\r':
				// skip
			case ' ', '\n':
				p.state = OP_UNSUB_SID
			default:
				goto parseError
			}
		case OP_UNSUB_SID:
			switch b {
			case '\r':
			// skip
			case ' ', '\n':
				p.state = OP_START
				sid, err := strconv.Atoi(string(p.arg))
				if err != nil {
					goto parseError
				}
				commands.UnsubHandler(config.connectConf.Verbose, sid, config.channel)
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
