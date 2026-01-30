package config

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"sync"

	_ "modernc.org/sqlite"
)

const (
	KeyDefaultCarrierDir        = "defaultCarrierDir"
	KeyDefaultOutputDir         = "defaultOutputDir"
	KeyDefaultEncryptPassword   = "defaultEncryptPassword"
	KeyDefaultDecryptPassword   = "defaultDecryptPassword"
	KeyDefaultEncryptOutputName = "defaultEncryptOutputName"
	KeyAuthor                   = "author"
	KeyRepository               = "repository"
	KeyContact                  = "contact"
	defaultCarrierDirValue      = "./images"
	defaultOutputDirValue       = "./output"
	defaultEncryptPasswordVal   = ""
	defaultDecryptPasswordVal   = ""
	defaultEncryptOutputNameVal = "encrypted"
	defaultAuthorValue          = ""
	defaultRepositoryValue      = ""
	defaultContactValue         = ""
	schemaInit                  = `CREATE TABLE IF NOT EXISTS kv (k TEXT PRIMARY KEY, v TEXT NOT NULL);`
)

type Store struct {
	mu  sync.Mutex
	db  *sql.DB
	mem map[string]string
}

func NewInMemoryStore() *Store {
	return &Store{mem: map[string]string{}}
}

func NewStore(dbPath string) (*Store, error) {
	if dbPath == "" {
		return nil, errors.New("db path is empty")
	}
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(schemaInit); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) GetAllWithDefaults() map[string]string {
	m := s.GetAll()
	if m[KeyDefaultCarrierDir] == "" {
		m[KeyDefaultCarrierDir] = defaultCarrierDirValue
	}
	if m[KeyDefaultOutputDir] == "" {
		m[KeyDefaultOutputDir] = defaultOutputDirValue
	}
	if _, ok := m[KeyDefaultEncryptPassword]; !ok {
		m[KeyDefaultEncryptPassword] = defaultEncryptPasswordVal
	}
	if _, ok := m[KeyDefaultDecryptPassword]; !ok {
		m[KeyDefaultDecryptPassword] = defaultDecryptPasswordVal
	}
	if _, ok := m[KeyDefaultEncryptOutputName]; !ok {
		m[KeyDefaultEncryptOutputName] = defaultEncryptOutputNameVal
	}
	if _, ok := m[KeyAuthor]; !ok {
		m[KeyAuthor] = defaultAuthorValue
	}
	if _, ok := m[KeyRepository]; !ok {
		m[KeyRepository] = defaultRepositoryValue
	}
	if _, ok := m[KeyContact]; !ok {
		m[KeyContact] = defaultContactValue
	}
	return m
}

func (s *Store) GetAll() map[string]string {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.db == nil {
		out := make(map[string]string, len(s.mem))
		for k, v := range s.mem {
			out[k] = v
		}
		return out
	}

	rows, err := s.db.Query(`SELECT k, v FROM kv`)
	if err != nil {
		return map[string]string{}
	}
	defer func() { _ = rows.Close() }()

	out := map[string]string{}
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			continue
		}
		out[k] = v
	}
	return out
}

func (s *Store) SaveAll(values map[string]string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.db == nil {
		if s.mem == nil {
			s.mem = map[string]string{}
		}
		for k, v := range values {
			s.mem[k] = v
		}
		return nil
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	stmt, err := tx.Prepare(`INSERT INTO kv(k, v) VALUES(?, ?) ON CONFLICT(k) DO UPDATE SET v=excluded.v`)
	if err != nil {
		return err
	}
	defer func() { _ = stmt.Close() }()

	for k, v := range values {
		if k == "" {
			continue
		}
		if _, err := stmt.Exec(k, v); err != nil {
			return err
		}
	}
	return tx.Commit()
}
