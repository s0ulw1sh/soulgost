package ws

import (
	"io"
	"net"
	"bufio"
)

type Conn struct {
	conn   net.Conn
	brw    *bufio.ReadWriter
	srv    bool
	errcnt int
}

func (c *Conn) Close() error {
	return c.conn.Close()
}

func (c *Conn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Conn) read(n int) ([]byte, error) {
	p, err := c.brw.Peek(n)
	if err == io.EOF {
		err = ErrRead
	}
	c.brw.Discard(len(p))
	return p, err
}