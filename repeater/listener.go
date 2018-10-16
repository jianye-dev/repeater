package repeater

import "net"

type vncListener struct {
	*net.TCPListener
}

func (ln vncListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	tc.SetNoDelay(true)
	tc.SetKeepAlive(true)
	return tc, nil
}

func newVNCListener(addr string) (l net.Listener, err error) {

	if l, err = net.Listen("tcp", addr); err != nil {
		return
	}

	l = &vncListener{l.(*net.TCPListener)}

	return
}
