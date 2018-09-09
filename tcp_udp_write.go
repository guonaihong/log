package log

import (
	"io"
	"net"
	"strings"
)

func ParseSocket(url string) (io.Writer, error) {
	if strings.HasPrefix(url, "udp://") {
		addr := url[len("udp://"):]
		raddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			return nil, err
		}

		return net.DialUDP("udp", nil, raddr)
	}

	if strings.HasPrefix(url, "tcp://") {
		return net.Dial("tcp", url[len("tcp://"):])
	}

	return net.Dial("tcp", url)
}
