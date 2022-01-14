package thirdparty

import (
	"time"
	"bytes"
	"strings"
	"strconv"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
)

func g2FATokenPfx(otp string) string {
	if len(otp) == 6 {
		return otp
	}
	for i := (6 - len(otp)); i > 0; i-- {
		otp = "0" + otp
	}
	return otp
}

func G2FAToken(secret string, interval int64) (string, error) {
	key, err := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
	if err != nil {
		return "", err
	}

	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, uint64(interval))

	hash := hmac.New(sha1.New, key)
	hash.Write(bs)
	h := hash.Sum(nil)

	o := (h[19] & 15)

	var header uint32

	r := bytes.NewReader(h[o : o+4])
	err = binary.Read(r, binary.BigEndian, &header)

	if err != nil {
		return "", err
	}

	h12 := (int(header) & 0x7fffffff) % 1000000

	otp := strconv.Itoa(int(h12))

	return g2FATokenPfx(otp), nil
}

func G2FAToken30Sec(secret string) (string, error) {
	return G2FAToken(secret, time.Now().Unix() / 30)
}