package httputil2

import "net"

func PickPort() int {
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		return -1
	}
	defer listener.Close()

	addr := listener.Addr().(*net.TCPAddr)
	return addr.Port
}
