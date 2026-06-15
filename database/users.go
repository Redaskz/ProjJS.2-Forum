package database

import (
	"database/sql"
	"forum/models"
	"forum/utils"
	"time"
)

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

func GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}

	err := DB.QueryRow(`
		SELECT id, email, username, password_hash, COALESCE(profile_photo, ''), created_at
		FROM users WHERE email = ?
	`, email).Scan(
		&user.ID, &user.Email, &user.Username,
		&user.PasswordHash, &user.ProfilePhoto, &user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByID(id string) (*models.User, error) {
	user := &models.User{}

	err := DB.QueryRow(`
		SELECT id, email, username, COALESCE(profile_photo, ''), created_at
		FROM users WHERE id = ?
	`, id).Scan(
		&user.ID, &user.Email, &user.Username,
		&user.ProfilePhoto, &user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func UpdatePassword(userID, hashedPassword string) error {
	_, err := DB.Exec(`UPDATE users SET password_hash = ? WHERE id = ?`, hashedPassword, userID)
	return err
}

func UpdateProfilePhoto(userID, photoPath string) error {
	_, err := DB.Exec(`UPDATE users SET profile_photo = ? WHERE id = ?`, photoPath, userID)
	return err
}
