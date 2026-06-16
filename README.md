# SportTalk — Forum Web

Forum de discussion sportif développé en Go avec SQLite, sans aucun framework.

## Fonctionnalités

- **Posts** : création avec image, association à une ou plusieurs catégories, modification, suppression
- **Commentaires** : ajout, modification, suppression sur chaque post
- **Votes** : like / dislike sur les posts et les commentaires (visible sans compte)
- **Filtrage** : par catégorie, mes posts, posts aimés
- **Authentification** : inscription, connexion, déconnexion — 1 session max par utilisateur, expiration 24 h
- **Profil** : changement de mot de passe, photo de profil
- **Rôles** : guest · user · moderator · admin
  - Les modérateurs peuvent supprimer n'importe quel post ou commentaire
  - L'admin gère les rôles depuis le panneau `/admin`
  - Le premier compte inscrit devient automatiquement admin
- **Calendrier** : page `/calendrier` avec tous les matchs de la FIFA World Cup 2026

## Stack technique

| Couche          | Technologie                                  |
|-----------------|----------------------------------------------|
| Backend         | Go — `net/http`, `html/template`             |
| Base de données | SQLite via `go-sqlite3` (CGO requis)         |
| Sécurité        | bcrypt pour les mots de passe                |
| IDs             | UUID v4 via `google/uuid`                    |
| Frontend        | HTML / CSS pur — zéro framework              |
| Sessions        | Cookie HttpOnly, 1 session max / utilisateur |
| Conteneur       | Docker + Docker Compose                      |

## Prérequis

- Go 1.21+
- GCC (`CGO_ENABLED=1` requis par go-sqlite3)  
  → Windows : installer [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) ou MinGW-w64
- Docker (optionnel)

## Lancer avec Docker (recommandé)

```bash
docker compose up --build
```

L'application est accessible sur [http://localhost:8080](http://localhost:8080).  
La base de données et les uploads sont persistés dans des volumes locaux.

## Lancer en local

```bash
go mod tidy
CGO_ENABLED=1 go run .
```

La base de données `forum.db` est créée automatiquement au premier lancement.

## Lancer les tests

```bash
CGO_ENABLED=1 go test ./...
```

## Structure du projet

```
.
├── main.go                  # Routeur HTTP + middleware de logging
├── schema.sql               # Schéma SQLite + données initiales (catégories)
├── Dockerfile               # Build multi-stage (builder Go + runtime alpine)
├── Docker compose.yml
├── database/
│   ├── database.go          # Initialisation + migrations silencieuses
│   ├── users.go
│   ├── posts.go
│   ├── comments.go
│   ├── likes.go
│   ├── sessions.go
│   ├── testsetup_test.go    # Helpers partagés pour les tests
│   ├── session_test.go
│   └── likes_test.go
├── handlers/
│   ├── auth.go              # Inscription / connexion / déconnexion
│   ├── posts.go
│   ├── comments.go
│   ├── profile.go
│   ├── admin.go
│   ├── calendar.go
│   └── helpers.go           # serveInternalError (page 500 HTML)
├── middleware/
│   └── auth.go              # WithUser · RequireAuth · RequireModerator · RequireAdmin
├── models/
│   ├── user.go
│   ├── post.go
│   ├── comment.go
│   └── like.go
├── utils/
│   ├── uuid.go
│   └── uuid_test.go
├── templates/
│   ├── index.html           # Accueil — 3 colonnes (sidebar / feed / calendrier)
│   ├── post.html            # Détail d'un post + commentaires
│   ├── create_post.html     # Formulaire de création
│   ├── login.html
│   ├── register.html
│   ├── profile.html
│   ├── admin.html           # Panneau d'administration des rôles
│   ├── calendrier.html      # FIFA World Cup 2026 — 48 matchs
│   ├── 404.html
│   └── 500.html
└── static/
    └── style.css            # Thème dark SportTalk (CSS pur)
```

## Règles techniques respectées

| Exigence | Statut |
|----------|--------|
| Langage Go uniquement, pas de framework backend | ✅ |
| Base de données SQLite uniquement | ✅ |
| Frontend HTML/CSS pur, pas de React/Bootstrap | ✅ |
| Bibliothèques autorisées : sqlite3, bcrypt, uuid | ✅ |
| `CGO_ENABLED=1` | ✅ |
| Erreurs HTTP 404 et 500 avec page HTML dédiée | ✅ |
| Sessions : cookie HttpOnly, expiration 24 h, 1 max par user | ✅ |
| Conteneurisation Docker | ✅ |
| Tests unitaires | ✅ |
