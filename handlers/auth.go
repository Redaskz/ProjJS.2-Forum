package handlers

import (
	"forum/database"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// ShowRegister affiche le formulaire d'inscription.
func ShowRegister(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/register.html")
}


// Register traite la soumission du formulaire d'inscription.
func Register(w http.ResponseWriter, r *http.Request) {
	// Accepter uniquement les requêtes POST
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	// Récupérer les champs du formulaire
	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Vérifier que les champs ne sont pas vides
	if email == "" || username == "" || password == "" {
		http.Error(w, "Tous les champs sont obligatoires", http.StatusBadRequest)
		return
	}

	// Vérifier que l'email n'est pas déjà utilisé
	existingUser, err := database.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	if existingUser != nil {
		http.Error(w, "Cet email est déjà utilisé", http.StatusBadRequest)
		return
	}

	// Hasher le mot de passe avec bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// Créer l'utilisateur en BDD
	_, err = database.CreateUser(email, username, string(hashedPassword))
	if err != nil {
		http.Error(w, "Erreur lors de la création du compte", http.StatusInternalServerError)
		return
	}

	// Rediriger vers la page de connexion
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// ShowLogin affiche le formulaire de connexion.
func ShowLogin(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/login.html")
}

// Login traite la soumission du formulaire de connexion.
func Login(w http.ResponseWriter, r *http.Request) {
	// Accepter uniquement les requêtes POST
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Récupérer les champs du formulaire
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Vérifier que les champs ne sont pas vides
	if email == "" || password == "" {
		http.Error(w, "Tous les champs sont obligatoires", http.StatusBadRequest)
		return
	}

	// Chercher le user par email
	user, err := database.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// Si le user n'existe pas ou mot de passe incorrect → même message (sécurité)
	if user == nil {
		http.Error(w, "Email ou mot de passe incorrect", http.StatusUnauthorized)
		return
	}

	// Comparer le mot de passe avec le hash bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		http.Error(w, "Email ou mot de passe incorrect", http.StatusUnauthorized)
		return
	}

	// Créer une session en BDD (supprime l'ancienne automatiquement)
	session, err := database.CreateSession(user.ID)
	if err != nil {
		http.Error(w, "Erreur lors de la connexion", http.StatusInternalServerError)
		return
	}

	// Créer le cookie de session
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true, // JavaScript ne peut pas lire le cookie
		Path:     "/",
	})

	// Rediriger vers la page d'accueil
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout supprime la session et le cookie.
func Logout(w http.ResponseWriter, r *http.Request) {
	// Récupérer le cookie de session
	cookie, err := r.Cookie("session_id")
	if err == nil {
		// Supprimer la session en BDD
		database.DeleteSession(cookie.Value)
	}

	// Expirer le cookie côté navigateur
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour), // Date dans le passé = cookie supprimé
		HttpOnly: true,
		Path:     "/",
	})

	// Rediriger vers la page d'accueil
	http.Redirect(w, r, "/", http.StatusSeeOther)
}