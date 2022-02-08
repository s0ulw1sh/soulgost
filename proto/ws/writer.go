package ws

import (
	"encoding/binary"
)

type Writer struct {
	Type int
	buf [defaultWriteBufferSize]byte
	pos int
	c   *Conn
	err error
}

func (w *Writer) ncopy(max int) (int, error) {
	n := defaultWriteBufferSize - w.pos
	if n <= 0 {
		if err := w.flush(false, nil); err != nil {
			return 0, err
		}
		n = defaultWriteBufferSize - w.pos
	}
	if n > max {
		n = max
	}
	return n, nil
}

func (w *Writer) Write(p []byte) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	ln := len(p)

	if ln > defaultWriteBufferSize*2 && w.c.srv {
		err := w.flush(false, p)
		if err != nil {
			return 0, err
		}
		return ln, nil
	}

	for len(p) > 0 {
		n, err := w.ncopy(len(p))
		if err != nil {
			return 0, err
		}
		copy(w.buf[w.pos:], p[:n])
		w.pos += n
		p = p[n:]
	}

	return ln, nil
}

func (w *Writer) Close() error {
	if w.err != nil {
		return w.err
	}
	return w.flush(true, nil)
}

func (w *Writer) errmsg(err error) error {
	if w.err != nil {
		return err
	}

	w.err = err

	return err
}

func (w *Writer) flush(final bool, extra []byte) error {

	ln := w.pos - maxFrameHeaderSize + len(extra)

	b0 := byte(w.Type)
	b1 := byte(0)
	if final {
		b0 |= finalBit
	}

	fpos := 0
	if !w.c.srv {
		b1 |= maskBit
	} else {
		fpos = 4
	}

	switch {
	case ln >= 65536:
		w.buf[fpos]   = b0
		w.buf[fpos+1] = b1 | 127
		binary.BigEndian.PutUint64(w.buf[fpos+2:], uint64(ln))
	case ln > 125:
		fpos         += 6
		w.buf[fpos]   = b0
		w.buf[fpos+1] = b1 | 126
		binary.BigEndian.PutUint16(w.buf[fpos+2:], uint16(ln))
	default:
		fpos         += 8
		w.buf[fpos]   = b0
		w.buf[fpos+1] = b1 | byte(ln)
	}

	if !w.c.srv {
		key := keyMask()
		copy(w.buf[maxFrameHeaderSize-4:], key[:])
		maskData(key, 0, w.buf[maxFrameHeaderSize:w.pos])
		if len(extra) > 0 {
			return w.errmsg(ErrFrame)
		}
	}

	_, err := w.c.conn.Write(w.buf[fpos:w.pos])

	if err != nil {
		return w.errmsg(err)
	}

	if len(extra) > 0 {
		_, err := w.c.conn.Write(extra)

		if err != nil {
			return w.errmsg(err)
		}
	}

	if final {
		w.errmsg(ErrClosed)
		return nil
	}

	w.pos  = maxFrameHeaderSize
	w.Type = ContinueFrame

	return nil
}