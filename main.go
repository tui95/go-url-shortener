package main

import (
	"github.com/tui95/go-url-shortener/internal/server"
)

func main() {
	port := ":8080"
	server.RunServer(port)
}
