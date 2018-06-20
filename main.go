package main

import (
	"log"
	"net/http"
	"os"
)

const (
	defaultListenAddr = ":8080"
)

func main() {
	listenAddr, ok := os.LookupEnv("LISTEN_ADDR")
	if !ok {
		listenAddr = defaultListenAddr
	}

	router := NewRouter()

	log.Fatal(http.ListenAndServe(listenAddr, router))
}
