package database

import (
	"testing"
	"time"
)

func TestCreateSession(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "u1", "a@a.com", "alice")

	session, err := CreateSession("u1")
	if err != nil {
		t.Fatalf("CreateSession : %v", err)
	}
	if session.ID == "" {
		t.Error("l'ID de session ne doit pas être vide")
	}
	if session.UserID != "u1" {
		t.Errorf("UserID attendu %q, obtenu %q", "u1", session.UserID)
	}
	if session.ExpiresAt.Before(time.Now()) {
		t.Error("la session doit expirer dans le futur")
	}
}

func TestCreateSession_ReplacesExisting(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "u1", "a@a.com", "alice")

	s1, _ := CreateSession("u1")
	s2, _ := CreateSession("u1")

	// La première session doit avoir été supprimée.
	old, _ := GetSessionByID(s1.ID)
	if old != nil {
		t.Error("l'ancienne session doit être supprimée après reconnexion")
	}
	current, _ := GetSessionByID(s2.ID)
	if current == nil {
		t.Error("la nouvelle session doit exister")
	}
}

func TestGetSessionByID_NotFound(t *testing.T) {
	setupTestDB(t)

	session, err := GetSessionByID("inexistant")
	if err != nil {
		t.Fatalf("GetSessionByID : %v", err)
	}
	if session != nil {
		t.Error("doit retourner nil pour un ID inexistant")
	}
}

func TestDeleteSession(t *testing.T) {
	setupTestDB(t)
	createTestUser(t, "u1", "a@a.com", "alice")

	session, _ := CreateSession("u1")
	if err := DeleteSession(session.ID); err != nil {
		t.Fatalf("DeleteSession : %v", err)
	}

	got, _ := GetSessionByID(session.ID)
	if got != nil {
		t.Error("la session supprimée ne doit plus exister")
	}
}
