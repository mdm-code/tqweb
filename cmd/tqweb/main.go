package main

import (
	"log"

	"github.com/mdm-code/tqweb/server"
)

func main() {
	s := server.Server()
	log.Fatal(s.Start(":8000"))
}
