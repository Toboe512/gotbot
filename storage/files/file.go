package files

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/toboe512/gotbot/lib/e"
	"github.com/toboe512/gotbot/storage"
)

type Storage struct {
	basePath string
}

const defaultPerm = 0774

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(ctx context.Context, row *storage.Row) (err error) {
	defer func() { err = e.WarpIfErr("can't save row in file: ", err) }()

	fPath := filepath.Join(s.basePath, row.UserName)

	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return err
	}

	fName, err := fileName(row)
	if err != nil {
		return err
	}

	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(row); err != nil {
		return err
	}

	return nil

}

func (s Storage) PickRandom(ctx context.Context, userName string) (row *storage.Row, err error) {
	defer func() { err = e.WarpIfErr("can't pick random row in file: ", err) }()

	fPath := filepath.Join(s.basePath, userName)
	files, err := os.ReadDir(fPath)

	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedData
	}

	n := rand.Intn(len(files))
	file := files[n]

	return s.decodePage(filepath.Join(fPath, file.Name()))

}

func (s *Storage) GetByPwd(ctx context.Context, userName string, pwd string) (row *storage.Row, err error) {
	//TODO реализовать метод
	return nil, nil
}

func (s Storage) Remove(ctx context.Context, row *storage.Row) error {
	fileName, err := fileName(row)
	if err != nil {
		return e.Warp("can't remove file", err)
	}

	fPath := filepath.Join(s.basePath, row.UserName, fileName)

	if err := os.Remove(fPath); err != nil {
		msg := fmt.Sprintf("can't remove file %s", fPath)
		return e.Warp(msg, err)
	}

	return nil
}

func (s Storage) IsExists(ctx context.Context, row *storage.Row) (bool, error) {
	fileName, err := fileName(row)
	if err != nil {
		return false, e.Warp("can't check if exists page in file", err)
	}

	fPath := filepath.Join(s.basePath, row.UserName, fileName)

	switch _, err = os.Stat(fPath); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if file %s exists", fPath)
		return false, e.Warp(msg, err)
	}

	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Row, error) {
	f, err := os.Open(filePath)

	if err != nil {
		return nil, e.Warp("can't decode page", err)
	}

	defer func() { _ = f.Close() }()

	var row storage.Row

	if err := gob.NewDecoder(f).Decode(&row); err != nil {
		return nil, e.Warp("can't decode page", err)
	}

	return &row, nil
}

func fileName(row *storage.Row) (string, error) {
	return row.Hash()
}
