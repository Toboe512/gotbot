package telegram

import (
	"context"
	"errors"
	"github.com/toboe512/gotbot/clients/telegram"
	"github.com/toboe512/gotbot/events"
	"github.com/toboe512/gotbot/lib/e"
	"github.com/toboe512/gotbot/storage"
	"github.com/toboe512/gotbot/utils"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

// Meta Структура с данными которые получаем от пользователя.
type Meta struct {
	ChatID   int
	Username string
	ImageID  string
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

func (p Processor) Process(ctx context.Context, event events.Event) error {
	switch event.Type {
	case events.Message, events.Image:
		return p.processMessage(ctx, event)
	default:
		return e.Warp("can't process message", ErrUnknownType)
	}
}

// event метод в котором по сути происходит мепинг Update в Event с заполнением структуры Meta.
func event(udp telegram.Update) events.Event {
	udpType := fetchType(udp)

	res := events.Event{
		Type: udpType,
		Text: fetchText(udp),
	}

	//TODO реализовать возможность сохранения нескольких фото
	img := utils.EmptyStr
	if udp.Message.Photo != nil {
		img = udp.Message.Photo[0].FileID
	}

	if udpType != events.Unknown {
		res.Meta = Meta{
			ChatID:   udp.Message.Chat.ID,
			Username: udp.Message.From.UserName,
			ImageID:  img,
		}
	}

	return res
}

func (p *Processor) processMessage(ctx context.Context, event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Warp("can't process message", err)
	}

	if err := p.DoCmd(ctx, event.Text, meta.ChatID, meta.Username, meta.ImageID); err != nil {
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
	switch fetchType(udp) {
	case events.Message:
		return udp.Message.Text
	case events.Image:
		return udp.Message.Caption
	default:
		return utils.EmptyStr
	}
}

func fetchType(udp telegram.Update) events.Type {
	if udp.Message == nil {
		return events.Unknown
	}

	if udp.Message.Photo == nil || len(udp.Message.Photo) == 0 {
		return events.Message
	}

	return events.Image
}
