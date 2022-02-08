package ws

import (
	"time"
	"net/http"
	"crypto/sha1"
	"encoding/base64"
)

const ws_magic = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

func HandshakeAccessKey(key string) string {
	h := sha1.New()
	h.Write([]byte(key))
	h.Write([]byte(ws_magic))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func WsUpgrader(w http.ResponseWriter, r *http.Request) (Conn, error) {

	if r.Method != http.MethodGet {
		http.Error(w, "Request method is not GET", http.StatusBadRequest)
		return Conn{}, http.ErrBodyNotAllowed
	}

	key := r.Header.Get("Sec-WebSocket-Key")

	if key == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return Conn{}, http.ErrBodyNotAllowed
	}

	h, ok := w.(http.Hijacker)

	if !ok {
		http.Error(w, "Hijack error", http.StatusInternalServerError)
		return Conn{}, http.ErrHijacked
	}

	netConn, brw, err := h.Hijack()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return Conn{}, http.ErrHijacked
	}

	if brw.Reader.Buffered() > 0 {
		netConn.Close()
		return Conn{}, http.ErrWriteAfterFlush
	}

	netConn.SetDeadline(time.Time{})

	out := "HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: "
	out += HandshakeAccessKey(key) + "\r\n\r\n"

	if _, err = netConn.Write([]byte(out)); err != nil {
		netConn.Close()
		return Conn{}, err
	}

	return Conn{netConn, brw, true, 0}, nil
}