package redir

import (
	"net"

	"github.com/Dreamacro/clash/adapters/inbound"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	"github.com/Dreamacro/clash/tunnel"
)

var (
	tun = tunnel.Instance()
)

type redirListener struct {
	net.Listener
	address string
	closed  bool
}

func NewRedirProxy(addr string) (*redirListener, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	rl := &redirListener{l, addr, false}

	go func() {
		log.Infoln("Redir proxy listening at: %s", addr)
		for {
			c, err := l.Accept()
			if err != nil {
				if rl.closed {
					break
				}
				continue
			}
			go handleRedir(c)
		}
	}()

	return rl, nil
}

func (l *redirListener) Close() {
	l.closed = true
	l.Listener.Close()
}

func (l *redirListener) Address() string {
	return l.address
}

func handleRedir(conn net.Conn) {
	target, err := parserPacket(conn)
	if err != nil {
		conn.Close()
		return
	}
	conn.(*net.TCPConn).SetKeepAlive(true)
	tun.Add(adapters.NewSocket(target, conn, C.REDIR))
}
