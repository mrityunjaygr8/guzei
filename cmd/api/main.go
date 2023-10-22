package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("woo")

	pgStore, err := NewPostgresStore(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("error creating postgres store: ", err)
	}

	server := NewApplication(pgStore)

	http.ListenAndServe(":3000", server)
}
