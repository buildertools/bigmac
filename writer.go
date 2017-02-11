package bigmac

import (
	"crypto"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/ecdsa"
	"encoding/base64"
	"fmt"
	"io"
)

func signHMAC(p []byte, s []byte) string {
	r := hmac.New(sha256.New, s)
	r.Write(p)
	return base64.StdEncoding.EncodeToString(r.Sum(nil))
}

func signPKCS1v15(p []byte, s *rsa.PrivateKey) (string, error) {
	hashed := sha256.Sum256(p)
	sig, err := rsa.SignPKCS1v15(rand.Reader, s, crypto.SHA256, hashed[:])
	return base64.StdEncoding.EncodeToString(sig), err
}

func signECDSA(p []byte, k *ecdsa.PrivateKey) (string, string, error) {
	hashed := sha256.Sum256(p)
	r, s, err := ecdsa.Sign(rand.Reader, k, hashed[:])
	return base64.StdEncoding.EncodeToString(r.Bytes()), base64.StdEncoding.EncodeToString(s.Bytes()), err
}

type SimpleSigner struct {
	writer io.Writer
	secret []byte
}

func (s SimpleSigner) Write(p []byte) (n int, err error) {
	mac := signHMAC(p, s.secret)
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
// Prepends the input slice with the HMAC and secret name as provided to the IdentifiedSigner
func (s IdentifiedSigner) Write(p []byte) (n int, err error) {
	mac := signHMAC(p, s.secret)
	
	i, err := s.writer.Write([]byte(fmt.Sprintf(rid, mac, s.name, string(p))))
	return i, err
}

func NewIdentifiedSigner(w io.Writer, name string, secret []byte) io.Writer {
	return &IdentifiedSigner{SimpleSigner: SimpleSigner{writer: w, secret: secret}, name: name}
}

type IdentifiedPKCS1v15Signer struct {
	writer io.Writer
	key    *rsa.PrivateKey
	name   string
}

func NewIdentifiedPKCS1v15Signer(w io.Writer, name string, secret *rsa.PrivateKey) io.Writer {
	return &IdentifiedPKCS1v15Signer{writer: w, name: name, key: secret}
}

const formatIdentifiedPKCS1v15 = "%v %v %v"
func (s IdentifiedPKCS1v15Signer) Write(p []byte) (int, error) {
	sig, err := signPKCS1v15(p, s.key)
	i, err := s.writer.Write([]byte(fmt.Sprintf(rid, sig, s.name, string(p))))
	return i, err
}

type IdentifiedECDSASigner struct {
	writer io.Writer
	key    *ecdsa.PrivateKey
	name   string
}

func NewIdentifiedECDSASigner(w io.Writer, name string, secret *ecdsa.PrivateKey) io.Writer {
	return &IdentifiedECDSASigner{writer: w, name: name, key: secret}
}

const ridecdsa = `%v %v %v %v`
func (s IdentifiedECDSASigner) Write(p []byte) (int, error) {
	rs, ss, err := signECDSA(p, s.key)
	i, err := s.writer.Write([]byte(fmt.Sprintf(ridecdsa, rs, ss, s.name, string(p))))
	return i, err
}
