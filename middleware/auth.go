package middleware

import (
	"context"
	"forum/database"
	"forum/models"
	"net/http"
	"time"
)

type contextKey string

const UserKey contextKey = "user"

// GetUserFromSession récupère l'utilisateur connecté depuis le cookie de session.
// Retourne nil si non connecté ou session expirée.
func GetUserFromSession(r *http.Request) *models.User {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil
	}

	var user models.User
	var expiresAt time.Time

	err = database.DB.QueryRow(`
		SELECT u.id, u.email, u.username, s.expires_at
		FROM sessions s
		JOIN users u ON u.id = s.user_id
		WHERE s.id = ?
	`, cookie.Value).Scan(&user.ID, &user.Email, &user.Username, &expiresAt)

	if err != nil {
		return nil
	}

	if time.Now().After(expiresAt) {
		database.DB.Exec("DELETE FROM sessions WHERE id = ?", cookie.Value)
		return nil
	}

	return &user
}

// RequireAuth redirige vers /login si l'utilisateur n'est pas connecté.
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := GetUserFromSession(r)
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), UserKey, user)
		next(w, r.WithContext(ctx))
	}
}

// WithUser injecte l'utilisateur dans le contexte (pour les routes publiques).
func WithUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := GetUserFromSession(r)
		ctx := context.WithValue(r.Context(), UserKey, user)
		next(w, r.WithContext(ctx))
	}
}
