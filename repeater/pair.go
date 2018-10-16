package repeater

import (
	"io"
)

const (
	readBufferSize = 32768
	bufChanSize    = 32
)

type pair struct {
	token    string
	isServer bool
	rwc      io.ReadWriteCloser
	peer     io.ReadWriteCloser
	bufc     chan []byte
	errc     chan error
	donec    chan<- string
}

func (p *pair) prepare() (err error) {

	if !p.isServer {
		if err = sendRepeaterVersion(p.rwc); err != nil {
			return
		}
	}

	p.token, err = fetchHostInfo(p.rwc)

	return
}

func (p *pair) localServe() {

	buf := make([]byte, readBufferSize)

	for {
		if nr, err := p.rwc.Read(buf); err != nil {
			p.errc <- err
			p.donec <- p.token
			return
		} else {
			rbuf := make([]byte, nr)
			copy(rbuf, buf[:nr])
			p.bufc <- rbuf
		}
	}
}

func (p *pair) peerServe() {
	for buf := range p.bufc {
		_, err := p.peer.Write(buf)
		if err != nil {
			p.errc <- err
			return
		}
	}
}

func (p *pair) serve() {

	go p.peerServe()

	go func() {
		_, err := io.Copy(p.rwc, p.peer)
		p.errc <- err
	}()

	<-p.errc

	p.rwc.Close()
	p.peer.Close()
}

func newPair(conn io.ReadWriteCloser, isServer bool, donec chan<- string) (item *pair, err error) {

	item = &pair{
		isServer: isServer,
		rwc:      conn,
		bufc:     make(chan []byte, bufChanSize),
		errc:     make(chan error, 1),
		donec:    donec,
	}

	err = item.prepare()

	return
}
