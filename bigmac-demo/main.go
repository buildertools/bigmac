package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"github.com/buildertools/bigmac"
	"log"
	"os"
)

func main() {
	sw := bigmac.NewSimpleSigner(os.Stdout, []byte("This is a demo secret"))
	simple := log.New(sw, ``, log.Flags())

	idw1 := bigmac.NewIdentifiedSigner(os.Stdout, "Author.1", []byte("This is a demo secret"))
	ident := log.New(idw1, ``, log.Flags())
	idw2 := bigmac.NewIdentifiedSigner(os.Stdout, "Author.2", []byte("This is a demo secret"))

	ecpk, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	ecdsa1 := bigmac.NewIdentifiedECDSASigner(os.Stdout, "generated-ECDSA", ecpk)

	pk, _ := rsa.GenerateKey(rand.Reader, 2048)
	pkcs1 := bigmac.NewIdentifiedPKCS1v15Signer(os.Stdout, "generated-rsa", pk)

	log.Println("SimpleSigner")
	simple.Println("Everyone shares a secret.")
	simple.Println("Useful in very simple scenarios.")
	simple.Println("But long lived processes are going to need key rotation.")

	log.Println("IdentifiedSigner")
	ident.Println("You can like totally trust that Author.1 created this entry.")
	ident.Println("You can be sure that nobody modified it or spoofed the identify.")
	ident.SetOutput(idw2)
	ident.Println("Even better, when you rotate and change the key version you can still read the whole log.")

	log.Println("IdentifiedECDSASigner")
	ident.SetOutput(ecdsa1)
	ident.Println("This uses an ECDSA signature. Pretty fancy.")

	log.Println("IdentifiedPKCS1v15Signer")
	ident.SetOutput(pkcs1)
	ident.Println("This example uses a 2048 bit RSA key and creates a PKCS1v15 signature.")
	ident.Println("This generates a really long signature and takes a long time.")

}
