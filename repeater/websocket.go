package repeater

import (
	"context"
	"io"
	"net/http"

	"github.com/gorilla/websocket"
)

// nolint
type websock struct {
	svr *http.Server
	r   chan<- io.ReadWriteCloser
}

var upgrader = websocket.Upgrader{
	CheckOrigin:  func(r *http.Request) bool { return true },
	Subprotocols: []string{"binary"},
}

// nolint
func (svc *websock) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		return
	}

	svc.r <- &wsConn{Conn: conn}
}

// nolint
func (svc *websock) ListenAndServe(addr string, r chan<- io.ReadWriteCloser) error {

	svc.svr = &http.Server{Addr: addr, Handler: svc}
	svc.r = r

	return svc.svr.ListenAndServe()
}

// nolint
func (svc *websock) Shutdown() error {
	if svc.svr != nil {
		return svc.svr.Shutdown(context.Background())
	}
	return nil
}
