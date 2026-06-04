package handlers

import (
	"net/http"
)

// ToggleLike gère le like/dislike sur un post ou commentaire.
// target_type : "post" ou "comment"
// value       : 1 (like) ou -1 (dislike)
// TODO M3 : INSERT ou UPDATE ou DELETE dans likes selon l'état actuel
func ToggleLike(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Non implémenté", http.StatusNotImplemented)
}
