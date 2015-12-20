package main

import (
	"os"
	"log"
	"net/http"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"net/http/httputil"
	"math/rand"
	"net/url"
)

type Endpoint struct {
	Path       string    `yaml:"path"`
	Components []string    `yaml:"components"`
}

var config map[string]Endpoint

func main() {

	filename, exists := os.LookupEnv("CONFIG")
	if !exists {
		filename = "/etc/endpoints.yaml"
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err.Error())
	}

	for key, config := range config {
		log.Printf("Creating proxies for %v", key)
		var proxies []*httputil.ReverseProxy
		for _, comp := range config.Components {
			url, err := url.Parse(comp)
			if err != nil {
				log.Fatal(err.Error())
			}
			log.Printf("Creating a reverse proxy for %v", url.String())
			proxy := httputil.NewSingleHostReverseProxy(url)
			proxies = append(proxies, proxy)
		}
		log.Printf("Creating a handler for %v", config.Path)
		http.HandleFunc(config.Path, func(w http.ResponseWriter, r *http.Request) {
			proxyNr := rand.Intn(len(proxies))
			proxies[proxyNr].ServeHTTP(w, r)
		})
	}

	port, exists := os.LookupEnv("SERVER_PORT")
	if !exists {
		port = "80"
	}

	log.Printf("Starting a server on port %v", port)
	log.Fatal(http.ListenAndServe(":" + port, nil))
}
