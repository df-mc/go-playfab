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

// New returns a new Client from the provided components.
func New(client *http.Client, title title.Title, src entity.TokenSource) *Client {
	return &Client{
		client: client,
		title:  title,
		src:    src,
	}
}

// Client implements a client communicating with the PlayFab Catalog API.
type Client struct {
	client *http.Client
	title  title.Title
	src    entity.TokenSource
}

// SearchItems searches for items in the catalog.
func (c *Client) SearchItems(ctx context.Context, filter SearchFilter, opts ...internal.RequestOption) (*SearchResult, error) {
	return internal.Post[*SearchResult](
		ctx,
		c.client,
		c.title.URL().JoinPath("/Catalog/SearchItems"),
		filter,
		append(opts,
			entity.RequestOption(c.src),
			internal.AcceptLanguage(internal.DefaultLanguage),
		),
	)
}

// ItemByID retrieves an Item by the ID.
func (c *Client) ItemByID(ctx context.Context, id uuid.UUID, opts ...internal.RequestOption) (*Item, error) {
	resp, err := internal.Post[*itemResponse](
		ctx,
		c.client,
		c.title.URL().JoinPath("/Catalog/GetItem"),
		itemRequest{ID: id},
		append(opts,
			entity.RequestOption(c.src),
			internal.AcceptLanguage(internal.DefaultLanguage),
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
	// itemRequest represents a request payload used for retrieving an item by ID.
	itemRequest struct {
		// AlternateIDs is a list of AlternateID that may be associated
		// with the resulting Item.
		AlternateIDs []AlternateID `json:"AlternateIds,omitempty"`
		// CustomTags are the custom tags associated with the request.
		CustomTags map[string]any `json:",omitempty"`
		// Entity specifies whose perspective is used for querying an Item.
		Entity entity.Key `json:",omitzero"`
		// ID is the UUID associated with the Item.
		ID uuid.UUID `json:"Id"`
	}
	// itemResponse represents a successful response for [Client.ItemByID].
	itemResponse struct {
		// Item is the resulting Item.
		Item *Item
	}
)

type (
	// SearchFilter is the search filter applied for the search.
	SearchFilter struct {
		// Count is the number of returned items included in the SearchResult.
		// The maximum value is 50, and defaulted to 10 by the service-side.
		Count int `json:",omitzero"`
		// ContinuationToken is the opaque token used for continuing the search,
		// if any are available. It is normally filled from [SearchResult.ContinuationToken].
		ContinuationToken string `json:",omitempty"`
		// CustomTags is the optional properties associated with the request.
		CustomTags map[string]any `json:",omitempty"`
		// Entity is an [entity.Key] to perform any actions using the Filter.
		// If left as nil and an [entity.Token] has been provided to [Filter.Search],
		// it will be filled from [entity.Token.Key].
		Entity entity.Key `json:",omitzero"`
		// Filter is an OData query for filtering the items included in SearchResult.
		// For example, "<Field of Item> eq '<Value of Item>'".
		Filter string `json:",omitempty"`
		// Language is the locale to be included in the dictionary of Items returned in
		// the SearchResult. It is also used as an 'Accept-Language' header of the request
		// sent from [SearchResult.Search].
		Language language.Tag `json:",omitempty"`
		// OrderBy is an OData sort query for sorting the index of SearchResult. Defaulted to relevance.
		OrderBy string `json:",omitempty"` // OData query
		// Term is the string terms to be searched.
		Term string `json:"Search,omitempty"`
		// Select is an OData selection query for filtering the fields of returned items included in the SearchResult.
		Select string `json:",omitempty"`
	}

	// SearchResult describes a successful response for [Client.SearchItems].
	SearchResult struct {
		// ContinuationToken provides an opaque token for continuing the next page of
		// SearchResult by specifying it to [Filter.ContinuationToken], if any are available.
		ContinuationToken string
		// Items is a paginated set of Item for the search query.
		Items []Item
	}
)

// MarshalJSON ...
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
