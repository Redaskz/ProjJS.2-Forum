package database

import (
	"database/sql"
	"forum/models"
	"forum/utils"
	"time"
)

// CreateSession crée une nouvelle session pour un utilisateur.
// Supprime d'abord toute session existante (1 session max par user).
// La session expire après 24 heures.
func CreateSession(userID string) (*models.Session, error) {
	// Supprimer l'ancienne session si elle existe
	_, err := DB.Exec(`DELETE FROM sessions WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}

	session := &models.Session{
		ID:        utils.NewID(),
		UserID:    userID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}

	_, err = DB.Exec(`
		INSERT INTO sessions (id, user_id, expires_at, created_at)
		VALUES (?, ?, ?, ?)
	`, session.ID, session.UserID, session.ExpiresAt, session.CreatedAt)

	if err != nil {
		return nil, err
	}

	return session, nil
}

// GetSessionByID récupère une session par son ID.
// Retourne nil si la session n'existe pas ou est expirée.
func GetSessionByID(id string) (*models.Session, error) {
	session := &models.Session{}

	err := DB.QueryRow(`
		SELECT id, user_id, expires_at, created_at
		FROM sessions
		WHERE id = ?
	`, id).Scan(
		&session.ID,
		&session.UserID,
		&session.ExpiresAt,
		&session.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Vérifier si la session est expirée
	if time.Now().After(session.ExpiresAt) {
		DeleteSession(id)
		return nil, nil
	}

	return session, nil
}

// DeleteSession supprime une session (déconnexion).
func DeleteSession(id string) error {
	_, err := DB.Exec(`DELETE FROM sessions WHERE id = ?`, id)
	return err
}

// DeleteExpiredSessions supprime toutes les sessions expirées.
// À appeler périodiquement pour nettoyer la BDD.
func DeleteExpiredSessions() error {
	_, err := DB.Exec(`DELETE FROM sessions WHERE expires_at < ?`, time.Now())
	return err
}