package main

import (
	"forum/database"
	"forum/handlers"
	"forum/middleware"
	"log"
	"net/http"
	"time"
)

// loggingMiddleware affiche dans le terminal chaque requête reçue.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)
		log.Printf("[%s] %s → %d (%s)", r.Method, r.URL.Path, rw.statusCode, time.Since(start))
	})
}

// responseWriter wrapper pour capturer le status code HTTP.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// notFoundHandler gère les routes inconnues (404).
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	http.ServeFile(w, r, "templates/404.html")
}

// internalErrorHandler gère les erreurs internes (500).
func internalErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	http.ServeFile(w, r, "templates/500.html")
}

func main() {
	database.Init("schema.sql")
	defer database.Close()

	mux := http.NewServeMux()

	// Fichiers statiques
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	// Routes publiques
	mux.HandleFunc("/", middleware.WithUser(handlers.ListPosts))
	mux.HandleFunc("/post/", middleware.WithUser(handlers.ShowPost))
	mux.HandleFunc("/login", handlers.ShowLogin)
	mux.HandleFunc("/login/submit", handlers.Login)
	mux.HandleFunc("/register", handlers.ShowRegister)
	mux.HandleFunc("/register/submit", handlers.Register)

	// Routes protégées (connecté uniquement)
	mux.HandleFunc("/logout", middleware.RequireAuth(handlers.Logout))
	mux.HandleFunc("/post/create", middleware.RequireAuth(handlers.ShowCreatePost))
	mux.HandleFunc("/post/create/submit", middleware.RequireAuth(handlers.CreatePost))
	mux.HandleFunc("/post/update", middleware.RequireAuth(handlers.UpdatePost))
	mux.HandleFunc("/post/delete", middleware.RequireAuth(handlers.DeletePost))
	mux.HandleFunc("/comment/create", middleware.RequireAuth(handlers.CreateComment))
	mux.HandleFunc("/comment/update", middleware.RequireAuth(handlers.UpdateComment))
	mux.HandleFunc("/comment/delete", middleware.RequireAuth(handlers.DeleteComment))
	mux.HandleFunc("/like", middleware.RequireAuth(handlers.ToggleLike))

	// Profil utilisateur
	mux.HandleFunc("/profile", middleware.RequireAuth(handlers.ShowProfile))
	mux.HandleFunc("/profile/password", middleware.RequireAuth(handlers.ChangePassword))
	mux.HandleFunc("/profile/photo", middleware.RequireAuth(handlers.UploadProfilePhoto))

	// Calendrier Coupe du Monde 2026
	mux.HandleFunc("/calendrier", middleware.WithUser(handlers.ShowCalendar))

	// Administration (admin uniquement)
	mux.HandleFunc("/admin", middleware.RequireAdmin(handlers.ShowAdmin))
	mux.HandleFunc("/admin/role", middleware.RequireAdmin(handlers.ChangeUserRole))

	// Pages d'erreur
	mux.HandleFunc("/404", notFoundHandler)
	mux.HandleFunc("/500", internalErrorHandler)

	// Logging sur toutes les requêtes
	loggedMux := loggingMiddleware(mux)

	log.Println("Serveur démarré sur http://localhost:8080")
	if err := http.ListenAndServe(":8080", loggedMux); err != nil {
		log.Fatal(err)
	}
}