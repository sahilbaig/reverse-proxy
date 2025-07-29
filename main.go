package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	reverseproxy "github.com/sahilbaig/reverse-proxy/internal/reverProxy"
)

func main() {
	target := os.Getenv("PROXY_TARGET")
	if target == "" {
		log.Fatal("PROXY_TARGET env variable not set")
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", reverseproxy.Reverseproxy(target))
	// r.Handle("/*")

	log.Println("Listening on :7001")
	log.Fatal(http.ListenAndServe(":7001", r))
}
