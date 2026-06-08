package handlers

import (
	"forum/database"
	"forum/middleware"
	"forum/models"
	"html/template"
	"net/http"
)

func ListPosts(w http.ResponseWriter, r *http.Request) {
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
		categoryID := r.URL.Query().Get("id")
		posts, err = database.GetPostsByCategory(categoryID, userID)

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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	categories, _ := database.GetAllCategories()

	data := struct {
		User       *models.User
		Posts      []models.Post
		Categories []models.Category
	}{
		User:       currentUser,
		Posts:      posts,
		Categories: categories,
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, data)
}

func ShowPost(w http.ResponseWriter, r *http.Request) {
	var currentUser *models.User

	if u := r.Context().Value(middleware.UserKey); u != nil {
		currentUser = u.(*models.User)
	}

	postID := r.URL.Query().Get("id")

	userID := ""
	if currentUser != nil {
		userID = currentUser.ID
	}

	post, err := database.GetPostByID(postID, userID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if post == nil {
		http.NotFound(w, r)
		return
	}

	comments, err := database.GetCommentsByPost(postID, userID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	data := struct {
		User     *models.User
		Post     *models.Post
		Comments []models.Comment
	}{
		User:     currentUser,
		Post:     post,
		Comments: comments,
	}

	tmpl := template.Must(template.ParseFiles("templates/post.html"))
	tmpl.Execute(w, data)
}

func ShowCreatePost(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/create_post.html")
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserKey).(*models.User)

	title := r.FormValue("title")
	content := r.FormValue("content")
	categoryIDs := r.Form["categories"]

	_, err := database.CreatePost(
		user.ID,
		title,
		content,
		"",
		categoryIDs,
	)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserKey).(*models.User)

	postID := r.FormValue("post_id")

	post, err := database.GetPostByID(postID, user.ID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if post == nil {
		http.NotFound(w, r)
		return
	}

	if post.UserID != user.ID {
		http.Error(w, "Accès interdit", http.StatusForbidden)
		return
	}

	err = database.UpdatePost(
		postID,
		r.FormValue("title"),
		r.FormValue("content"),
	)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/post/?id="+postID, http.StatusSeeOther)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserKey).(*models.User)

	postID := r.FormValue("post_id")

	post, err := database.GetPostByID(postID, user.ID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if post == nil {
		http.NotFound(w, r)
		return
	}

	if post.UserID != user.ID {
		http.Error(w, "Accès interdit", http.StatusForbidden)
		return
	}

	err = database.DeletePost(postID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}