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

	log.Println("Base de données initialisée")
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
