package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	store := newPostgresStore()
	srv := newServer(store)
	handler := srv.router()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("backend startet auf :%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("server-fehler: %v", err)
	}
}
