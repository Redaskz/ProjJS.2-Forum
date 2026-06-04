package handlers

import (
	"net/http"
)

// CreateComment crée un commentaire sur un post.
// TODO M3 : INSERT comment
func CreateComment(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Non implémenté", http.StatusNotImplemented)
}

// UpdateComment modifie un commentaire.
// TODO M3 : vérifier propriétaire, UPDATE comment
func UpdateComment(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Non implémenté", http.StatusNotImplemented)
}

// DeleteComment supprime un commentaire.
// TODO M3 : vérifier propriétaire, DELETE comment
func DeleteComment(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Non implémenté", http.StatusNotImplemented)
}
