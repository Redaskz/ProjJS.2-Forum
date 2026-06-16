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
)

func ListPosts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		http.ServeFile(w, r, "templates/404.html")
		return
	}

	var currentUser *models.User
	if u := r.Context().Value(middleware.UserKey); u != nil {
		currentUser = u.(*models.User)
	}

	userID := ""
	if currentUser != nil {
		userID = currentUser.ID
	}

	filter := r.URL.Query().Get("filter")

	var (
		posts []models.Post
		err   error
	)

	switch filter {
	case "category":
		posts, err = database.GetPostsByCategory(r.URL.Query().Get("id"), userID)
	case "myposts":
		if currentUser == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		posts, err = database.GetPostsByUser(currentUser.ID, userID)
	case "liked":
		if currentUser == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		posts, err = database.GetLikedPostsByUser(currentUser.ID)
	default:
		posts, err = database.GetAllPosts(userID)
	}

	if err != nil {
		serveInternalError(w)
		return
	}

	categories, _ := database.GetAllCategories()

	data := struct {
		User       *models.User
		Posts      []models.Post
		Categories []models.Category
		Filter     string
		FilterID   string
	}{
		User:       currentUser,
		Posts:      posts,
		Categories: categories,
		Filter:     filter,
		FilterID:   r.URL.Query().Get("id"),
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, data)
}

func ShowPost(w http.ResponseWriter, r *http.Request) {
	var currentUser *models.User
	if u := r.Context().Value(middleware.UserKey); u != nil {
		currentUser = u.(*models.User)
	}

	userID := ""
	if currentUser != nil {
		userID = currentUser.ID
	}

	post, err := database.GetPostByID(r.URL.Query().Get("id"), userID)
	if err != nil {
		serveInternalError(w)
		return
	}
	if post == nil {
		w.WriteHeader(http.StatusNotFound)
		http.ServeFile(w, r, "templates/404.html")
		return
	}

	comments, err := database.GetCommentsByPost(post.ID, userID)
	if err != nil {
		serveInternalError(w)
		return
	}

	categories, _ := database.GetAllCategories()

	data := struct {
		User       *models.User
		Post       *models.Post
		Comments   []models.Comment
		Categories []models.Category
	}{
		User:       currentUser,
		Post:       post,
		Comments:   comments,
		Categories: categories,
	}

	tmpl := template.Must(template.ParseFiles("templates/post.html"))
	tmpl.Execute(w, data)
}

func ShowCreatePost(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserKey).(*models.User)
	categories, _ := database.GetAllCategories()

	data := struct {
		User       *models.User
		Categories []models.Category
		Error      string
	}{
		User:       user,
		Categories: categories,
		Error:      r.URL.Query().Get("error"),
	}

	tmpl := template.Must(template.ParseFiles("templates/create_post.html"))
	tmpl.Execute(w, data)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserKey).(*models.User)

	r.ParseMultipartForm(20 << 20)

	imagePath := ""
	file, handler, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		ext := strings.ToLower(filepath.Ext(handler.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" {
			http.Redirect(w, r, "/post/create?error=format_non_supporte", http.StatusSeeOther)
			return
		}
		os.MkdirAll("uploads", 0755)
		filename := utils.NewID() + ext
		dst, err := os.Create("uploads/" + filename)
		if err == nil {
			io.Copy(dst, file)
			dst.Close()
			imagePath = filename
		}
	}

	_, err = database.CreatePost(user.ID, r.FormValue("title"), r.FormValue("content"), imagePath, r.Form["categories"])
	if err != nil {
		serveInternalError(w)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserKey).(*models.User)

	postID := r.FormValue("post_id")

	post, err := database.GetPostByID(postID, user.ID)
	if err != nil {
		serveInternalError(w)
		return
	}
	if post == nil {
		w.WriteHeader(http.StatusNotFound)
		http.ServeFile(w, r, "templates/404.html")
		return
	}
	if post.UserID != user.ID {
		http.Error(w, "Accès interdit", http.StatusForbidden)
		return
	}

	if err = database.UpdatePost(postID, r.FormValue("title"), r.FormValue("content")); err != nil {
		serveInternalError(w)
		return
	}

	http.Redirect(w, r, "/post/?id="+postID, http.StatusSeeOther)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserKey).(*models.User)

	postID := r.FormValue("post_id")

	post, err := database.GetPostByID(postID, user.ID)
	if err != nil {
		serveInternalError(w)
		return
	}
	if post == nil {
		w.WriteHeader(http.StatusNotFound)
		http.ServeFile(w, r, "templates/404.html")
		return
	}
	if post.UserID != user.ID && user.Role != "moderator" && user.Role != "admin" {
		http.Error(w, "Accès interdit", http.StatusForbidden)
		return
	}

	if err = database.DeletePost(postID); err != nil {
		serveInternalError(w)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
