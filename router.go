package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const (
	defaultStaticDir = "./app/"
)

func newRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	var (
		staticDir string
		ok        bool
	)

	staticDir, ok = os.LookupEnv("STATIC_DIR")
	if !ok {
		staticDir = defaultStaticDir
	}

	router.PathPrefix("/").Handler(logger(http.FileServer(http.Dir(staticDir)), "Default"))

	return router
}
