package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/greenblat17/alarm-notification-bot/internal/storage"
	e "github.com/greenblat17/alarm-notification-bot/pkg/errors"
)

type Storage struct {
	basePath string
}

const (
	defaultPerm = 0774
)

var (
	ErrNoSavedPages = errors.New("no saved pages")
)

func New(basePath string) *Storage {
	return &Storage{basePath: basePath}
}

func (s *Storage) Save(page *storage.Page) (err error) {
	defer func() { err = e.WrapWithErr("cannot save page", err) }()

	fPath := filepath.Join(s.basePath, page.Username)

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

func (s *Storage) PickRandom(username string) (page *storage.Page, err error) {
	defer func() { err = e.WrapWithErr("cannot pick random page", err) }()

	fPath := filepath.Join(s.basePath, username)

	files, err := os.ReadDir(fPath)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, ErrNoSavedPages
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))
	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(fPath, file.Name()))
}

func (s *Storage) Remove(p *storage.Page) error {
	const errMsg = "cannot remove page"

	fileName, err := fileName(p)
	if err != nil {
		return e.Wrap(errMsg, err)
	}

	fPath := filepath.Join(s.basePath, p.Username, fileName)

	if err := os.Remove(fPath); err != nil {
		msg := fmt.Sprintf("%s %s", errMsg, fPath)

		return e.Wrap(msg, err)
	}

	return nil
}

func (s *Storage) IsExists(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("cannot check if file exists", err)
	}

	fPath := filepath.Join(s.basePath, p.Username, fileName)

	switch _, err = os.Stat(fPath); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("cannot check if file %s exists", fPath)

		return false, e.Wrap(msg, err)
	}

	return true, nil
}

func (s *Storage) decodePage(filePath string) (*storage.Page, error) {
	const errMsg = "cannot decode page"

	f, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap(errMsg, err)
	}
	defer func() { err = f.Close() }()

	var p storage.Page
	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Wrap(errMsg, err)
	}

	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
