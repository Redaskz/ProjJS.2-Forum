package handlers

import (
	"net/http"
)

// ListPosts affiche tous les posts (page d'accueil).
// TODO M3 : SELECT posts avec filtres (catégorie, mes posts, posts aimés)
func ListPosts(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

// ShowPost affiche un post et ses commentaires.
// TODO M3 : SELECT post + comments + likes
func ShowPost(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/post.html")
}

// ShowCreatePost affiche le formulaire de création.
func ShowCreatePost(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/create_post.html")
}

// CreatePost traite la création d'un post.
// TODO M3 : valider, upload image, INSERT post + post_categories
func CreatePost(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Non implémenté", http.StatusNotImplemented)
}

// UpdatePost traite la modification d'un post.
// TODO M3 : vérifier propriétaire, UPDATE post
func UpdatePost(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Non implémenté", http.StatusNotImplemented)
}

// DeletePost supprime un post.
// TODO M3 : vérifier propriétaire, DELETE post
func DeletePost(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Non implémenté", http.StatusNotImplemented)
}
