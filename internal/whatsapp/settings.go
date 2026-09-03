package whatsapp

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

const (
	settingKeyAutoReply = "autoreply_enabled"
	settingKeyReplyText = "autoreply_text"
	defaultReplyText    = "Pesan Anda sudah diterima. Kami akan segera merespons."
	gmt7Offset          = 7 * time.Hour
)

type SettingsStore struct {
	db *sql.DB
}

func NewSettingsStore(dbPath string) (*SettingsStore, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("gagal buka database settings: %w", err)
	}

	if err := db.PingContext(context.Background()); err != nil {
		db.Close()
		return nil, fmt.Errorf("gagal koneksi database settings: %w", err)
	}

	store := &SettingsStore{db: db}
	if err := store.migrate(); err != nil {
		db.Close()
		return nil, err
	}

	return store, nil
}

func (s *SettingsStore) migrate() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS settings (
			key   TEXT PRIMARY KEY,
			value TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS autoreply_daily_log (
			chat_jid TEXT PRIMARY KEY,
			reply_date TEXT NOT NULL
		)`,
	}

	for _, q := range queries {
		if _, err := s.db.ExecContext(context.Background(), q); err != nil {
			return fmt.Errorf("gagal migrate settings: %w", err)
		}
	}

	// Set default enabled jika belum ada
	if !s.keyExists(settingKeyAutoReply) {
		s.setSetting(settingKeyAutoReply, "true")
	}
	if !s.keyExists(settingKeyReplyText) {
		s.setSetting(settingKeyReplyText, defaultReplyText)
	}

	return nil
}

func (s *SettingsStore) keyExists(key string) bool {
	var count int
	s.db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM settings WHERE key = ?", key).Scan(&count)
	return count > 0
}

func (s *SettingsStore) setSetting(key, value string) {
	s.db.ExecContext(context.Background(), "INSERT OR REPLACE INTO settings (key, value) VALUES (?, ?)", key, value)
}

func (s *SettingsStore) getSetting(key string) string {
	var value string
	err := s.db.QueryRowContext(context.Background(), "SELECT value FROM settings WHERE key = ?", key).Scan(&value)
	if err != nil {
		return ""
	}
	return value
}

func (s *SettingsStore) IsAutoReplyEnabled() bool {
	return strings.ToLower(s.getSetting(settingKeyAutoReply)) == "true"
}

func (s *SettingsStore) SetAutoReplyEnabled(enabled bool) {
	s.setSetting(settingKeyAutoReply, fmt.Sprintf("%v", enabled))
}

func (s *SettingsStore) GetAutoReplyText() string {
	text := s.getSetting(settingKeyReplyText)
	if text == "" {
		return defaultReplyText
	}
	return text
}

func (s *SettingsStore) SetAutoReplyText(text string) {
	s.setSetting(settingKeyReplyText, text)
}

// todayKey mengembalikan tanggal hari ini dalam format YYYY-MM-DD (GMT+7)
func todayKey() string {
	now := time.Now().UTC().Add(gmt7Offset)
	return now.Format("2006-01-02")
}

func (s *SettingsStore) HasRepliedToday(chatJID string) bool {
	var replyDate string
	err := s.db.QueryRowContext(context.Background(), "SELECT reply_date FROM autoreply_daily_log WHERE chat_jid = ?", chatJID).Scan(&replyDate)
	if err != nil {
		return false
	}
	return replyDate == todayKey()
}

func (s *SettingsStore) MarkReplied(chatJID string) {
	s.db.ExecContext(context.Background(), "INSERT OR REPLACE INTO autoreply_daily_log (chat_jid, reply_date) VALUES (?, ?)", chatJID, todayKey())
}

func (s *SettingsStore) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
