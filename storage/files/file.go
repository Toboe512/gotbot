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

func (s Storage) Save(ctx context.Context, page *storage.Page) (err error) {
	defer func() { err = e.WarpIfErr("can't save page in file: ", err) }()

	fPath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return err
	}

	fName, err := fileName(page)
	if err != nil {
		return err
	}

	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil

}

func (s Storage) PickRandom(ctx context.Context, userName string) (page *storage.Page, err error) {
	defer func() { err = e.WarpIfErr("can't pick random page in file: ", err) }()

	fPath := filepath.Join(s.basePath, userName)
	files, err := os.ReadDir(fPath)

	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	n := rand.Intn(len(files))
	file := files[n]

	return s.decodePage(filepath.Join(fPath, file.Name()))

}

func (s Storage) Remove(ctx context.Context, p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return e.Warp("can't remove file", err)
	}

	fPath := filepath.Join(s.basePath, p.UserName, fileName)

	if err := os.Remove(fPath); err != nil {
		msg := fmt.Sprintf("can't remove file %s", fPath)
		return e.Warp(msg, err)
	}

	return nil
}

func (s Storage) IsExists(ctx context.Context, p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, e.Warp("can't check if exists page in file", err)
	}

	fPath := filepath.Join(s.basePath, p.UserName, fileName)

	switch _, err = os.Stat(fPath); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if file %s exists", fPath)
		return false, e.Warp(msg, err)
	}

	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)

	if err != nil {
		return nil, e.Warp("can't decode page", err)
	}

	defer func() { _ = f.Close() }()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Warp("can't decode page", err)
	}

	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
