package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

var (
	handlerInstance http.Handler
	once            sync.Once
)

func initApp() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		host := envOrDefault("DB_HOST", "localhost")
		port := envOrDefault("DB_PORT", "5432")
		user := envOrDefault("DB_USER", "restgeld")
		password := envOrDefault("DB_PASSWORD", "restgeld")
		dbname := envOrDefault("DB_NAME", "restgeld")

		connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("fehler beim db-verbindungsaufbau: %v", err)
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := RunMigrations(db); err != nil {
		log.Printf("warnung/fehler bei db-migration: %v", err)
	}

	store := &postgresStore{db: db}
	srv := newServer(store)
	handlerInstance = srv.router()
}

func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(initApp)
	if handlerInstance != nil {
		handlerInstance.ServeHTTP(w, r)
	} else {
		http.Error(w, "server nicht initialisiert", http.StatusInternalServerError)
	}
}
