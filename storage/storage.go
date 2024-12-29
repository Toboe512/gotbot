package storage

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/toboe512/gotbot/lib/e"
	"github.com/toboe512/gotbot/utils"
	"io"
)

type Storage interface {
	Save(ctx context.Context, r *Row) error
	GetByPwd(ctx context.Context, userName string, pwd string) (*Row, error)
	Remove(ctx context.Context, r *Row) error
	IsExists(ctx context.Context, r *Row) (bool, error)
}

type Row struct {
	UserName string
	PWD      string
	Data     string
}

var ErrNoSavedData = errors.New("no saved pages")
var ErrNoLoadData = errors.New("no load data")

func (p Row) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.Data); err != nil {
		return utils.EmptyStr, e.Warp("can't calculate hash", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return utils.EmptyStr, e.Warp("can't calculate hash", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
