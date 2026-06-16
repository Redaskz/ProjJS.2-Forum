package handlers

import (
	"net/http"
	"os"
)

// serveInternalError affiche le template 500.html avec le bon code HTTP.
// On lit le fichier directement pour éviter de dépendre d'une requête.
func serveInternalError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	content, err := os.ReadFile("templates/500.html")
	if err != nil {
		http.Error(w, "Erreur interne", http.StatusInternalServerError)
		return
	}
	w.Write(content)
}
