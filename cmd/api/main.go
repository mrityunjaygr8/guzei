package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("woo")

	server := NewApplication()

	http.ListenAndServe(":3000", server)
}
