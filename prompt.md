# Journal des prompts IA — Projet Forum Go B1

Ce fichier documente les prompts utilisés avec Claude (Anthropic) pour nous aider à concevoir et développer ce projet. L'IA a été utilisée comme outil d'aide à la conception et à la génération de code, que chaque membre a ensuite lu, compris et adapté.

---

## Prompt 1 — Conception de l'architecture et gestion de projet

**Contexte :** Début du projet, nous devions organiser le travail entre 4 membres et définir la structure technique avant de commencer à coder.

**Prompt utilisé :**
> "Nous sommes 4 étudiants en B1 et nous devons réaliser un forum web complet en Go pur avec SQLite, sans framework backend ni frontend. Le forum doit permettre la communication entre utilisateurs via des posts et commentaires, la gestion de catégories, un système de likes/dislikes et le filtrage des posts. Génère-nous une gestion de projet complète avec la répartition des tâches par membre, un planning sur 4 semaines en sprints, une matrice RACI et une analyse des risques. Chaque membre doit avoir un rôle clair et des tâches précises sans chevauchement."

**Résultat obtenu :** Répartition en 4 rôles distincts (Backend/BDD, Auth/Sécurité, Posts/Commentaires, Frontend/Tests), planning sprint semaine par semaine, matrice RACI complète et identification des dépendances entre membres (ex: M1 doit livrer le schéma BDD avant que M2 et M3 puissent commencer).

**Ce qu'on a adapté :** Nous avons ajusté la répartition des tâches selon les disponibilités réelles de chaque membre et revu les délais en fonction de notre niveau en Go.

---

## Prompt 2 — Conception du schéma entité-relation SQLite

**Contexte :** Avant de commencer à coder, il fallait concevoir un schéma BDD cohérent que tous les membres allaient utiliser. C'est la base sur laquelle tout le projet repose.

**Prompt utilisé :**
> "Je dois concevoir le schéma entité-relation SQLite pour un forum web en Go. Le forum doit gérer : des utilisateurs avec email unique et mot de passe hashé, des sessions de connexion avec expiration, des posts avec titre, contenu, image optionnelle et plusieurs catégories possibles, des commentaires sur les posts, un système de likes et dislikes sur les posts ET les commentaires, et des catégories prédéfinies. Génère le fichier schema.sql complet avec les CREATE TABLE, les contraintes NOT NULL, UNIQUE, CHECK, les clés étrangères avec ON DELETE CASCADE, et les données de base pour les catégories. Utilise des UUID en TEXT comme clés primaires."

**Résultat obtenu :** 7 tables (users, sessions, posts, comments, categories, post_categories, likes) avec toutes les contraintes, clés étrangères et suppressions en cascade. La table `likes` utilise un champ `target_type` pour gérer les likes sur posts ET commentaires avec une seule table, et un `CHECK (value IN (1, -1))` pour garantir l'intégrité des votes.

**Ce qu'on a adapté :** Nous avons validé le schéma en groupe avant de commencer à coder, et ajouté 5 catégories de base avec `INSERT OR IGNORE` pour éviter les doublons au redémarrage.

---

## Prompt 3 — Génération de la couche d'accès aux données (database/)

**Contexte :** Une fois le schéma validé, il fallait écrire toutes les fonctions Go qui communiquent avec SQLite. C'est le travail du Membre 1 qui conditionne le travail des autres membres.

**Prompt utilisé :**
> "Je développe un forum en Go pur avec SQLite. Mon schéma BDD contient 7 tables : users, sessions, posts, comments, categories, post_categories, likes. Je dois créer la couche d'accès aux données dans un dossier database/ avec un fichier par entité. Pour chaque fichier, génère les fonctions CRUD nécessaires en utilisant des requêtes préparées avec ? pour éviter les injections SQL. Voici les besoins précis : users.go doit avoir CreateUser, GetUserByEmail, GetUserByID. sessions.go doit avoir CreateSession qui supprime l'ancienne session (1 max par user), GetSessionByID avec vérification d'expiration automatique, DeleteSession, DeleteExpiredSessions. posts.go doit avoir CreatePost avec transaction pour les catégories, GetAllPosts avec les likes agrégés, GetPostsByCategory, GetPostsByUser, GetLikedPostsByUser, GetPostByID, UpdatePost, DeletePost, GetAllCategories. comments.go et likes.go avec ToggleLike qui gère automatiquement les 3 cas : INSERT si pas de vote, DELETE si même vote, UPDATE si vote opposé."

**Résultat obtenu :** 5 fichiers Go dans database/ avec toutes les fonctions, requêtes préparées, gestion des erreurs sql.ErrNoRows, et une fonction helper `getUserVote` réutilisée dans plusieurs requêtes. La fonction `ToggleLike` gère les 3 cas en une seule fonction.

