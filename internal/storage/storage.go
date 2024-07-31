package storage

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"time"

	e "github.com/greenblat17/alarm-notification-bot/pkg/errors"
)

type Storage interface {
	Save(p *Page) error
	PickRandom(username string) (*Page, error)
	Remove(p *Page) error
	IsExists(p *Page) (bool, error)
}

type Page struct {
	URL       string
	Username  string
	CreatedAt time.Time
}

var (
	ErrNoSavedPages = errors.New("no saved pages")
)

func (p *Page) Hash() (string, error) {
	const errMsg = "cannot calculate hash"

	h := sha256.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap(errMsg, err)
	}

	if _, err := io.WriteString(h, p.Username); err != nil {
		return "", e.Wrap(errMsg, err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
