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
		SELECT id, email, username, password_hash, COALESCE(role,'user'), COALESCE(profile_photo, ''), created_at
		FROM users WHERE email = ?
	`, email).Scan(
		&user.ID, &user.Email, &user.Username,
		&user.PasswordHash, &user.Role, &user.ProfilePhoto, &user.CreatedAt,
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
		SELECT id, email, username, COALESCE(role,'user'), COALESCE(profile_photo, ''), created_at
		FROM users WHERE id = ?
	`, id).Scan(
		&user.ID, &user.Email, &user.Username,
		&user.Role, &user.ProfilePhoto, &user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetAllUsers() ([]models.User, error) {
	rows, err := DB.Query(`
		SELECT id, email, username, COALESCE(role,'user'), COALESCE(profile_photo,''), created_at
		FROM users ORDER BY created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Username, &u.Role, &u.ProfilePhoto, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func UpdateUserRole(userID, role string) error {
	_, err := DB.Exec(`UPDATE users SET role = ? WHERE id = ?`, role, userID)
	return err
}

func UpdatePassword(userID, hashedPassword string) error {
	_, err := DB.Exec(`UPDATE users SET password_hash = ? WHERE id = ?`, hashedPassword, userID)
	return err
}

func UpdateProfilePhoto(userID, photoPath string) error {
	_, err := DB.Exec(`UPDATE users SET profile_photo = ? WHERE id = ?`, photoPath, userID)
	return err
}
