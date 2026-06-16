package handlers

import (
	"forum/database"
	"forum/middleware"
	"forum/models"
	"net/http"
)

func CreateComment(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserKey).(*models.User)

	postID := r.FormValue("post_id")
	content := r.FormValue("content")

	if content == "" {
		http.Error(w, "Commentaire vide", http.StatusBadRequest)
		return
	}

	if _, err := database.CreateComment(postID, user.ID, content); err != nil {
		serveInternalError(w)
		return
	}

	http.Redirect(w, r, "/post/?id="+postID, http.StatusSeeOther)
}

func UpdateComment(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserKey).(*models.User)

	commentID := r.FormValue("comment_id")

	comment, err := database.GetCommentByID(commentID)
	if err != nil {
		serveInternalError(w)
		return
	}
	if comment == nil {
		w.WriteHeader(http.StatusNotFound)
		http.ServeFile(w, r, "templates/404.html")
		return
	}
	if comment.UserID != user.ID {
		http.Error(w, "Accès interdit", http.StatusForbidden)
		return
	}

	if err = database.UpdateComment(commentID, r.FormValue("content")); err != nil {
		serveInternalError(w)
		return
	}

	http.Redirect(w, r, "/post/?id="+comment.PostID, http.StatusSeeOther)
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserKey).(*models.User)

	commentID := r.FormValue("comment_id")

	comment, err := database.GetCommentByID(commentID)
	if err != nil {
		serveInternalError(w)
		return
	}
	if comment == nil {
		w.WriteHeader(http.StatusNotFound)
		http.ServeFile(w, r, "templates/404.html")
		return
	}
	if comment.UserID != user.ID && user.Role != "moderator" && user.Role != "admin" {
		http.Error(w, "Accès interdit", http.StatusForbidden)
		return
	}

	if err = database.DeleteComment(commentID); err != nil {
		serveInternalError(w)
		return
	}

	http.Redirect(w, r, "/post/?id="+comment.PostID, http.StatusSeeOther)
}
