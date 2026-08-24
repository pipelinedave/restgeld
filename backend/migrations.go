package main

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"sort"
	"strings"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

// RunMigrations führt alle noch nicht angewandten Datenbank-Migrationen aus.
func RunMigrations(db *sql.DB) error {
	// 1. Tabelle schema_migrations erstellen, falls nicht vorhanden
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version VARCHAR(255) PRIMARY KEY,
		applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`
	if _, err := db.Exec(createTableSQL); err != nil {
		return fmt.Errorf("schema_migrations tabelle erstellen: %w", err)
	}

	// 2. Bereits angewandte Migrationen abrufen
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return fmt.Errorf("angewandte migrationen lesen: %w", err)
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return fmt.Errorf("migrations-version scannen: %w", err)
		}
		applied[v] = true
	}

	// 3. Migrationsdateien aus dem eingebetteten Dateisystem laden und sortieren
	entries, err := fs.ReadDir(migrationFS, "migrations")
	if err != nil {
		return fmt.Errorf("migrations-verzeichnis lesen: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)

	// 4. Ausstehende Migrationen der Reihe nach anwenden
	for _, file := range files {
		if applied[file] {
			continue
		}

		content, err := migrationFS.ReadFile("migrations/" + file)
		if err != nil {
			return fmt.Errorf("migrationsdatei %s lesen: %w", file, err)
		}

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("transaktion fuer %s starten: %w", file, err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("migration %s ausfuehren: %w", file, err)
		}

		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", file); err != nil {
			tx.Rollback()
			return fmt.Errorf("migrations-status fuer %s speichern: %w", file, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("transaktion fuer %s committen: %w", file, err)
		}

		log.Printf("migration erfolgreich angewendet: %s", file)
	}

	return nil
}
