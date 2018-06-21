package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	defaultPort = 8080
)

func main() {
	var port int
	if tPort, ok := os.LookupEnv("PORT"); ok {
		var err error
		port, err = strconv.Atoi(tPort)
		if err != nil {
			panic(err)
		}
	} else {
		port = defaultPort
	}

	listenAddr := fmt.Sprintf(":%d", port)

	router := newRouter()

	log.Fatal(http.ListenAndServe(listenAddr, router))
}
