package main

import (
	"github.com/buildertools/bigmac"
	"log"
	"os"
)

func main() {
	w := bigmac.NewSigner(os.Stdout, []byte("This is a demo secret"))
	log.SetOutput(w)

	log.Println("Hello, World!")
	log.Println("Different message")
}
