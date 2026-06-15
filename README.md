# ForumFoot — Projet Go B1

Forum web complet en Go pur + SQLite, inspiré de l'univers football / Coupe du Monde 2026.

## Prérequis

| Outil | Version minimale | Notes |
|-------|-----------------|-------|
| Go | 1.21 | `CGO_ENABLED=1` obligatoire |
| GCC | Toute version récente | Requis pour go-sqlite3 (CGO) |
| Docker + Docker Compose | 24+ | Pour le lancement conteneurisé |

> **Windows** : installer [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) ou MinGW-w64 pour avoir GCC.

## Lancement avec Docker (recommandé)

```bash
docker compose up --build
```

Accès : [http://localhost:8080](http://localhost:8080)

Aucune configuration manuelle requise. La base de données SQLite et les uploads sont persistés dans des volumes Docker.

## Lancement sans Docker (développement)

```bash
# 1. Télécharger les dépendances
go mod tidy

# 2. Compiler et lancer (CGO obligatoire)
CGO_ENABLED=1 go run .
```

Accès : [http://localhost:8080](http://localhost:8080)

La base de données `forum.db` est créée automatiquement au premier lancement.

## Variables d'environnement

Aucune variable d'environnement obligatoire. Toute la configuration est dans le code.

| Variable | Valeur par défaut | Description |
|----------|------------------|-------------|
| `DB_PATH` | `./forum.db` | Chemin de la base SQLite |

## Structure du projet

```
forum/
├── main.go                  # Routing HTTP, middleware logging
├── schema.sql               # Schéma BDD SQLite + catégories par défaut
├── Dockerfile               # Build multi-stage (builder + runtime alpine)
├── docker-compose.yml       # Orchestration Docker
├── go.mod / go.sum          # Dépendances Go
│
├── database/                # Couche d'accès aux données
│   ├── database.go          # Initialisation SQLite + migration
│   ├── users.go             # CRUD utilisateurs
│   ├── sessions.go          # Gestion sessions
│   ├── posts.go             # CRUD posts + filtres
│   ├── comments.go          # CRUD commentaires
│   ├── likes.go             # Système like/dislike
│   └── categories.go        # Catégories
│
├── handlers/                # Handlers HTTP (logique métier)
│   ├── auth.go              # Inscription / connexion / déconnexion
│   ├── posts.go             # Posts (liste, détail, création, modification, suppression)
│   ├── comments.go          # Commentaires
│   ├── likes.go             # Votes
│   └── profile.go           # Profil utilisateur (photo, mot de passe)
│
├── middleware/              # Middlewares
│   └── auth.go              # RequireAuth + WithUser (session → contexte)
│
├── models/                  # Structs de données
│   ├── user.go
│   ├── session.go
│   ├── post.go
│   ├── comment.go
│   └── category.go
│
├── utils/                   # Utilitaires
│   └── uuid.go              # Génération UUID
│
├── templates/               # Templates HTML (Go html/template)
│   ├── index.html           # Page d'accueil (3 colonnes)
│   ├── post.html            # Détail d'un post
│   ├── create_post.html     # Formulaire de création
│   ├── login.html           # Connexion
│   ├── register.html        # Inscription
│   ├── profile.html         # Profil utilisateur
│   ├── 404.html             # Page introuvable
│   └── 500.html             # Erreur serveur
│
├── static/                  # Fichiers statiques
│   └── style.css            # CSS pur (pas de framework)
│
└── uploads/                 # Images uploadées (créé automatiquement)
```

## Fonctionnalités

- Inscription / Connexion / Déconnexion (bcrypt, session cookie HttpOnly)
- Session unique par utilisateur, expiration 24h
- Création / modification / suppression de posts (auteur uniquement)
- Commentaires avec modification et suppression
- Like / dislike sur posts et commentaires (toggle, visible en guest)
- Catégories multi-sélection et filtres (par catégorie, mes posts, aimés)
- Upload d'image dans les posts (JPEG/PNG/GIF/WebP, 20 Mo max)
- Page profil : changement de mot de passe + photo de profil
- Pages d'erreur 404 et 500 personnalisées
- Layout responsive 3 colonnes (matchs WC2026, feed, actualités foot)

## Technologies

- **Langage :** Go 1.21 (pas de framework HTTP)
- **Base de données :** SQLite3 via `github.com/mattn/go-sqlite3` (CGO)
- **Hashage :** `golang.org/x/crypto/bcrypt`
- **UUID :** `github.com/google/uuid`
- **Templates :** `html/template` (stdlib Go)
- **Frontend :** HTML/CSS pur (pas de React, Bootstrap ou autre framework)
