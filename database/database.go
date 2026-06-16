package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init(schemaPath string) {
	var err error
	DB, err = sql.Open("sqlite3", "./forum.db?_foreign_keys=on")
	if err != nil {
		log.Fatal("Erreur ouverture BDD:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Erreur connexion BDD:", err)
	}

	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		log.Fatal("Erreur lecture schema.sql:", err)
	}

	if _, err = DB.Exec(string(schema)); err != nil {
		log.Fatal("Erreur exécution schema.sql:", err)
	}

	// Migrations silencieuses pour les bases créées avant l'ajout de ces colonnes.
	DB.Exec(`ALTER TABLE users ADD COLUMN profile_photo TEXT NOT NULL DEFAULT ''`)
	DB.Exec(`ALTER TABLE users ADD COLUMN role TEXT NOT NULL DEFAULT 'user'`)

	// Si la base existante n'a aucun admin, promouvoir le plus ancien utilisateur.
	var adminCount int
	if DB.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'admin'").Scan(&adminCount); adminCount == 0 {
		DB.Exec("UPDATE users SET role = 'admin' WHERE id = (SELECT id FROM users ORDER BY created_at ASC LIMIT 1)")
	}

	log.Println("Base de données initialisée")
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
