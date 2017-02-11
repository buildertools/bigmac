package bigmac

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

func sign(p []byte, s []byte) string {
	r := hmac.New(sha256.New, s)
	r.Write(p)
	return base64.StdEncoding.EncodeToString(r.Sum(nil))
}

type SimpleSigner struct {
	writer io.Writer
	secret []byte
}

func (s SimpleSigner) Write(p []byte) (n int, err error) {
	mac := sign(p, s.secret)
	i, err := s.writer.Write([]byte(mac))
	if err != nil {
		return i, err
	}
	j, err := s.writer.Write(p)
	return i+j, err
}

func NewSimpleSigner(w io.Writer, secret []byte) io.Writer {
	return &SimpleSigner{writer: w, secret: secret}
}

type IdentifiedSigner struct {
	SimpleSigner
	name string
}

const rid = "%v %v %v"
func (s IdentifiedSigner) Write(p []byte) (n int, err error) {
	mac := sign(p, s.secret)
	
	i, err := s.writer.Write([]byte(fmt.Sprintf(rid, mac, s.name, string(p))))
	return i, err
}

func NewIdentifiedSigner(w io.Writer, name string, secret []byte) io.Writer {
	return &IdentifiedSigner{SimpleSigner: SimpleSigner{writer: w, secret: secret}, name: name}
}

