package store

import (
	"path/filepath"
	"testing"
)

func TestOpenCreatesCollaborationSchema(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "hyperion.db"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer db.Close()

	for _, table := range []string{"agents", "tasks", "messages", "results"} {
		var count int
		err := db.DB.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?`, table).Scan(&count)
		if err != nil {
			t.Fatalf("query schema for %s: %v", table, err)
		}
		if count != 1 {
			t.Errorf("expected table %q to exist", table)
		}
	}
}
