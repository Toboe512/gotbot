package telegram

import (
	"context"
	"errors"
	"github.com/toboe512/gotbot/clients/telegram"
	"github.com/toboe512/gotbot/events"
	"github.com/toboe512/gotbot/lib/e"
	"github.com/toboe512/gotbot/storage"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUnknownType     = errors.New("unknown event type")
	ErrUnknownMetaType = errors.New("unknown meta type")
)

func New(client *telegram.Client, storage storage.Storage) *Processor {

	return &Processor{
		tg:      client,
		offset:  0,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Warp("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))

	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func event(udp telegram.Update) events.Event {
	udpType := fetchType(udp)

	res := events.Event{
		Type: udpType,
		Text: fetchText(udp),
	}

	if udpType == events.Message {
		res.Meta = Meta{
			ChatID:   udp.Message.Chat.ID,
			Username: udp.Message.From.UserName,
		}
	}

	return res
}

func (p Processor) Process(ctx context.Context, event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(ctx, event)
	default:
		return e.Warp("can't process message", ErrUnknownType)
	}
}

func (p *Processor) processMessage(ctx context.Context, event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Warp("can't process message", err)
	}

	if err := p.doCmd(ctx, event.Text, meta.ChatID, meta.Username); err != nil {
		return e.Warp("can't process message", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Warp("can't get meta", ErrUnknownMetaType)
	}

	return res, nil
}

func fetchText(udp telegram.Update) string {
	if udp.Message == nil {
		return ""
	}

	return udp.Message.Text
}

func fetchType(udp telegram.Update) events.Type {
	if udp.Message == nil {
		return events.Unknown
	}

	return events.Message
}
