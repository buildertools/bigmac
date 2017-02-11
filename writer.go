package bigmac

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
)

type Signer struct {
	writer io.Writer
	secret []byte
}

func (s Signer) Write(p []byte) (n int, err error) {
	raw := hmac.New(sha256.New, s.secret)
	raw.Write(p)
	mac := base64.StdEncoding.EncodeToString(raw.Sum(nil))
	i, err := s.writer.Write([]byte(mac))
	if err != nil {
		return i, err
	}
	j, err := s.writer.Write(p)
	return i+j, err
}

func NewSigner(w io.Writer, secret []byte) *Signer {
	return &Signer{writer: w, secret: secret}
}
