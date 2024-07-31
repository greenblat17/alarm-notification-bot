package telegram

import (
	"errors"
	"log"
	"net/url"
	"strings"

	"github.com/greenblat17/alarm-notification-bot/internal/storage"
	e "github.com/greenblat17/alarm-notification-bot/pkg/errors"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *EventProcessor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command [%s] from [%s]", text, username)

	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.sendUnknown(chatID)
	}
}

func (p *EventProcessor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.Wrap("cannot do command: save page", err) }()

	page := &storage.Page{
		URL:      pageURL,
		Username: username,
	}

	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}
	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *EventProcessor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.Wrap("cannot do command: send random", err) }()

	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(page)
}

func (p *EventProcessor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *EventProcessor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func (p *EventProcessor) sendUnknown(chatID int) error {
	return p.tg.SendMessage(chatID, msgUnknownCommand)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
