package catalog

import (
	"github.com/df-mc/go-playfab/entity"
	"github.com/df-mc/go-playfab/internal"
	"github.com/df-mc/go-playfab/title"
)

type Query struct {
	// AlternateID is an alternate ID associated with the Item being retrieved
	// from [Query.Item].
	AlternateID *AlternateID `json:"AlternateId,omitempty"`
	// CustomTags is the optional custom tags associated with the request sent
	// from [Query.Item].
	CustomTags map[string]any `json:"CustomTags,omitempty"`
	// Entity is an [entity.Key] to perform any actions using the Query.
	// If left as nil and an [entity.Token] has been provided to [Query.Item],
	// it will be filled from [entity.Token.Key].
	Entity *entity.Key `json:"Entity,omitempty"`
	// ID is the unique ID of the Item being retrieved from [Query.Item].
	ID string `json:"Id,omitempty"`
}

// Item retrieves an Item from the public catalog. It does not work off a cache of the catalog
// and should be used when trying to retrieve recent updates of Item. However, please note that
// Item references data is cached and may take few moments for changes to propagate.
func (q Query) Item(t title.Title, tok *entity.Token) (zero Item, err error) {
	if q.Entity == nil && tok != nil {
		q.Entity = &tok.Entity
	}

	res, err := internal.Post[*queryResponse](t.URL().JoinPath("/Catalog/GetItem"), q, tok.SetAuthHeader)
	if err != nil {
		return zero, err
	}
	return res.Item, nil
}

type queryResponse struct {
	Item Item `json:"Item,omitempty"`
}
