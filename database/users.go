package database

import (
	"database/sql"
	"forum/models"
	"forum/utils"
	"time"
)

// CreateUser insère un nouvel utilisateur en BDD.
// Le mot de passe doit déjà être hashé avec bcrypt avant d'appeler cette fonction.
func CreateUser(email, username, passwordHash string) (*models.User, error) {
	user := &models.User{
		ID:           utils.NewID(),
		Email:        email,
		Username:     username,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}

	_, err := DB.Exec(`
		INSERT INTO users (id, email, username, password_hash, created_at)
		VALUES (?, ?, ?, ?, ?)
	`, user.ID, user.Email, user.Username, user.PasswordHash, user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByEmail récupère un utilisateur par son email.
// Utilisé lors de la connexion pour vérifier les credentials.
// Retourne sql.ErrNoRows si l'email n'existe pas.
func GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}

	err := DB.QueryRow(`
		SELECT id, email, username, password_hash, created_at
		FROM users
		WHERE email = ?
	`, email).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID récupère un utilisateur par son ID.
// Utilisé par le middleware de session.
func GetUserByID(id string) (*models.User, error) {
	user := &models.User{}

	err := DB.QueryRow(`
		SELECT id, email, username, created_at
		FROM users
		WHERE id = ?
	`, id).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}