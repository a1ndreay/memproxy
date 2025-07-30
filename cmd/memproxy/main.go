package main

import (
	"flag"
	"log"
	"os"

	fc "github.com/a1ndreay/memproxy/internal/fake"
	mc "github.com/a1ndreay/memproxy/internal/memcached"
	"github.com/a1ndreay/memproxy/internal/server"
	"github.com/a1ndreay/memproxy/pkg/cache"
)

func main() {
	// CLI flags
	var listenAddr string
	var cacheAddr string
	var originAddr string
	var useBackend string

	flag.StringVar(&listenAddr, "listen", ":8080", "address to listen on (e.g. :8080)")
	flag.StringVar(&cacheAddr, "cache-address", "localhost:11211", "cache address")
	flag.StringVar(&originAddr, "origin", "http://localhost:8081", "origin server address")
	flag.StringVar(&useBackend, "backend", "memcached", "cache to use (default memcached). values: memcached, inmemory")
	flag.Parse()

	// initialize backend
	var backend cache.Backend
	switch useBackend {
	case "inmemory":
		backend = fc.New()
	default:
		backend = mc.New(cacheAddr)
	}

	// start HTTP server
	srv := server.New(backend, originAddr)

	log.Printf("Starting memproxy on %s, in mode %s %s, origin was %s", listenAddr, useBackend, cacheAddr, originAddr)
	if err := srv.ListenAndServe(listenAddr); err != nil {
		log.Fatalf("Server error: %v", err)
		os.Exit(1)
	}
}
