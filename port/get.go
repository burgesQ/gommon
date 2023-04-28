package port

import (
	"errors"
	"net"
)

var ErrGetTCP = errors.New("get tcp address")

func GetFree() (port int, err error) {
	var a *net.TCPAddr

	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()

			addr, ok := l.Addr().(*net.TCPAddr)
			if !ok {
				return 0, ErrGetTCP
			}

			return addr.Port, nil
		}
	}

	return
}
