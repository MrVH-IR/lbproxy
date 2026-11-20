package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	port := flag.Int("port", 9101, "port to listen on")
	name := flag.String("name", "b1", "backend name")
	flag.Parse()

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", *port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "backend=%s path=%s\n", *name, r.URL.Path)
		}),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("backend %s listening on %s", *name, srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
