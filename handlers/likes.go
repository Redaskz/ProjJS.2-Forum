package handlers

import (
	"forum/database"
	"forum/middleware"
	"forum/models"
	"net/http"
	"strconv"
)

func ToggleLike(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserKey).(*models.User)

	targetID := r.FormValue("target_id")
	targetType := r.FormValue("target_type")

	value, err := strconv.Atoi(r.FormValue("value"))
	if err != nil || (value != 1 && value != -1) {
		http.Error(w, "Vote invalide", http.StatusBadRequest)
		return
	}

	err = database.ToggleLike(
		user.ID,
		targetID,
		targetType,
		value,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	back := r.Referer()
	if back == "" {
		back = "/"
	}

	http.Redirect(w, r, back, http.StatusSeeOther)
}