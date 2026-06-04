# Forum — Projet Go

Forum web complet en Go + SQLite + Docker.

## Lancement avec Docker

```bash
docker compose up --build
```

Accès : http://localhost:8080

## Lancement sans Docker (développement)

```bash
go mod tidy
go run .
```

## Structure

```
forum/
├── main.go               # Routing HTTP principal
├── schema.sql            # Schéma BDD SQLite
├── Dockerfile
├── docker-compose.yml
├── database/             # Connexion et init BDD
├── handlers/             # Logique HTTP (auth, posts, comments, likes)
├── models/               # Structs Go
├── middleware/           # Auth session
├── utils/                # UUID
├── templates/            # HTML
├── static/               # CSS
└── uploads/              # Images uploadées
```

## Technologies

- Go 1.21 (pas de framework)
- SQLite3 via `github.com/mattn/go-sqlite3`
- Hashage mot de passe : `golang.org/x/crypto/bcrypt`
- UUID : `github.com/google/uuid`
