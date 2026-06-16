package handlers

import (
	"forum/middleware"
	"forum/models"
	"html/template"
	"net/http"
)

func ShowCalendar(w http.ResponseWriter, r *http.Request) {
	var currentUser *models.User
	if u := r.Context().Value(middleware.UserKey); u != nil {
		currentUser = u.(*models.User)
	}
	tmpl := template.Must(template.ParseFiles("templates/calendrier.html"))
	tmpl.Execute(w, struct{ User *models.User }{User: currentUser})
}
