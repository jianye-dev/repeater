package repeater

import (
	"io"
	"net"
)

type vncserver struct {
	ln net.Listener
}

func (svc *vncserver) ListenAndServe(addr string, r chan<- io.ReadWriteCloser) (err error) {
	var conn net.Conn

	if svc.ln, err = newVNCListener(addr); err != nil {
		return
	}

	for {
		if conn, err = svc.ln.Accept(); err != nil {
			break
		}
		r <- conn
	}

	return
}

func (svc *vncserver) Shutdown() error {
	if svc.ln != nil {
		return svc.ln.Close()
	}
	return nil
}
