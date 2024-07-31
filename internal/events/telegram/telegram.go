package telegram

import (
	"errors"

	"github.com/greenblat17/alarm-notification-bot/internal/clients/telegram"
	"github.com/greenblat17/alarm-notification-bot/internal/events"
	"github.com/greenblat17/alarm-notification-bot/internal/storage"
	e "github.com/greenblat17/alarm-notification-bot/pkg/errors"
)

type EventProcessor struct {
	tg      *telegram.Client
	storage storage.Storage
	offset  int
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrNewsEmpty        = errors.New("news are empty")
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *telegram.Client, storage storage.Storage) *EventProcessor {
	return &EventProcessor{
		tg:      client,
		storage: storage,
	}
}

func (p *EventProcessor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("cannot get events", err)
	}

	if len(updates) == 0 {
		return nil, e.Wrap("cannot get events", ErrNewsEmpty)
	}

	res := make([]events.Event, 0, len(updates))
	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *EventProcessor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return e.Wrap("cannot process message", ErrUnknownEventType)
	}
}

func (p *EventProcessor) processMessage(event events.Event) error {
	const errMsg = "cannot process message"

	meta, err := meta(event)
	if err != nil {
		return e.Wrap(errMsg, err)
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return e.Wrap(errMsg, err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("cannot get meta", ErrUnknownMetaType)
	}

	return res, nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: fetchType(upd),
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}

	return res
}

func fetchType(upd telegram.Update) events.Type {
	if upd.CheckIfMessageNil() {
		return events.Unknown
	}

	return events.Message
}

func fetchText(upd telegram.Update) string {
	if upd.CheckIfMessageNil() {
		return ""
	}

	return upd.Message.Text
}
