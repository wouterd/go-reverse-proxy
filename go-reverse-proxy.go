package main

import (
	"os"
	"log"
	"net/http"
	"net/url"
	"net/http/httputil"
	"strings"
	"math/rand"
)

func main() {
	backend, exists := os.LookupEnv("BACKEND")
	if !exists {
		backend = "http://localhost:1234"
	}

	backends := strings.Fields(backend)
	var proxies []*httputil.ReverseProxy
	for _, bkend := range backends {
		bk, err := url.Parse(bkend)
		if err != nil {
			log.Fatalln(err.Error())
		}
		log.Printf("Creating a reverse proxy for %v", bk)
		proxy := httputil.NewSingleHostReverseProxy(bk)
		proxies = append(proxies, proxy)
	}

	http.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		proxyNr := rand.Intn(len(backends))
		proxies[proxyNr].ServeHTTP(w, r)
	})

	http.HandleFunc("/foo/", func(w http.ResponseWriter, r *http.Request) {
		proxyNr := rand.Intn(len(backends))
		proxies[proxyNr].ServeHTTP(w, r)
	})

	port, exists := os.LookupEnv("SERVER_PORT")
	if !exists {
		port = "80"
	}

	log.Printf("Starting a server on port %v", port)
	log.Fatal(http.ListenAndServe(":" + port, nil))
}
