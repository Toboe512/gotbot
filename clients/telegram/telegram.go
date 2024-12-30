package telegram

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/toboe512/gotbot/lib/e"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

const (
	getUpdatesMethod  = "getUpdates"
	sensMessageMethod = "sendMessage"
	sendPhoto         = "sendPhoto"
	setMyCommands     = "setMyCommands"
)

func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func (c *Client) Updates(offset int, limit int) (updates []Update, err error) {
	defer func() { err = e.WarpIfErr("can't do updates ", err) }()

	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(http.MethodGet, getUpdatesMethod, q, nil)
	if err != nil {
		return nil, err
	}

	var resp UpdatesResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return resp.Result, nil
}

func (c Client) SendMessage(chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := c.doRequest(http.MethodGet, sensMessageMethod, q, nil)
	if err != nil {
		return e.Warp("can't send message", err)
	}

	return nil
}

func (c Client) SendPhoto(chatID int, photoID string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("photo", photoID)

	_, err := c.doRequest(http.MethodGet, sendPhoto, q, nil)
	if err != nil {
		return e.Warp("can't send photo", err)
	}

	return nil
}

func (c Client) SetCmd(commands []BotCommand) (err error) {
	defer func() { err = e.WarpIfErr("can't do set commands", err) }()

	body, err := json.Marshal(Commands{Commands: commands})
	if err != nil {
		return err
	}

	data, err := c.doRequest(http.MethodPost, setMyCommands, nil, bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	var resp BaseResult
	if err := json.Unmarshal(data, &resp); err != nil {
		return err
	}

	if !resp.OK && !resp.Result {
		return errors.New(fmt.Sprint("Response failed"))
	}

	return nil
}

func (c Client) DeleteCmd() (err error) {
	defer func() { err = e.WarpIfErr("can't do delete commands", err) }()

	data, err := c.doRequest(http.MethodPost, setMyCommands, nil, nil)

	if err != nil {
		return err
	}

	var resp BaseResult
	if err := json.Unmarshal(data, &resp); err != nil {
		return err
	}

	if !resp.OK && !resp.Result {
		return errors.New(fmt.Sprint("Response failed"))
	}

	return nil
}

func (c *Client) doRequest(nttpMethod string, method string, query url.Values, body io.Reader) (data []byte, err error) {
	defer func() { err = e.WarpIfErr("can't do request", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(nttpMethod, u.String(), body)
	if err != nil {
		return nil, err
	}

	if query != nil {
		req.URL.RawQuery = query.Encode()
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, e.WarpIfErr("can't do request", err)
	}

	defer func() { _ = resp.Body.Close() }()

	rsBody, err := io.ReadAll(resp.Body)
	if err != nil {

		return nil, e.WarpIfErr("can't do request", err)
	}

	return rsBody, nil
}

func newBasePath(token string) string {
	return "bot" + token
}
