package main

import (
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

	log.Println("SimpleSigner")
	simple.Println("Everyone shares a secret.")
	simple.Println("Useful in very simple scenarios.")
	simple.Println("But long lived processes are going to need key rotation.")

	log.Println("IdentifiedSigner")
	ident.Println("You can like totally trust that Author.1 created this entry.")
	ident.Println("You can be sure that nobody modified it or spoofed the identify.")
	ident.SetOutput(idw2)
	ident.Println("Even better, when you rotate and change the key version you can still read the whole log.")
}