**Ce qu'on a adapté :** Nous avons corrigé les imports manquants, ajusté les noms de champs pour correspondre exactement à nos structs Go, et testé chaque fonction en lançant le serveur et en vérifiant les logs.

---

## Prompt 4 — Mise en place du serveur HTTP, routing et middleware d'authentification

**Contexte :** Il fallait connecter toutes les couches ensemble : définir les routes HTTP, protéger les routes réservées aux connectés, et mettre en place un système de logging pour débugger pendant le développement.

**Prompt utilisé :**
> "Je développe un forum en Go pur sans framework. J'ai besoin de mettre en place le serveur HTTP complet dans main.go avec : toutes les routes publiques (/, /post/, /login, /register) et protégées (/post/create, /comment/create, /like, /logout, etc.), un middleware de logging qui affiche [METHOD] /chemin → statusCode (durée) dans le terminal pour chaque requête, un wrapper responseWriter pour capturer le status code HTTP, et des handlers 404/500 avec des pages HTML personnalisées. J'ai aussi besoin d'un middleware d'authentification dans middleware/auth.go avec trois fonctions : GetUserFromSession qui vérifie le cookie session_id et son expiration en BDD, RequireAuth qui redirige vers /login si non connecté, et WithUser qui injecte le user dans le contexte même pour les routes publiques pour que les templates sachent si quelqu'un est connecté."

**Résultat obtenu :** main.go avec routing complet, loggingMiddleware wrappant tout le serveur, middleware/auth.go avec les 3 fonctions utilisant `context.WithValue` pour passer le user aux handlers, et pages d'erreur 404/500 personnalisées.

**Ce qu'on a adapté :** Nous avons ajouté les routes de profil utilisateur (/profile, /profile/password, /profile/photo) qui n'étaient pas prévues initialement, et corrigé l'ordre des routes pour que `/post/` ne capture pas `/post/create`.

---

## Prompt 5 — Génération du frontend HTML/CSS style forum sportif

**Contexte :** Le Membre 4 devait créer tous les templates HTML en Go templates (html/template) et un CSS complet sans framework. Nous voulions un design moderne inspiré des plateformes sportives type Betclic, avec un thème sombre.

**Prompt utilisé :**
> "Je dois créer le frontend d'un forum web en HTML/CSS pur, sans React, Bootstrap ou autre framework. Les templates utilisent la syntaxe Go html/template avec {{range}}, {{if}}, {{.Variable}}. Je veux un design inspiré des plateformes de paris sportifs (style Betclic) avec un thème sombre : fond #0A0E1A (nuit de stade), accent vert #00FF87 style score en direct avec animation pulse, typographie Rajdhani pour les titres (condensé, sportif) et Inter pour le texte. Génère : style.css complet et responsive avec sidebar, cards de posts, boutons like/dislike, formulaires, pages auth centrées, pages d'erreur 404/500. Génère aussi index.html avec liste des posts et filtres par sidebar, post.html avec détail et commentaires, create_post.html avec upload image et checkboxes catégories, login.html et register.html avec carte centrée, 404.html et 500.html. Les templates doivent utiliser {{if .User}} pour l'affichage conditionnel guest vs connecté, et {{range .Posts}} pour boucler sur les données."

**Résultat obtenu :** CSS de 500+ lignes avec variables CSS, animations, responsive mobile, scrollbar personnalisée. 7 templates HTML complets avec syntaxe Go templates correcte, affichage conditionnel selon le rôle connecté/guest, boutons like/dislike avec états actifs colorés, et formulaire d'upload avec style drag-and-drop.

**Ce qu'on a adapté :** Nous avons ajusté les actions des formulaires pour correspondre exactement aux routes définies dans main.go, corrigé les noms des variables Go templates (`.Post.UserID` vs `.UserID`) et ajouté les messages d'erreur dans les formulaires d'auth via les paramètres URL.

---

## Bilan de l'utilisation de l'IA

**Points positifs :**
- Gain de temps considérable sur les parties répétitives (requêtes SQL, structs Go)
- L'IA a proposé des patterns que nous n'aurions pas pensé seuls (ToggleLike en une fonction, ON DELETE CASCADE)
- Aide précieuse pour débugger les erreurs de compilation Go

**Limites rencontrées et corrections apportées :**
- L'IA a généré des imports incorrects que nous avons dû corriger manuellement
- Les noms de variables dans les Go templates ne correspondaient pas toujours aux structs — nous avons dû faire correspondre les deux
- L'IA ne connaissait pas notre environnement exact (chemins de fichiers Windows) et nous avons dû adapter les commandes
- Certaines suggestions de l'IA utilisaient des patterns trop complexes que nous avons simplifiés pour mieux comprendre le code

**Conclusion :** L'IA a été utilisée comme un assistant technique, pas comme un remplaçant. Chaque membre a lu, compris et testé le code généré avant de l'intégrer au projet.
