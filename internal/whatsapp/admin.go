package whatsapp

import (
	"context"
	"database/sql"
)

type Admin struct {
	db *sql.DB
}

func NewAdmin(db *sql.DB) *Admin {
	return &Admin{
		db: db,
	}
}

func (s *Admin) Migrate(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS admins (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		jid TEXT NOT NULL UNIQUE,
		role TEXT NOT NULL DEFAULT 'admin',
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)

	return err
}

func (s *Admin) AddAdmin(ctx context.Context, jid string, role string) error {
	_, err := s.db.ExecContext(ctx, `
	INSERT INTO admins (jid, role)
	VALUES (?, ?)
	`, jid, role)

	return err
}

func (s *Admin) RemoveAdmin(ctx context.Context, jid string) error {
	_, err := s.db.ExecContext(ctx, `
	DELETE FROM admins
	WHERE jid = ?
	`, jid)

	return err
}

func (s *Admin) IsAdmin(ctx context.Context, jid string) (bool, error) {
	var exist bool

	err := s.db.QueryRowContext(ctx, `
	SELECT EXISTS(
    SELECT 1
    FROM admins
    WHERE jid = ?
)
	`, jid).Scan(&exist)

	return exist, err
}
