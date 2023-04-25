package main

import (
	"log"

	"github.com/rafaelbeecker/mwskit/internal/cmd/toolkit"
)

func main() {
	if err := toolkit.Run(); err != nil {
		log.Fatalln(err)
	}
}
