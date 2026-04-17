package catalog

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/df-mc/go-playfab/entity"
	"github.com/df-mc/go-playfab/internal"
	"github.com/df-mc/go-playfab/title"
	"github.com/google/uuid"
	"golang.org/x/text/language"
)

func New(client *http.Client, title title.Title, src entity.TokenSource) *Client {
	return &Client{
		client: client,
		title:  title,
		src:    src,
	}
}

type Client struct {
	client *http.Client
	title  title.Title
	src    entity.TokenSource
}

func (c *Client) SearchItems(ctx context.Context, filter SearchFilter, opts ...internal.RequestOption) (*SearchResult, error) {
	return internal.Post[*SearchResult](
		ctx,
		c.client,
		c.title.URL().JoinPath("/Catalog/SearchItems"),
		filter,
		append(opts,
			entity.RequestOption(c.src),
		),
	)
}

func (c *Client) ItemByID(ctx context.Context, id uuid.UUID, opts ...internal.RequestOption) (*Item, error) {
	resp, err := internal.Post[*itemResponse](
		ctx,
		c.client,
		c.title.URL().JoinPath("/Catalog/GetItem"),
		itemRequest{ID: id},
		append(opts,
			entity.RequestOption(c.src),
		),
	)
	if err != nil {
		return nil, err
	}
	if resp == nil || resp.Item == nil {
		return nil, errors.New("catalog: invalid Item response")
	}
	return resp.Item, nil
}

type (
	itemRequest struct {
		AlternateIDs []AlternateID  `json:"AlternateIds,omitempty"`
		CustomTags   map[string]any `json:",omitempty"`
		Entity       entity.Key     `json:",omitzero"`
		ID           uuid.UUID      `json:"Id"`
	}
	itemResponse struct {
		Item *Item
	}
)

type (
	SearchFilter struct {
		Count             int            `json:",omitzero"`
		ContinuationToken string         `json:",omitempty"`
		CustomTags        map[string]any `json:",omitempty"`
		Entity            entity.Key     `json:",omitzero"`
		Filter            string         `json:",omitempty"`
		Language          language.Tag   `json:",omitempty"`
		OrderBy           string         `json:",omitempty"` // OData query
		Text              string         `json:"Search,omitempty"`
		Select            string         `json:",omitempty"`
	}

	SearchResult struct {
		ContinuationToken string
		Items             []Item
	}
)

func (f SearchFilter) MarshalJSON() ([]byte, error) {
	type Alias SearchFilter
	data := struct {
		Alias
		Language string `json:",omitempty"`
	}{Alias: (Alias)(f)}
	if f.Language != language.Und {
		data.Language = f.Language.String()
	}
	return json.Marshal(data)
}
