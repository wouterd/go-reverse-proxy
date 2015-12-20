package main

import (
	"os"
	"log"
	"net/http"
	"net/url"
	"net/http/httputil"
)

func main() {
	backend, exists := os.LookupEnv("BACKEND")
	if !exists {
		backend = "http://localhost:1234"
	}

	bk, err := url.Parse(backend)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Printf("Creating a reverse proxy for %v", bk)
	proxy := httputil.NewSingleHostReverseProxy(bk)

	http.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	port, exists := os.LookupEnv("SERVER_PORT")
	if !exists {
		port = "80"
	}

	log.Printf("Starting a server on port %v", port)
	log.Fatal(http.ListenAndServe(":" + port, nil))
}
