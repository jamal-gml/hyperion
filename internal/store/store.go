package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type Store struct{ DB *sql.DB }

func Open(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if _, err = db.Exec(`PRAGMA journal_mode=WAL; PRAGMA foreign_keys=ON; PRAGMA busy_timeout=5000;`); err != nil {
		db.Close()
		return nil, err
	}
	s := &Store{DB: db}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}
func (s *Store) Close() error { return s.DB.Close() }
func (s *Store) migrate() error {
	_, err := s.DB.Exec(`
CREATE TABLE IF NOT EXISTS agents (
 id TEXT PRIMARY KEY, project TEXT NOT NULL, provider TEXT NOT NULL, role TEXT NOT NULL,
 status TEXT NOT NULL DEFAULT 'seen', last_seen TEXT NOT NULL, created_at TEXT NOT NULL,
 UNIQUE(project,id));
CREATE INDEX IF NOT EXISTS idx_agents_project ON agents(project);
CREATE TABLE IF NOT EXISTS tasks (
 id TEXT PRIMARY KEY, project TEXT NOT NULL, title TEXT NOT NULL, description TEXT NOT NULL,
 created_by TEXT NOT NULL, assigned_to TEXT NOT NULL DEFAULT '', status TEXT NOT NULL DEFAULT 'open',
 priority TEXT NOT NULL DEFAULT 'normal', created_at TEXT NOT NULL, updated_at TEXT NOT NULL);
CREATE INDEX IF NOT EXISTS idx_tasks_project_status ON tasks(project,status);
CREATE TABLE IF NOT EXISTS messages (
 id TEXT PRIMARY KEY, project TEXT NOT NULL, task_id TEXT NOT NULL DEFAULT '', sender TEXT NOT NULL,
 recipient TEXT NOT NULL DEFAULT '', message TEXT NOT NULL, created_at TEXT NOT NULL);
CREATE INDEX IF NOT EXISTS idx_messages_project_time ON messages(project,created_at);
CREATE TABLE IF NOT EXISTS results (
 id TEXT PRIMARY KEY, task_id TEXT NOT NULL, agent_id TEXT NOT NULL, content TEXT NOT NULL, created_at TEXT NOT NULL);
`)
	if err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}
	return nil
}
func Now() string { return time.Now().UTC().Format(time.RFC3339Nano) }
