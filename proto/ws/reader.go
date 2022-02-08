package ws

import (
	"io"
)

type Reader struct {
	Type    int
	mask    [4]byte
	usemask bool
	ln      int
	rm      int
	final   bool
	c       *Conn
}

func (m *Reader) Read(b []byte) (int, error) {
	var err error

	for err == nil {
		if m.rm > 0 {
			if len(b) > m.rm {
				b = b[:m.rm]
			}

			n, e := m.c.brw.Read(b)

			if e != nil {
				err = e
			}

			if m.c.srv {
				for i := 0; i < n; i++ {
					b[i] ^= m.mask[i % 4]
				}
			}

			m.rm -= n

			if m.rm > 0 && e == io.EOF {
				e = ErrRead
			}

			return n, e
		}

		if m.final {
			return 0, io.EOF
		}

		err = m.c.NextReader(m)
	}

	if err == io.EOF {
		err = ErrRead
	}

	return 0, err
}