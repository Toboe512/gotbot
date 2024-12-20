package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/toboe512/gotbot/storage"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)

	if err != nil {
		return nil, fmt.Errorf("can't open database sqlite: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't open database sqlite: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Save(ctx context.Context, row *storage.Row) error {
	q := `INSERT INTO usr_data (user_name, pwd, data) VALUES (?, ?, ?)`

	if _, err := s.db.ExecContext(ctx, q, row.UserName, row.PWD, row.Data); err != nil {
		return fmt.Errorf("can't save row in sqlite: %w", err)
	}

	return nil
}

func (s *Storage) PickRandom(ctx context.Context, userName string) (*storage.Row, error) {
	q := `SELECT data FROM usr_data WHERE  user_name = ? ORDER BY RANDOM() LIMIT 1`

	var imgID string

	err := s.db.QueryRowContext(ctx, q, userName).Scan(&imgID)

	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSavedPages
	}

	if err != nil {
		return nil, fmt.Errorf("can't pick random row in sqlite: %w", err)
	}

	return &storage.Row{
		Data:     imgID,
		UserName: userName,
	}, nil
}

func (s *Storage) GetByPwd(ctx context.Context, userName string, pwd string) (*storage.Row, error) {
	q := `SELECT data FROM usr_data WHERE user_name = ? AND pwd = ?`

	var data string

	err := s.db.QueryRowContext(ctx, q, userName, pwd).Scan(&data)

	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSavedPages
	}

	if err != nil {
		return nil, fmt.Errorf("can't get by row pwd in sqlite: %w", err)
	}

	return &storage.Row{
		Data:     data,
		UserName: userName,
		PWD:      pwd,
	}, nil
}

func (s *Storage) Remove(ctx context.Context, row *storage.Row) error {
	q := `DELETE FROM usr_data WHERE user_name = ? AND pwd = ?`

	if _, err := s.db.ExecContext(ctx, q, row.UserName, row.PWD); err != nil {
		return fmt.Errorf("can't remove row in sqlite: %w", err)
	}

	return nil
}

func (s *Storage) IsExists(ctx context.Context, p *storage.Row) (bool, error) {
	q := `SELECT COUNT(*) FROM usr_data WHERE user_name = ? AND pwd = ?`

	var count int

	if err := s.db.QueryRowContext(ctx, q, p.UserName, p.PWD).Scan(&count); err != nil {
		return false, fmt.Errorf("can't check if exists row in sqlite: %w", err)
	}

	return count > 0, nil
}

func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS usr_data (user_name TEXT, pwd TEXT, data TEXT)`

	_, err := s.db.ExecContext(ctx, q)

	if err != nil {
		return fmt.Errorf("can't create table in sqlite: %w", err)
	}

	return nil
}
