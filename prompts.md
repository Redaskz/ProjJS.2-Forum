# Journal de bord — Utilisation de l'IA (Claude)

Ce fichier documente les prompts utilisés avec Claude (claude.ai / Claude Code) durant le développement du projet Forum Foot.

---

## Prompt 1 — Analyse du projet et état des lieux

**Contexte :** Début du projet, besoin de comprendre ce qui était déjà fait par les autres membres.

**Prompt :**
> "Voici le projet de forum en Go. Analyse le code existant et dis-moi ce qui est fait, ce qui manque, et ce qui est cassé. Ne te base pas sur ma description, vérifie toi-même les fichiers."

**Résultat obtenu :**
Claude a lu chaque fichier (handlers/, database/, models/, main.go) et a produit un audit complet : le module `golang.org/x/crypto` manquait dans go.mod malgré son utilisation dans auth.go. Il a aussi listé les templates manquants.

**Itération :** Correction manuelle du go.mod confirmée par `go build`.

---

## Prompt 2 — Ajout de la dépendance bcrypt manquante

**Contexte :** Le projet ne compilait pas car `golang.org/x/crypto` n'était pas dans go.mod.

**Prompt :**
> "Le module bcrypt est utilisé dans handlers/auth.go mais absent de go.mod. Ajoute-le correctement sans casser le reste."

**Résultat obtenu :**
Claude a exécuté `go get golang.org/x/crypto` et mis à jour go.mod/go.sum. Build confirmé sans erreur.

**Correction apportée :** Aucune — la suggestion était correcte.

---

## Prompt 3 — Création des templates HTML/CSS

**Contexte :** Aucun template fonctionnel n'existait (fichiers statiques sans directives Go template).

**Prompt :**
> "Crée tous les templates HTML/CSS pour le forum. Inspiré de betclic.fr : 3 colonnes (matchs WC2026 à gauche, posts au centre, actu foot à droite). HTML/CSS pur uniquement, pas de framework. Couleurs claires, base blanche, rouge comme couleur principale."

**Résultat obtenu :**
Claude a généré index.html, post.html, create_post.html, login.html, register.html avec un layout 3 colonnes CSS Grid, variables CSS, effets hover, responsive.

**Correction apportée :** Première version trop sombre → demande explicite de passer en thème clair avec base blanche. Deuxième itération correcte.

---

## Prompt 4 — Page profil et upload de photo

**Contexte :** Besoin d'une page profil permettant de changer le mot de passe et uploader une photo.

**Prompt :**
> "Il faut une page profil pour l'utilisateur : changer son mot de passe (vérification bcrypt du mot de passe actuel), uploader une photo de profil (jpg/png/gif/webp, 10Mo max). La photo doit remplacer l'avatar initiale partout sur le site. Respecte les contraintes : Go pur, pas de framework, SQLite via go-sqlite3."

**Résultat obtenu :**
Claude a créé : handlers/profile.go (3 handlers), migration SQLite silencieuse pour la colonne profile_photo, mise à jour de tous les SELECT users avec COALESCE, routes dans main.go, template profile.html complet.

**Limite détectée :** Claude avait oublié de mettre à jour la requête dans middleware/auth.go → correction immédiate après signalement.

---

## Prompt 5 — Authentification : accès différencié guest/user

**Contexte :** Un utilisateur non connecté pouvait accéder aux routes de création.

**Prompt :**
> "Assure-toi que seuls les utilisateurs connectés peuvent créer, modifier, supprimer des posts et des commentaires, et voter. Dans les templates, si l'utilisateur n'est pas connecté, affiche un message l'invitant à se connecter plutôt que de lui donner le formulaire."

**Résultat obtenu :**
Les templates ont été mis à jour avec `{{if .User}}` / `{{else}}` pour les formulaires de vote, commentaire et modification. Le middleware RequireAuth était déjà en place côté backend.

**Regard critique :** Claude a parfois généré du code trop verbeux pour les templates. J'ai simplifié certaines conditions redondantes.

---

## Prompt 6 — Upload d'images dans les posts

**Contexte :** La fonction CreatePost passait une chaîne vide comme imagePath, l'upload n'était pas implémenté.

**Prompt :**
> "Dans handlers/posts.go, la fonction CreatePost passe toujours '' comme imagePath. Implémente l'upload d'image optionnel : validation d'extension (jpg/jpeg/png/gif/webp), nom UUID, sauvegarde dans uploads/."

**Résultat obtenu :**
Claude a ajouté les imports nécessaires (io, os, path/filepath, strings, utils) et implémenté la logique d'upload avec gestion silencieuse des erreurs pour garder l'upload optionnel.

---

## Bilan critique de l'utilisation de l'IA

**Ce qui a bien fonctionné :**
- Génération rapide de code boilerplate (templates, handlers CRUD)
- Identification de bugs (module manquant, requête SQL incomplète)
- Cohérence entre les fichiers (mêmes variables CSS partout)

**Limites rencontrées :**
- Claude a oublié de mettre à jour middleware/auth.go lors de l'ajout de profile_photo → bug silencieux
- Les premières versions CSS étaient en thème sombre alors que la demande était thème clair → 2 itérations nécessaires
- Claude propose parfois des abstractions inutiles (helper functions) → refus et demande de rester simple
- Il faut toujours vérifier que le code compilé fonctionne réellement, pas juste que Claude dit "c'est bon"
