package handlers

import (
	"net/http"
)

// ShowRegister affiche le formulaire d'inscription.
func ShowRegister(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/register.html")
}

// Register traite la soumission du formulaire d'inscription.
// TODO M2 : valider email unique, hasher mot de passe avec bcrypt, créer session
func Register(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Non implémenté", http.StatusNotImplemented)
}

// ShowLogin affiche le formulaire de connexion.
func ShowLogin(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/login.html")
}

// Login traite la soumission du formulaire de connexion.
// TODO M2 : vérifier credentials, créer cookie de session UUID
func Login(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Non implémenté", http.StatusNotImplemented)
}

// Logout supprime la session et le cookie.
// TODO M2 : DELETE session en BDD, expirer le cookie
func Logout(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
