package main

import (
	"DarcyProxy/internal/testbackend"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	port := flag.Int("port", 9101, "port to listen on")
	name := flag.String("name", "b1", "name of the backend")
	flag.Parse()

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", *port),
		Handler:           testbackend.New(*name),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("backend %s listening on %s", *name, srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
