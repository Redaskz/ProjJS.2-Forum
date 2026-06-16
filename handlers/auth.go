package handlers

import (
	"forum/database"
	"html/template"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func ShowRegister(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/register.html"))
	tmpl.Execute(w, struct{ Error string }{Error: r.URL.Query().Get("error")})
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")

	if email == "" || username == "" || password == "" {
		http.Error(w, "Tous les champs sont obligatoires", http.StatusBadRequest)
		return
	}

	existingUser, err := database.GetUserByEmail(email)
	if err != nil {
		http.Redirect(w, r, "/register?error=erreur_serveur", http.StatusSeeOther)
		return
	}
	if existingUser != nil {
		http.Redirect(w, r, "/register?error=email_deja_utilise", http.StatusSeeOther)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		serveInternalError(w)
		return
	}

	_, err = database.CreateUser(email, username, string(hashedPassword))
	if err != nil {
		serveInternalError(w)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func ShowLogin(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/login.html"))
	tmpl.Execute(w, struct{ Error string }{Error: r.URL.Query().Get("error")})
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		http.Error(w, "Tous les champs sont obligatoires", http.StatusBadRequest)
		return
	}

	user, err := database.GetUserByEmail(email)
	if err != nil {
		serveInternalError(w)
		return
	}

	// Même message pour utilisateur inconnu et mauvais mot de passe (sécurité)
	if user == nil {
		http.Redirect(w, r, "/login?error=identifiants_incorrects", http.StatusSeeOther)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		http.Redirect(w, r, "/login?error=identifiants_incorrects", http.StatusSeeOther)
		return
	}

	session, err := database.CreateSession(user.ID)
	if err != nil {
		serveInternalError(w)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		database.DeleteSession(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
