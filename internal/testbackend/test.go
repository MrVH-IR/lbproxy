package testbackend

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

type server struct {
	name    string
	healthy atomic.Bool
}

func New(name string) http.Handler {
	s := &server{name: name}
	s.healthy.Store(true)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("POST /health/toggle", s.toggle)
	mux.HandleFunc("GET /slow", s.slow)
	mux.HandleFunc("GET /", s.index)
	return mux
}

func (s *server) index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Backend", s.name)
	fmt.Fprintf(w, "backend=%s path=%s", s.name, r.URL.Path)
}

func (s *server) health(w http.ResponseWriter, r *http.Request) {
	if !s.healthy.Load() {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintln(w, "unhealthy")
		return
	}
	fmt.Fprintln(w, "OK")
}

func (s *server) toggle(w http.ResponseWriter, r *http.Request) {
	s.healthy.Store(!s.healthy.Load())
	fmt.Fprintf(w, "healthy=%t\n", s.healthy.Load())
}

func (s *server) slow(w http.ResponseWriter, r *http.Request) {
	d, err := time.ParseDuration(r.URL.Query().Get("d"))
	if err != nil {
		http.Error(w, "invalid duration", http.StatusBadRequest)
		return
	}

	select {
	case <-time.After(d):
		fmt.Fprintf(w, "backend=%s slept=%s\n", s.name, d)
	case <-r.Context().Done():
		return
	}
}
