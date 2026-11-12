package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"libproxy/internal/backend"
	"libproxy/internal/balancer"
	"libproxy/internal/config"
)

func main() {
	configPath := flag.String("config", "config/config.yml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	// TODO: از cfg.Backends یک []*backend.Backend بساز.
	// راهنما: حلقه با range روی cfg.Backends، صدا زدن backend.New(raw)،
	// اگه خطا بود log.Fatalf، وگرنه append به slice.
	var backends []*backend.Backend

	bal := balancer.NewRoundRobin(backends)

	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			b := bal.Next()
			r.SetURL(b.URL)
			r.SetXForwarded()
		},
	}

	srv := &http.Server{
		Addr:              cfg.Listen,
		Handler:           proxy,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("proxy listening on %s with %d backends (%s)", srv.Addr, len(backends), cfg.Algorithm)
	log.Fatal(srv.ListenAndServe())
}
