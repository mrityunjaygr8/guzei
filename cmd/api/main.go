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

func helloWorldServer(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Hello World"))
}

func pingPongHandler(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("PONG"))
}
