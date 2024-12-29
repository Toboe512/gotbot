package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/toboe512/gotbot/storage"
	"github.com/toboe512/gotbot/utils"
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
	q := `INSERT INTO usr_data (user_name, pwd, salt, data) VALUES (?, ?, ?, ?)`

	data, salt, err := utils.EncryptStrAes(row.PWD, row.Data)
	if err != nil {
		return err
	}

	pwd := utils.StringToSha256(row.PWD)
	if _, err := s.db.ExecContext(ctx, q, row.UserName, pwd, salt, data); err != nil {
		return fmt.Errorf("can't save row in sqlite: %w", err)
	}

	return nil
}

func (s *Storage) GetByPwd(ctx context.Context, userName string, pwd string) (*storage.Row, error) {
	q := `SELECT data, salt FROM usr_data WHERE user_name = ? AND pwd = ?`
	pwdHashed := utils.StringToSha256(pwd)

	var data string
	var salt string

	err := s.db.QueryRowContext(ctx, q, userName, pwdHashed).Scan(&data, &salt)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoLoadData
	}

	if err != nil {
		return nil, fmt.Errorf("can't get by row pwd in sqlite: %w", err)
	}

	data, err = utils.DecryptStrAes(pwd, salt, data)
	if err != nil {
		return nil, err
	}

	return &storage.Row{
		Data:     data,
		UserName: userName,
		PWD:      pwd,
	}, nil
}

func (s *Storage) Remove(ctx context.Context, row *storage.Row) error {
	q := `DELETE FROM usr_data WHERE user_name = ? AND pwd = ?`

	pwd := utils.StringToSha256(row.PWD)
	if _, err := s.db.ExecContext(ctx, q, row.UserName, pwd); err != nil {
		return fmt.Errorf("can't remove row in sqlite: %w", err)
	}

	return nil
}

func (s *Storage) IsExists(ctx context.Context, row *storage.Row) (bool, error) {
	q := `SELECT COUNT(*) FROM usr_data WHERE user_name = ? AND pwd = ?`

	var count int
	pwdHashed := utils.StringToSha256(row.PWD)

	if err := s.db.QueryRowContext(ctx, q, row.UserName, pwdHashed).Scan(&count); err != nil {
		return false, fmt.Errorf("can't check if exists row in sqlite: %w", err)
	}

	return count > 0, nil
}

func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS usr_data (user_name TEXT, pwd TEXT, salt TEXT, data TEXT)`

	_, err := s.db.ExecContext(ctx, q)

	if err != nil {
		return fmt.Errorf("can't create table in sqlite: %w", err)
	}

	return nil
}
