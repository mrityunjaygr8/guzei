package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	log.Println("woo")

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Hello World"))
	})

	http.ListenAndServe(":3000", r)
}
