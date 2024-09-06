package entity

import (
	"github.com/df-mc/go-playfab/internal"
	"github.com/df-mc/go-playfab/title"
)

// Exchange exchanges a Token of TypeMasterPlayerAccount with the ID.
func (tok *Token) Exchange(t title.Title, id string) (_ *Token, err error) {
	r := exchange{
		Entity: Key{
			Type: TypeMasterPlayerAccount,
			ID:   id,
		},
	}

	return internal.Post[*Token](t.URL().JoinPath("/Authentication/GetEntityToken"), r, tok.SetAuthHeader)
}

type exchange struct {
	CustomTags map[string]any `json:"CustomTags,omitempty"`
	Entity     Key            `json:"Entity,omitempty"`
}
