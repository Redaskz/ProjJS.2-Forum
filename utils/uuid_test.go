package utils

import (
	"strings"
	"testing"
)

func TestNewID_Format(t *testing.T) {
	id := NewID()
	parts := strings.Split(id, "-")
	if len(parts) != 5 {
		t.Errorf("UUID attendu en 5 parties, obtenu %d dans %q", len(parts), id)
	}
}

func TestNewID_Unique(t *testing.T) {
	seen := make(map[string]bool, 1000)
	for i := 0; i < 1000; i++ {
		id := NewID()
		if seen[id] {
			t.Fatalf("ID dupliqué généré : %s", id)
		}
		seen[id] = true
	}
}

func TestNewID_NotEmpty(t *testing.T) {
	if id := NewID(); id == "" {
		t.Error("NewID ne doit pas retourner une chaîne vide")
	}
}
