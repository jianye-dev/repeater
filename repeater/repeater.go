package repeater

import (
	"io"
	"sync"
)

type Repeater struct {
	mu         sync.Mutex
	activeConn map[string]*pair
	doneChan   chan struct{}
	tokenDone  chan string
	chServer   chan io.ReadWriteCloser
	chViewer   chan io.ReadWriteCloser
}

func (r *Repeater) onNewConnection(item *pair) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.activeConn == nil {
		r.activeConn = make(map[string]*pair)
	}

	if _item, loaded := r.activeConn[item.token]; loaded {

		if _item.peer != nil {
			item.rwc.Close()
			return
		}

		if _item.isServer == item.isServer {
			_item.rwc.Close()
			r.activeConn[item.token] = item
			return
		}

		_item.peer = item.rwc

		go _item.serve()
	} else {
		r.activeConn[item.token] = item
		go item.localServe()
	}
}

func (r *Repeater) onCloseConnection(token string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.activeConn == nil {
		return
	}

	delete(r.activeConn, token)
}

func (r *Repeater) onNewServer(conn io.ReadWriteCloser) {
	if item, err := newPair(conn, true, r.tokenDone); err != nil {
		conn.Close()
	} else {
		r.onNewConnection(item)
	}
}

func (r *Repeater) onNewViewer(conn io.ReadWriteCloser) {
	if item, err := newPair(conn, false, r.tokenDone); err != nil {
		conn.Close()
	} else {
		r.onNewConnection(item)
	}
}

// nolint
func (r *Repeater) ListenAndServe(svrAddr, wsAddr, vwAddr string) (err error) {

	r.doneChan = make(chan struct{})
	r.tokenDone = make(chan string)

	r.chServer = make(chan io.ReadWriteCloser)
	r.chViewer = make(chan io.ReadWriteCloser)

	server := &vncserver{}
	websock := &websock{}
	viewer := &vncserver{}

	if svrAddr != "" {
		go server.ListenAndServe(svrAddr, r.chServer)
	}

	if wsAddr != "" {
		go websock.ListenAndServe(wsAddr, r.chViewer)
	}

	if vwAddr != "" {
		go viewer.ListenAndServe(vwAddr, r.chViewer)
	}

loop:
	for {
		select {
		case conn := <-r.chServer:
			r.onNewServer(conn)
		case conn := <-r.chViewer:
			r.onNewViewer(conn)
		case token := <-r.tokenDone:
			r.onCloseConnection(token)
		case <-r.doneChan:
			break loop
		}
	}

	server.Shutdown()
	viewer.Shutdown()
	websock.Shutdown()

	return
}

// nolint
func (r *Repeater) Shutdown() {
	close(r.doneChan)
loop:
	for {
		select {
		case token := <-r.tokenDone:
			r.onCloseConnection(token)
		default:
			r.mu.Lock()

			if len(r.activeConn) == 0 {
				r.mu.Unlock()
				break loop
			}

			for _, v := range r.activeConn {
				v.rwc.Close()
			}
			r.mu.Unlock()
		}
	}

	close(r.chServer)
	close(r.chViewer)
	close(r.tokenDone)
}

// nolint
func New() *Repeater {
	return &Repeater{}
}
