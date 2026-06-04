package database

import (
	"forum/utils"
	"time"
)

// ToggleLike gère le like/dislike sur un post ou commentaire.
// - Si l'utilisateur n'a pas encore voté → INSERT le vote
// - Si l'utilisateur a voté pareil → DELETE (annule le vote)
// - Si l'utilisateur a voté différemment → UPDATE (change le vote)
func ToggleLike(userID, targetID, targetType string, value int) error {
	var existingValue int
	var existingID string

	err := DB.QueryRow(`
		SELECT id, value FROM likes
		WHERE user_id = ? AND target_id = ? AND target_type = ?
	`, userID, targetID, targetType).Scan(&existingID, &existingValue)

	if err != nil {
		// Pas de vote existant → INSERT
		_, err = DB.Exec(`
			INSERT INTO likes (id, user_id, target_id, target_type, value, created_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`, utils.NewID(), userID, targetID, targetType, value, time.Now())
		return err
	}

	if existingValue == value {
		// Même vote → annuler (DELETE)
		_, err = DB.Exec(`DELETE FROM likes WHERE id = ?`, existingID)
		return err
	}

	// Vote différent → changer (UPDATE)
	_, err = DB.Exec(`
		UPDATE likes SET value = ? WHERE id = ?
	`, value, existingID)
	return err
}