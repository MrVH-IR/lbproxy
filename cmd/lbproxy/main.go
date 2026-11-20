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
	"libproxy/internal/middleware"
)

func main() {
	configPath := flag.String("config", "config/config.yml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	var backends []*backend.Backend
	for _, raw := range cfg.Backends {
		b, err := backend.New(raw)
		if err != nil {
			log.Fatalf("backend config: %v", err)
		}
		backends = append(backends, b)
	}

	var bal balancer.Balancer
	if cfg.Algorithm == "round_robin" {
		bal = balancer.NewRoundRobin(backends)
	} else {
		bal = balancer.NewLeastConn(backends)
	}
	logger, err := middleware.NewLogger(cfg.LogDir, cfg.LogRetentionDays)
	if err != nil {
		log.Fatalf("initialize logs: %v", err)
	}
	defer logger.Close()

	proxy := &httputil.ReverseProxy{
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			logger.Errorf("proxy error ip=%q method=%s path=%q error=%v", r.RemoteAddr, r.Method, r.URL.RequestURI(), err)
			http.Error(w, "Bad Gateway", http.StatusBadGateway)
		},
	}
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b := bal.Next()
		b.Start()
		defer b.Done()
		p := *proxy
		p.Rewrite = func(pr *httputil.ProxyRequest) { pr.SetURL(b.URL); pr.SetXForwarded() }
		p.ServeHTTP(w, r)
	})

	srv := &http.Server{
		Addr:              cfg.Listen,
		Handler:           logger.Middleware(baseHandler),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("proxy listening on %s with %d backends (%s)", srv.Addr, len(backends), cfg.Algorithm)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
