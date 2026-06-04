package database

import (
	"database/sql"
	"forum/models"
	"forum/utils"
	"time"
)

// CreateComment insère un nouveau commentaire sur un post.
func CreateComment(postID, userID, content string) (*models.Comment, error) {
	comment := &models.Comment{
		ID:        utils.NewID(),
		PostID:    postID,
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := DB.Exec(`
		INSERT INTO comments (id, post_id, user_id, content, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, comment.ID, comment.PostID, comment.UserID, comment.Content, comment.CreatedAt, comment.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return comment, nil
}

// GetCommentsByPost récupère tous les commentaires d'un post.
func GetCommentsByPost(postID, userID string) ([]models.Comment, error) {
	rows, err := DB.Query(`
		SELECT c.id, c.post_id, c.user_id, u.username, c.content, c.created_at,
			COALESCE(SUM(CASE WHEN l.value = 1 THEN 1 ELSE 0 END), 0) as likes,
			COALESCE(SUM(CASE WHEN l.value = -1 THEN 1 ELSE 0 END), 0) as dislikes
		FROM comments c
		JOIN users u ON u.id = c.user_id
		LEFT JOIN likes l ON l.target_id = c.id AND l.target_type = 'comment'
		WHERE c.post_id = ?
		GROUP BY c.id
		ORDER BY c.created_at ASC
	`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		err := rows.Scan(
			&c.ID, &c.PostID, &c.UserID, &c.Username,
			&c.Content, &c.CreatedAt, &c.Likes, &c.Dislikes,
		)
		if err != nil {
			return nil, err
		}
		c.UserVote = getUserVote(userID, c.ID, "comment")
		comments = append(comments, c)
	}
	return comments, nil
}

// GetCommentByID récupère un commentaire par son ID.
func GetCommentByID(commentID string) (*models.Comment, error) {
	comment := &models.Comment{}

	err := DB.QueryRow(`
		SELECT id, post_id, user_id, content, created_at, updated_at
		FROM comments WHERE id = ?
	`, commentID).Scan(
		&comment.ID, &comment.PostID, &comment.UserID,
		&comment.Content, &comment.CreatedAt, &comment.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return comment, nil
}

// UpdateComment modifie le contenu d'un commentaire.
func UpdateComment(commentID, content string) error {
	_, err := DB.Exec(`
		UPDATE comments SET content = ?, updated_at = ? WHERE id = ?
	`, content, time.Now(), commentID)
	return err
}

// DeleteComment supprime un commentaire.
func DeleteComment(commentID string) error {
	_, err := DB.Exec(`DELETE FROM comments WHERE id = ?`, commentID)
	return err
}