package log

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type Entry struct {
	ID        int64     `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Module    string    `json:"module"`
	Message   string    `json:"message"`
	Details   string    `json:"details,omitempty"`
}

type Store struct {
	db *sql.DB
}

func NewStore(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			level TEXT NOT NULL,
			module TEXT NOT NULL,
			message TEXT NOT NULL,
			details TEXT
		);

		CREATE INDEX IF NOT EXISTS idx_logs_timestamp ON logs(timestamp DESC);
		CREATE INDEX IF NOT EXISTS idx_logs_level ON logs(level);
	`)
	if err != nil {
		return nil, fmt.Errorf("create logs table: %w", err)
	}

	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *Store) Add(level, module, message, details string) error {
	if details == "" {
		details = "NULL"
	} else {
		details = "'" + escapeString(details) + "'"
	}

	query := fmt.Sprintf(
		"INSERT INTO logs (timestamp, level, module, message, details) VALUES (?,?,?,?,%s)",
		details,
	)

	_, err := s.db.Exec(query, time.Now(), level, module, message)
	return err
}

func (s *Store) Get(limit int, offset int) ([]Entry, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}

	query := `
		SELECT id, timestamp, level, module, message, details
		FROM logs
		ORDER BY timestamp DESC
		LIMIT ? OFFSET ?
	`

	rows, err := s.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var logs []Entry
	for rows.Next() {
		var e Entry
		var details sql.NullString
		err := rows.Scan(&e.ID, &e.Timestamp, &e.Level, &e.Module, &e.Message, &details)
		if err != nil {
			return nil, err
		}
		if details.Valid {
			e.Details = details.String
		}
		logs = append(logs, e)
	}

	return logs, nil
}

func (s *Store) GetByTimeRange(start, end time.Time, limit, offset int) ([]Entry, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}

	query := `
		SELECT id, timestamp, level, module, message, details
		FROM logs
		WHERE timestamp >= ? AND timestamp <= ?
		ORDER BY timestamp DESC
		LIMIT ? OFFSET ?
	`

	rows, err := s.db.Query(query, start, end, limit, offset)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var logs []Entry
	for rows.Next() {
		var e Entry
		var details sql.NullString
		err := rows.Scan(&e.ID, &e.Timestamp, &e.Level, &e.Module, &e.Message, &details)
		if err != nil {
			return nil, err
		}
		if details.Valid {
			e.Details = details.String
		}
		logs = append(logs, e)
	}

	return logs, nil
}

func (s *Store) GetByLevel(level string, limit, offset int) ([]Entry, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}

	query := `
		SELECT id, timestamp, level, module, message, details
		FROM logs
		WHERE level = ?
		ORDER BY timestamp DESC
		LIMIT ? OFFSET ?
	`

	rows, err := s.db.Query(query, level, limit, offset)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var logs []Entry
	for rows.Next() {
		var e Entry
		var details sql.NullString
		err := rows.Scan(&e.ID, &e.Timestamp, &e.Level, &e.Module, &e.Message, &details)
		if err != nil {
			return nil, err
		}
		if details.Valid {
			e.Details = details.String
		}
		logs = append(logs, e)
	}

	return logs, nil
}

func (s *Store) GetCount() (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM logs").Scan(&count)
	return count, err
}

func (s *Store) Clear() error {
	_, err := s.db.Exec("DELETE FROM logs")
	return err
}

func (s *Store) DeleteBefore(before time.Time) (int64, error) {
	result, err := s.db.Exec("DELETE FROM logs WHERE timestamp < ?", before)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Store) ExportAsText(start, end time.Time) (string, error) {
	query := `
		SELECT timestamp, level, module, message, details
		FROM logs
		WHERE timestamp >= ? AND timestamp <= ?
		ORDER BY timestamp ASC
	`

	rows, err := s.db.Query(query, start, end)
	if err != nil {
		return "", err
	}
	defer func() { _ = rows.Close() }()

	var result string
	for rows.Next() {
		var timestamp time.Time
		var level, module, message string
		var details sql.NullString

		err := rows.Scan(&timestamp, &level, &module, &message, &details)
		if err != nil {
			return "", err
		}

		line := fmt.Sprintf("[%s] %s %s: %s",
			timestamp.Format("2006-01-02 15:04:05"),
			level,
			module,
			message,
		)

		if details.Valid && details.String != "" {
			line += " - " + details.String
		}

		result += line + "\n"
	}

	return result, nil
}

func (s *Store) ExportAsJSON(start, end time.Time) (string, error) {
	logs, err := s.GetByTimeRange(start, end, 0, 0)
	if err != nil {
		return "", err
	}

	if len(logs) == 0 {
		return "[]", nil
	}

	result := "[\n"
	for i, log := range logs {
		result += fmt.Sprintf(`  {
    "id": %d,
    "timestamp": "%s",
    "level": "%s",
    "module": "%s",
    "message": "%s"`,
			log.ID,
			log.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
			log.Level,
			log.Module,
			escapeJSON(log.Message),
		)

		if log.Details != "" {
			result += fmt.Sprintf(`,\n    "details": "%s"`, escapeJSON(log.Details))
		}

		result += "\n  }"

		if i < len(logs)-1 {
			result += ","
		}
		result += "\n"
	}
	result += "]"

	return result, nil
}

func escapeJSON(s string) string {
	s = replaceAll(s, "\\", "\\\\")
	s = replaceAll(s, "\"", "\\\"")
	s = replaceAll(s, "\n", "\\n")
	s = replaceAll(s, "\r", "\\r")
	s = replaceAll(s, "\t", "\\t")
	return s
}

func replaceAll(s, old, new string) string {
	result := ""
	for {
		idx := indexOf(s, old)
		if idx == -1 {
			result += s
			break
		}
		result += s[:idx] + new
		s = s[idx+len(old):]
	}
	return result
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func escapeString(s string) string {
	s = strings.ReplaceAll(s, "'", "''")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	return s
}
