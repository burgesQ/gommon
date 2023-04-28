package port

import (
	"errors"
	"fmt"
	"net"
)

var (
	ErrGetTCP = errors.New("get tcp address")
	ErrNotIP  = errors.New("provided value isn't a valid ip")
)

// GetFree return an available free port to listen on.
// It use the synthax 'localhost:0' under the hood to attempt to listen on a
// random port. The list ip, which default to localhost, can be parametrized.
func GetFree(ip ...string) (port int, err error) {
	var (
		a   *net.TCPAddr
		lip = "localhost"
	)

	if len(ip) > 0 {
		if net.ParseIP(ip[0]) == nil {
			return 0, fmt.Errorf("%q: %w", ip[0], ErrNotIP)
		}

		lip = ip[0]
	}

	if a, err = net.ResolveTCPAddr("tcp", lip+":0"); err == nil {
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
