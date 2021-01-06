package mastodon

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	mast "github.com/mattn/go-mastodon"
)

type StatusURL struct {
	Server string
	ID     string
}

func ParseStatusURL(ref string) (*StatusURL, error) {
	parsedURL, err := url.Parse(ref)
	if err != nil {
		return nil, err
	}
	u := StatusURL{}
	u.Server = fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
	segments := strings.Split(parsedURL.Path, "/")
	if len(segments) < 3 {
		return nil, fmt.Errorf("not enough segments in the status URL")
	}
	u.ID = segments[2]
	return &u, nil
}

type StatusAccount struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Acct     string `json:"acct"`
	Avatar   string `json:"avatar"`
}

type Status struct {
	ID        string        `json:"id"`
	Content   string        `json:"content"`
	CreatedAt time.Time     `json:"created_at"`
	URL       string        `json:"url"`
	Account   StatusAccount `json:"account"`
}

type Client struct{}

func New() *Client {
	return &Client{}
}

func (c *Client) GetStatus(ctx context.Context, u string) (*Status, error) {
	su, err := ParseStatusURL(u)
	if err != nil {
		return nil, err
	}
	mc := mast.NewClient(&mast.Config{
		Server: su.Server,
	})
	s, err := mc.GetStatus(ctx, mast.ID(su.ID))
	if err != nil {
		return nil, err
	}
	cs := convertStatus(su.Server, s)
	return &cs, nil
}

func (c *Client) GetDescendants(ctx context.Context, u string) ([]Status, error) {
	su, err := ParseStatusURL(u)
	if err != nil {
		return nil, err
	}
	mc := mast.NewClient(&mast.Config{
		Server: su.Server,
	})
	sc, err := mc.GetStatusContext(ctx, mast.ID(su.ID))
	if err != nil {
		return nil, err
	}
	return convertStatuses(su.Server, sc.Descendants), nil
}

func convertStatuses(serverURL string, input []*mast.Status) []Status {
	converted := make([]Status, 0, len(input))
	for _, i := range input {
		converted = append(converted, convertStatus(serverURL, i))
	}
	return converted
}

func convertStatus(serverURL string, i *mast.Status) Status {
	return Status{
		ID:        string(i.ID),
		Content:   i.Content,
		CreatedAt: i.CreatedAt,
		URL:       i.URL,
		Account: StatusAccount{
			ID:       string(i.Account.ID),
			Acct:     NormalizeAcct(serverURL, i.Account.Acct),
			Username: i.Account.Username,
			Avatar:   i.Account.AvatarStatic,
		},
	}
}
