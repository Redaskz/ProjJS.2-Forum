package database

import (
	"testing"
)

func setupLikeTest(t *testing.T) {
	t.Helper()
	setupTestDB(t)
	createTestUser(t, "u1", "a@a.com", "alice")
	_, err := DB.Exec(
		`INSERT INTO posts (id, user_id, title, content, created_at, updated_at) VALUES ('p1', 'u1', 'Test', 'Contenu', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
	)
	if err != nil {
		t.Fatalf("impossible de créer le post de test : %v", err)
	}
}

func countLikes(t *testing.T, targetID string) int {
	t.Helper()
	var n int
	DB.QueryRow(`SELECT COUNT(*) FROM likes WHERE target_id = ?`, targetID).Scan(&n)
	return n
}

func getLikeValue(t *testing.T, userID, targetID string) int {
	t.Helper()
	var v int
	DB.QueryRow(`SELECT value FROM likes WHERE user_id = ? AND target_id = ?`, userID, targetID).Scan(&v)
	return v
}

func TestToggleLike_New(t *testing.T) {
	setupLikeTest(t)

	if err := ToggleLike("u1", "p1", "post", 1); err != nil {
		t.Fatalf("ToggleLike : %v", err)
	}
	if countLikes(t, "p1") != 1 {
		t.Error("un like doit être enregistré")
	}
	if getLikeValue(t, "u1", "p1") != 1 {
		t.Error("la valeur du like doit être 1")
	}
}

func TestToggleLike_SameVoteCancels(t *testing.T) {
	setupLikeTest(t)

	ToggleLike("u1", "p1", "post", 1)
	ToggleLike("u1", "p1", "post", 1) // même vote → annulation

	if countLikes(t, "p1") != 0 {
		t.Error("le vote doit être annulé après un double clic")
	}
}

func TestToggleLike_DifferentVoteUpdates(t *testing.T) {
	setupLikeTest(t)

	ToggleLike("u1", "p1", "post", 1)
	ToggleLike("u1", "p1", "post", -1) // change en dislike

	if countLikes(t, "p1") != 1 {
		t.Error("il doit rester exactement un vote")
	}
	if getLikeValue(t, "u1", "p1") != -1 {
		t.Error("la valeur doit être mise à jour à -1")
	}
}

func TestToggleLike_Dislike(t *testing.T) {
	setupLikeTest(t)

	if err := ToggleLike("u1", "p1", "post", -1); err != nil {
		t.Fatalf("ToggleLike dislike : %v", err)
	}
	if getLikeValue(t, "u1", "p1") != -1 {
		t.Error("la valeur du dislike doit être -1")
	}
}
