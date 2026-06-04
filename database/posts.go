package database

import (
	"database/sql"
	"forum/models"
	"forum/utils"
	"time"
)

// CreatePost insère un nouveau post avec ses catégories.
func CreatePost(userID, title, content, imagePath string, categoryIDs []string) (*models.Post, error) {
	post := &models.Post{
		ID:        utils.NewID(),
		UserID:    userID,
		Title:     title,
		Content:   content,
		ImagePath: imagePath,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tx, err := DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO posts (id, user_id, title, content, image_path, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, post.ID, post.UserID, post.Title, post.Content, post.ImagePath, post.CreatedAt, post.UpdatedAt)
	if err != nil {
		return nil, err
	}

	for _, catID := range categoryIDs {
		_, err = tx.Exec(`
			INSERT INTO post_categories (post_id, category_id)
			VALUES (?, ?)
		`, post.ID, catID)
		if err != nil {
			return nil, err
		}
	}

	return post, tx.Commit()
}

// GetAllPosts récupère tous les posts avec leur auteur et leurs catégories.
func GetAllPosts(userID string) ([]models.Post, error) {
	rows, err := DB.Query(`
		SELECT p.id, p.user_id, u.username, p.title, p.content, p.image_path, p.created_at,
			COALESCE(SUM(CASE WHEN l.value = 1 THEN 1 ELSE 0 END), 0) as likes,
			COALESCE(SUM(CASE WHEN l.value = -1 THEN 1 ELSE 0 END), 0) as dislikes
		FROM posts p
		JOIN users u ON u.id = p.user_id
		LEFT JOIN likes l ON l.target_id = p.id AND l.target_type = 'post'
		GROUP BY p.id
		ORDER BY p.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanPosts(rows, userID)
}

// GetPostsByCategory récupère les posts d'une catégorie donnée.
func GetPostsByCategory(categoryID, userID string) ([]models.Post, error) {
	rows, err := DB.Query(`
		SELECT p.id, p.user_id, u.username, p.title, p.content, p.image_path, p.created_at,
			COALESCE(SUM(CASE WHEN l.value = 1 THEN 1 ELSE 0 END), 0) as likes,
			COALESCE(SUM(CASE WHEN l.value = -1 THEN 1 ELSE 0 END), 0) as dislikes
		FROM posts p
		JOIN users u ON u.id = p.user_id
		JOIN post_categories pc ON pc.post_id = p.id
		LEFT JOIN likes l ON l.target_id = p.id AND l.target_type = 'post'
		WHERE pc.category_id = ?
		GROUP BY p.id
		ORDER BY p.created_at DESC
	`, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanPosts(rows, userID)
}

// GetPostsByUser récupère les posts créés par un utilisateur.
func GetPostsByUser(authorID, userID string) ([]models.Post, error) {
	rows, err := DB.Query(`
		SELECT p.id, p.user_id, u.username, p.title, p.content, p.image_path, p.created_at,
			COALESCE(SUM(CASE WHEN l.value = 1 THEN 1 ELSE 0 END), 0) as likes,
			COALESCE(SUM(CASE WHEN l.value = -1 THEN 1 ELSE 0 END), 0) as dislikes
		FROM posts p
		JOIN users u ON u.id = p.user_id
		LEFT JOIN likes l ON l.target_id = p.id AND l.target_type = 'post'
		WHERE p.user_id = ?
		GROUP BY p.id
		ORDER BY p.created_at DESC
	`, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanPosts(rows, userID)
}

// GetLikedPostsByUser récupère les posts aimés par un utilisateur.
func GetLikedPostsByUser(userID string) ([]models.Post, error) {
	rows, err := DB.Query(`
		SELECT p.id, p.user_id, u.username, p.title, p.content, p.image_path, p.created_at,
			COALESCE(SUM(CASE WHEN l.value = 1 THEN 1 ELSE 0 END), 0) as likes,
			COALESCE(SUM(CASE WHEN l.value = -1 THEN 1 ELSE 0 END), 0) as dislikes
		FROM posts p
		JOIN users u ON u.id = p.user_id
		LEFT JOIN likes l ON l.target_id = p.id AND l.target_type = 'post'
		WHERE EXISTS (
			SELECT 1 FROM likes
			WHERE user_id = ? AND target_id = p.id AND target_type = 'post' AND value = 1
		)
		GROUP BY p.id
		ORDER BY p.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanPosts(rows, userID)
}

// GetPostByID récupère un post par son ID.
func GetPostByID(postID, userID string) (*models.Post, error) {
	post := &models.Post{}

	err := DB.QueryRow(`
		SELECT p.id, p.user_id, u.username, p.title, p.content, p.image_path, p.created_at,
			COALESCE(SUM(CASE WHEN l.value = 1 THEN 1 ELSE 0 END), 0) as likes,
			COALESCE(SUM(CASE WHEN l.value = -1 THEN 1 ELSE 0 END), 0) as dislikes
		FROM posts p
		JOIN users u ON u.id = p.user_id
		LEFT JOIN likes l ON l.target_id = p.id AND l.target_type = 'post'
		WHERE p.id = ?
		GROUP BY p.id
	`, postID).Scan(
		&post.ID, &post.UserID, &post.Username,
		&post.Title, &post.Content, &post.ImagePath,
		&post.CreatedAt, &post.Likes, &post.Dislikes,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	post.Categories, _ = GetCategoriesByPost(postID)
	post.UserVote = getUserVote(userID, postID, "post")

	return post, nil
}

// UpdatePost modifie le titre et le contenu d'un post.
func UpdatePost(postID, title, content string) error {
	_, err := DB.Exec(`
		UPDATE posts SET title = ?, content = ?, updated_at = ?
		WHERE id = ?
	`, title, content, time.Now(), postID)
	return err
}

// DeletePost supprime un post et tout ce qui y est lié (cascade).
func DeletePost(postID string) error {
	_, err := DB.Exec(`DELETE FROM posts WHERE id = ?`, postID)
	return err
}

// GetCategoriesByPost récupère les catégories d'un post.
func GetCategoriesByPost(postID string) ([]models.Category, error) {
	rows, err := DB.Query(`
		SELECT c.id, c.name
		FROM categories c
		JOIN post_categories pc ON pc.category_id = c.id
		WHERE pc.post_id = ?
	`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	return categories, nil
}

// GetAllCategories récupère toutes les catégories disponibles.
func GetAllCategories() ([]models.Category, error) {
	rows, err := DB.Query(`SELECT id, name FROM categories ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	return categories, nil
}

// getUserVote retourne le vote d'un utilisateur sur une cible (1, -1, ou 0).
func getUserVote(userID, targetID, targetType string) int {
	if userID == "" {
		return 0
	}
	var value int
	err := DB.QueryRow(`
		SELECT value FROM likes
		WHERE user_id = ? AND target_id = ? AND target_type = ?
	`, userID, targetID, targetType).Scan(&value)
	if err != nil {
		return 0
	}
	return value
}

// scanPosts est un helper qui lit les rows SQL et retourne une liste de posts.
func scanPosts(rows *sql.Rows, userID string) ([]models.Post, error) {
	var posts []models.Post
	for rows.Next() {
		var p models.Post
		err := rows.Scan(
			&p.ID, &p.UserID, &p.Username,
			&p.Title, &p.Content, &p.ImagePath,
			&p.CreatedAt, &p.Likes, &p.Dislikes,
		)
		if err != nil {
			return nil, err
		}
		p.Categories, _ = GetCategoriesByPost(p.ID)
		p.UserVote = getUserVote(userID, p.ID, "post")
		posts = append(posts, p)
	}
	return posts, nil
}