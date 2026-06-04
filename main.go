package main

import (
	"forum/database"
	"forum/handlers"
	"forum/middleware"
	"log"
	"net/http"
)

func main() {
	database.Init("schema.sql")
	defer database.Close()

	mux := http.NewServeMux()

	// Fichiers statiques (CSS, images uploadées)
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

	// Gestion d'erreurs HTTP
	log.Println("Serveur démarré sur http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
