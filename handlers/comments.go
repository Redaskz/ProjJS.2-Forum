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

	_, err := database.CreateComment(
		postID,
		user.ID,
		content,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/post/?id="+postID, http.StatusSeeOther)
}

func UpdateComment(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserKey).(*models.User)

	commentID := r.FormValue("comment_id")
	content := r.FormValue("content")

	comment, err := database.GetCommentByID(commentID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if comment == nil {
		http.NotFound(w, r)
		return
	}

	if comment.UserID != user.ID {
		http.Error(w, "Accès interdit", http.StatusForbidden)
		return
	}

	err = database.UpdateComment(
		commentID,
		content,
	)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/post/?id="+comment.PostID, http.StatusSeeOther)
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserKey).(*models.User)

	commentID := r.FormValue("comment_id")

	comment, err := database.GetCommentByID(commentID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if comment == nil {
		http.NotFound(w, r)
		return
	}

	if comment.UserID != user.ID {
		http.Error(w, "Accès interdit", http.StatusForbidden)
		return
	}

	err = database.DeleteComment(commentID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/post/?id="+comment.PostID, http.StatusSeeOther)
}