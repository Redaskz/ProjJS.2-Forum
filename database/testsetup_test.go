package database

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// setupTestDB initialise une base SQLite en mémoire avec le schéma complet.
// La connexion est fermée automatiquement à la fin du test.
func setupTestDB(t *testing.T) {
	t.Helper()
	var err error
	DB, err = sql.Open("sqlite3", ":memory:?_foreign_keys=on")
	if err != nil {
		t.Fatalf("impossible d'ouvrir la base de test : %v", err)
	}

	schema, err := os.ReadFile("../schema.sql")
	if err != nil {
		t.Fatalf("impossible de lire schema.sql : %v", err)
	}

	if _, err = DB.Exec(string(schema)); err != nil {
		t.Fatalf("impossible d'initialiser le schéma : %v", err)
	}

	t.Cleanup(func() { DB.Close() })
}

// createTestUser insère un utilisateur minimal pour les tests qui en ont besoin.
func createTestUser(t *testing.T, id, email, username string) {
	t.Helper()
	_, err := DB.Exec(
		`INSERT INTO users (id, email, username, password_hash, created_at) VALUES (?, ?, ?, 'hash', CURRENT_TIMESTAMP)`,
		id, email, username,
	)
	if err != nil {
		t.Fatalf("impossible de créer l'utilisateur de test : %v", err)
	}
}
