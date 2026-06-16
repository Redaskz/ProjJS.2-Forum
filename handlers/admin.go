package handlers

import (
	"forum/database"
	"forum/middleware"
	"forum/models"
	"html/template"
	"net/http"
)

func ShowAdmin(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserKey).(*models.User)

	users, err := database.GetAllUsers()
	if err != nil {
		serveInternalError(w)
		return
	}

	data := struct {
		User    *models.User
		Users   []models.User
		Success string
		Error   string
	}{
		User:    user,
		Users:   users,
		Success: r.URL.Query().Get("success"),
		Error:   r.URL.Query().Get("error"),
	}

	tmpl := template.Must(template.ParseFiles("templates/admin.html"))
	tmpl.Execute(w, data)
}

func ChangeUserRole(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	currentUser := r.Context().Value(middleware.UserKey).(*models.User)
	targetID := r.FormValue("user_id")
	role := r.FormValue("role")

	validRoles := map[string]bool{"guest": true, "user": true, "moderator": true, "admin": true}
	if !validRoles[role] {
		http.Error(w, "Rôle invalide", http.StatusBadRequest)
		return
	}

	if targetID == currentUser.ID {
		http.Redirect(w, r, "/admin?error=self", http.StatusSeeOther)
		return
	}

	if err := database.UpdateUserRole(targetID, role); err != nil {
		serveInternalError(w)
		return
	}

	http.Redirect(w, r, "/admin?success=1", http.StatusSeeOther)
}
