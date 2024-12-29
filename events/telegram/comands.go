package telegram

import (
	"context"
	"errors"
	"github.com/toboe512/gotbot/clients/telegram"
	"github.com/toboe512/gotbot/lib/e"
	"github.com/toboe512/gotbot/storage"
	"github.com/toboe512/gotbot/utils"
	"log"
	"strings"
)

const (
	StartCmd = "/start"
	HelpCmd  = "/help"
	RmvCmd   = "/remove"
	SaveCmd  = "/save"
	GetCmd   = "/get"
)

func (p *Processor) doCmd(ctx context.Context, text string, chatID int, username string, photo string) error {
	text = strings.TrimSpace(text)
	cmd := strings.Split(text, utils.SpaseStr)

	log.Printf("got new command '%s from '%s'", text, username)

	switch cmd[0] {
	case StartCmd:
		return p.sendHello(chatID)
	case HelpCmd:
		return p.sendHelp(chatID)
	case RmvCmd:
		if len(cmd) < 2 {
			return p.sendHelp(chatID)
		}

		return p.removeRow(ctx, chatID, username, cmd[1])
	case SaveCmd:
		if len(cmd) < 2 {
			return p.sendHelp(chatID)
		}

		if photo != utils.EmptyStr {
			return p.saveRow(ctx, chatID, photo, username, cmd[1])
		}

		return nil
	case GetCmd:
		if len(cmd) < 2 {
			return p.sendHelp(chatID)
		}

		return p.getPhoto(ctx, chatID, username, cmd[1])

	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) saveRow(ctx context.Context, chatID int, data string, username string, pwd string) (err error) {
	defer func() { err = e.WarpIfErr("can't do command: save row", err) }()

	sendMsg := NewMessageSender(chatID, p.tg)

	row := &storage.Row{
		Data:     data,
		UserName: username,
		PWD:      pwd,
	}

	isExists, err := p.storage.IsExists(ctx, row)

	if err != nil {
		return err
	}

	if isExists {
		return sendMsg(msgAlreadyExists)
	}

	if err := p.storage.Save(ctx, row); err != nil {
		return err
	}

	if err := sendMsg(msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) removeRow(ctx context.Context, chatID int, username string, pwd string) (err error) {
	defer func() { err = e.WarpIfErr("can't do command: remove", err) }()

	row := &storage.Row{
		UserName: username,
		PWD:      pwd,
	}

	sendMsg := NewMessageSender(chatID, p.tg)
	isExists, err := p.storage.IsExists(ctx, row)
	if err != nil {
		return err
	}
	if !isExists {
		return sendMsg(msgNoRemoved)
	}

	err = p.storage.Remove(ctx, row)
	if err != nil {
		return err
	}
	if err := sendMsg(msgRemoved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func (p *Processor) getPhoto(ctx context.Context, chatID int, username string, pwd string) (err error) {
	defer func() { err = e.WarpIfErr("can't do command: can't send photo", err) }()

	row, err := p.storage.GetByPwd(ctx, username, pwd)
	if err != nil && !errors.Is(err, storage.ErrNoLoadData) {
		return err
	}

	if errors.Is(err, storage.ErrNoLoadData) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	sendPhoto := NewPhotoSender(chatID, p.tg)
	if err := sendPhoto(row.Data); err != nil {
		return err
	}

	return nil
}

func NewMessageSender(chatID int, tg *telegram.Client) func(string) error {
	return func(msg string) error {
		return tg.SendMessage(chatID, msg)
	}
}

func NewPhotoSender(chatID int, tg *telegram.Client) func(string) error {
	return func(photoID string) error {
		return tg.SendPhoto(chatID, photoID)
	}
}
