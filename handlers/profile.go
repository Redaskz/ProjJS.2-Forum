package handlers

import (
	"forum/database"
	"forum/middleware"
	"forum/models"
	"forum/utils"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func ShowProfile(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserKey).(*models.User)

	fullUser, err := database.GetUserByID(user.ID)
	if err != nil || fullUser == nil {
		http.Error(w, "Erreur", http.StatusInternalServerError)
		return
	}

	data := struct {
		User    *models.User
		Success string
		Error   string
	}{
		User:    fullUser,
		Success: r.URL.Query().Get("success"),
		Error:   r.URL.Query().Get("error"),
	}

	tmpl := template.Must(template.ParseFiles("templates/profile.html"))
	tmpl.Execute(w, data)
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserKey).(*models.User)

	current := r.FormValue("current_password")
	newPass := r.FormValue("new_password")

	if current == "" || newPass == "" {
		http.Redirect(w, r, "/profile?error=champs_vides", http.StatusSeeOther)
		return
	}

	if len(newPass) < 6 {
		http.Redirect(w, r, "/profile?error=mot_de_passe_trop_court", http.StatusSeeOther)
		return
	}

	fullUser, err := database.GetUserByEmail(user.Email)
	if err != nil || fullUser == nil {
		http.Redirect(w, r, "/profile?error=erreur_serveur", http.StatusSeeOther)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(fullUser.PasswordHash), []byte(current)); err != nil {
		http.Redirect(w, r, "/profile?error=mot_de_passe_incorrect", http.StatusSeeOther)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
	if err != nil {
		http.Redirect(w, r, "/profile?error=erreur_serveur", http.StatusSeeOther)
		return
	}

	if err := database.UpdatePassword(user.ID, string(hashed)); err != nil {
		http.Redirect(w, r, "/profile?error=erreur_serveur", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/profile?success=mot_de_passe_modifie", http.StatusSeeOther)
}

func UploadProfilePhoto(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserKey).(*models.User)

	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("photo")
	if err != nil {
		http.Redirect(w, r, "/profile?error=fichier_invalide", http.StatusSeeOther)
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(handler.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" {
		http.Redirect(w, r, "/profile?error=format_non_supporte", http.StatusSeeOther)
		return
	}

	os.MkdirAll("uploads", 0755)
	filename := utils.NewID() + ext

	dst, err := os.Create("uploads/" + filename)
	if err != nil {
		http.Redirect(w, r, "/profile?error=erreur_serveur", http.StatusSeeOther)
		return
	}
	defer dst.Close()
	io.Copy(dst, file)

	database.UpdateProfilePhoto(user.ID, filename)

	http.Redirect(w, r, "/profile?success=photo_mise_a_jour", http.StatusSeeOther)
}
