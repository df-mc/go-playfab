package entity

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/df-mc/go-playfab/internal"
	"github.com/df-mc/go-playfab/title"
)

// A Token represents an entity token that can be used for authenticating as an entity
// in the title.
type Token struct {
	Entity     Key       `json:"Entity,omitempty"`
	Token      string    `json:"EntityToken,omitempty"`
	Expiration time.Time `json:"TokenExpiration,omitempty"`
}

// Exchange exchanges the entity token with a token responsible for another entity identified by the given [Key].
func (t *Token) Exchange(ctx context.Context, title title.Title, key Key, opts ...internal.RequestOption) (*Token, error) {
	type exchangeRequest struct {
		Entity Key `json:"Entity"`
	}
	token, err := internal.Post[*Token](ctx, internal.ContextClient(ctx), title.URL().JoinPath("/Authentication/GetEntityToken"), exchangeRequest{
		Entity: key,
	}, append(opts,
		func(req *http.Request) error {
			if req.Header.Get("X-EntityToken") != "" {
				return nil
			}
			req.Header.Set("X-EntityToken", t.Token)
			return nil
		},
	))
	if err != nil {
		return nil, err
	}
	if !token.Valid() {
		return nil, errors.New("entity: invalid token result")
	}
	return token, nil
}

// RequestOption is an [internal.RequestOption] that sets the 'X-EntityToken' header from the token
// supplied by the given [TokenSource]. If the header already exists in the request, it will be no-op.
func RequestOption(src TokenSource) internal.RequestOption {
	return func(req *http.Request) error {
		if req.Header.Get("X-EntityToken") != "" {
			return nil
		}
		token, err := src.EntityToken(req.Context())
		if err != nil {
			return fmt.Errorf("request entity token: %w", err)
		}
		req.Header.Set("X-EntityToken", token.Token)
		return nil
	}
}

// Valid reports whether the [Token] is still valid.
func (t *Token) Valid() bool {
	return t != nil && t.Token != "" && !t.Expired()
}

// Expired determines whether the entity token has expired.
// Users should exchange this token as soon as possible at
// certain interval for avoid expiring.
func (t *Token) Expired() bool {
	return time.Now().After(t.Expiration.Add(-time.Minute * 15))
}

// SetAuthHeader sets an 'X-EntityToken' header on the request.
// This is the primary method for authenticating with PlayFab API.
func (t *Token) SetAuthHeader(req *http.Request) {
	req.Header.Set("X-EntityToken", t.Token)
}

// Key identifies the entity within the PlayFab platform.
type Key struct {
	// ID is the unique ID of the entity.
	ID string `json:"Id,omitempty"`
	// Type is the type of entity. It is one of constants defined below.
	Type string `json:"Type,omitempty"`
}

const (
	// TypeNamespace indicates an entity has access to all global information
	// for all titles within a studio in PlayFab. Consumers of the title normally
	// are not allowed to sign in with this type. When this type is used, the
	// [Key.ID] refers to the Publisher ID of a PlayFab studio.
	//
	// See: https://learn.microsoft.com/en-us/gaming/playfab/live-service-management/game-configuration/entities/available-built-in-entity-types#namespace
	TypeNamespace = "namespace"

	// TypeTitle indicates an entity has access to all global information for
	// a title. When this type is used for the entity, the [Key.ID] refers to
	// the ID for the PlayFab title within a studio. Consumers of the title
	// normally are not allowed sign in with this type.
	//
	// See:https://learn.microsoft.com/en-us/gaming/playfab/live-service-management/game-configuration/entities/available-built-in-entity-types#title
	TypeTitle = "title"

	// TypeMasterPlayerAccount indicates an entity is a player entity that is
	// shared by all titles within a studio. When this type is used for the entity
	// key, the [Key.ID] refers to the [playfab.Identity.PlayFabID].
	//
	// See: https://learn.microsoft.com/en-us/gaming/playfab/live-service-management/game-configuration/entities/available-built-in-entity-types#master_player_account
	TypeMasterPlayerAccount = "master_player_account"

	// TypeTitlePlayerAccount indicates an entity that is representing a player
	// within a title in the most traditional way. When this type is used for the
	// entity key, the [Key.ID] refers to the unique ID of the player within a title.
	//
	// See: https://learn.microsoft.com/en-us/gaming/playfab/live-service-management/game-configuration/entities/available-built-in-entity-types#title_player_account
	TypeTitlePlayerAccount = "title_player_account"

	// TypeCharacter indicates that an entity is a sub-entity of TypeTitlePlayerAccount.
	// When this type is used for the entity key, the [Key.ID] refers to the ID of the character
	// owned by the user.
	//
	// See: https://learn.microsoft.com/en-us/gaming/playfab/live-service-management/game-configuration/entities/available-built-in-entity-types#character
	TypeCharacter = "character"

	// TypeGroup indicates that an entity is a container for other entities.
	// When this type is used for the entity key, the [Key.ID] refers to the ID of the group.
	//
	// See: https://learn.microsoft.com/en-us/gaming/playfab/live-service-management/game-configuration/entities/available-built-in-entity-types#group
	TypeGroup = "group"

	// TypeGameServer indicates that an entity is used by game servers primarily for use
	// in the Matchmaking and Lobby features of the PlayFab. Future scenarios may be added
	// to support other PlayFab features.
	//
	// See: https://learn.microsoft.com/en-us/gaming/playfab/live-service-management/game-configuration/entities/available-built-in-entity-types#game_server
	TypeGameServer = "game_server"
)
