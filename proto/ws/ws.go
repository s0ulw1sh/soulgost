package ws

import (
	"io"
	"io/ioutil"
	"errors"
	"encoding/binary"
)

const (
	TextMessage   = 1
	BinaryMessage = 2
	CloseMessage  = 8
	PingMessage   = 9
	PongMessage   = 10
	ContinueFrame = 0
)

const (
	CloseNormalClosure           = 1000
	CloseGoingAway               = 1001
	CloseProtocolError           = 1002
	CloseUnsupportedData         = 1003
	CloseNoStatusReceived        = 1005
	CloseAbnormalClosure         = 1006
	CloseInvalidFramePayloadData = 1007
	ClosePolicyViolation         = 1008
	CloseMessageTooBig           = 1009
	CloseMandatoryExtension      = 1010
	CloseInternalServerErr       = 1011
	CloseServiceRestart          = 1012
	CloseTryAgainLater           = 1013
	CloseTLSHandshake            = 1015
)

const (
	finalBit = 1 << 7
	rsv1Bit  = 1 << 6
	rsv2Bit  = 1 << 5
	rsv3Bit  = 1 << 4
	maskBit  = 1 << 7

	maxFrameHeaderSize         = 2 + 8 + 4
	maxControlFramePayloadSize = 125

	//writeWait = time.Second

	defaultReadBufferSize  = 4096
	defaultWriteBufferSize = 4096
)

var ErrFrame  = errors.New("websocket: invalid frame")
var ErrRead   = errors.New("websocket: invalid read")
var ErrClosed = errors.New("websocket: write closed")

func (c *Conn) NextReader(msg *Reader) error {
	var err error = nil

	msg.c  = c

	for err == nil {
		e := c.nextFrame(msg)
		
		if err != nil {
			err = e
			break
		}

		if msg.Type == TextMessage || msg.Type == BinaryMessage {
			return nil
		}
	}

	c.errcnt += 1

	return err
}

func (c *Conn) NextWriter(ft int, msg *Writer) {
	msg.c    = c
	msg.Type = ft
	msg.pos  = maxFrameHeaderSize
	msg.err  = nil
}

func (c *Conn) DirectWrite(mtype int, data []byte) error {
	var w Writer

	c.NextWriter(mtype, &w)

	if _, err := w.Write(data); err != nil {
		return err
	}

	return w.Close()
}

func (c *Conn) WriteControl(mtype int, data []byte) error {
	if len(data) > maxControlFramePayloadSize {
		return ErrFrame
	}

	b0 := byte(mtype) | finalBit
	b1 := byte(len(data))
	if !c.srv {
		b1 |= maskBit
	}

	buf := make([]byte, 0, maxFrameHeaderSize+maxControlFramePayloadSize)
	buf = append(buf, b0, b1)

	if c.srv {
		buf = append(buf, data...)
	} else {
		key := keyMask()
		buf = append(buf, key[:]...)
		buf = append(buf, data...)
		maskData(key, 0, buf[6:])
	}

	_, err := c.conn.Write(buf)

	return err
}

func (c *Conn) WriteClose(code int, payload string) error {
	var buf []byte = []byte{}

	if code != CloseNoStatusReceived {
		buf = make([]byte, 2+len(payload))
		binary.BigEndian.PutUint16(buf, uint16(code))
		copy(buf[2:], payload)
	}

	return c.WriteControl(CloseMessage, buf)
}

func (c *Conn) nextFrame(msg *Reader) error {

	if msg.rm > 0 {
		if _, err := io.CopyN(ioutil.Discard, c.brw, int64(msg.rm)); err != nil {
			return err
		}
	}

	p, err := c.read(2)

	if err != nil {
		return err
	}

	msg.Type    = int(p[0] & 0xf)
	msg.final   = p[0]&finalBit != 0
	msg.usemask = p[1]&maskBit != 0
	msg.rm      = int(p[1] & 0x7f)

	// rsv1      := p[0]&rsv1Bit != 0
	// rsv2      := p[0]&rsv2Bit != 0
	// rsv3      := p[0]&rsv3Bit != 0

	switch msg.Type {
	case CloseMessage, PingMessage, PongMessage:
		if msg.rm > maxControlFramePayloadSize || !msg.final {
			return ErrFrame
		}
	case TextMessage, BinaryMessage:
		if !msg.final {
			return ErrFrame
		}
	case ContinueFrame:
		if msg.final {
			return ErrFrame
		}
	default:
		return ErrFrame
	}

	switch msg.rm {
	case 126:
		p, err := c.read(2)
		if err != nil {
			return err
		}

		msg.rm = int(binary.BigEndian.Uint16(p))
	case 127:
		p, err := c.read(8)
		if err != nil {
			return err
		}

		msg.rm = int(binary.BigEndian.Uint64(p))
	}

	if msg.usemask {
		p, err := c.read(4)
		if err != nil {
			return err
		}
		copy(msg.mask[:], p)
	}

	if msg.Type == ContinueFrame || msg.Type == TextMessage || msg.Type == BinaryMessage {
		msg.ln += msg.rm

		if msg.ln < 0 {
			return ErrFrame
		}

		return nil
	}

	var payload []byte

	if msg.rm > 0 {
		payload, err = c.read(int(msg.rm))
		msg.rm = 0
		
		if err != nil {
			return err
		}

		if c.srv {
			maskData(msg.mask, 0, payload)
		}
	}

	switch msg.Type {
	case PingMessage:
		if err := c.WriteControl(PongMessage, payload); err != nil {
			return err
		}
	}

	return nil
}