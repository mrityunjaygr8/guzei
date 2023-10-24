package main

import (
	"github.com/mrityunjaygr8/guzei/internal/postgres_store"
	"github.com/mrityunjaygr8/guzei/internal/server"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("woo")

	pgStore, err := postgres_store.NewPostgresStore(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("error creating postgres store: ", err)
	}

	srv := server.NewServer(pgStore)

	log.Fatal(http.ListenAndServe(":3000", srv))

}
