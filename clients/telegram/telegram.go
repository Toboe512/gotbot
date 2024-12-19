package telegram

import (
	"encoding/json"
	"github.com/toboe512/gotbot/lib/e"
	"io"
	"log"
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

	data, err := c.doRequest(getUpdatesMethod, q)

	if err != nil {
		return nil, err
	}

	var resp UpdatesResponse

	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	if len(resp.Result) != 0 {
		log.Print(resp.Result[0].Message.Chat.IsForum)
	}

	return resp.Result, nil
}

func (c Client) SendMessage(chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := c.doRequest(sensMessageMethod, q)

	if err != nil {
		return e.Warp("can't send message", err)
	}

	return nil
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	defer func() { err = e.WarpIfErr("can't do request", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()
	resp, err := c.client.Do(req)

	if err != nil {
		return nil, e.WarpIfErr("can't do request", err)
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {

		return nil, e.WarpIfErr("can't do request", err)
	}

	return body, nil
}
