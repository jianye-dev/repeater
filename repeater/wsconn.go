package repeater

import (
	"io"

	"github.com/gorilla/websocket"
)

type wsConn struct {
	Conn *websocket.Conn
}

func (ws wsConn) Read(p []byte) (n int, err error) {
	var r io.Reader

	if _, r, err = ws.Conn.NextReader(); err != nil {
		return
	}
	if n, err = r.Read(p); err != nil {
		return
	}

	return
}

func (ws wsConn) Write(p []byte) (n int, err error) {
	var w io.WriteCloser

	if w, err = ws.Conn.NextWriter(websocket.BinaryMessage); err != nil {
		return
	}
	defer func() {
		err = w.Close()
	}()

	if n, err = w.Write(p); err != nil {
		return
	}

	return
}

func (ws wsConn) Close() error {
	return ws.Conn.Close()
}
