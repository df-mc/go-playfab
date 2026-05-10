package playfab

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/df-mc/go-playfab/v2/entity"
	"github.com/df-mc/go-playfab/v2/internal"
	"github.com/df-mc/go-playfab/v2/title"
)

// LoginRequest represents a base structure used for signing in to the PlayFab through a method.
// An implementation of IdentityProvider may embed LoginRequest with additional fields required
// to fill in.
type LoginRequest struct {
	// Title is the title ID to perform this action on.
	Title title.Title `json:"TitleId"`
	// CreateAccount specifies to automatically create a PlayFab account if
	// one is not currently linked to the ID.
	CreateAccount bool `json:",omitempty"`
	// CustomTags is the optional custom tags associated with the request.
	CustomTags map[string]any `json:",omitempty"`
	// EncryptedRequest is a base64-encoded body that is encrypted with the
	// Title's public RSA key (Enterprise Only).
	EncryptedRequest []byte `json:",omitempty"`
	// InfoParameters specifies parameters to request additional data
	// on the LoginResult.
	InfoParameters *LoginInfoRequest `json:"InfoRequestParameters,omitempty"`
	// PlayerSecret that is used to verify API request signatures (Enterprise Only).
	PlayerSecret string `json:",omitempty"`
}

// Login logs in to PlayFab account and returns LoginResult.
func (l LoginRequest) Login(ctx context.Context, client *http.Client, u *url.URL, reqBody any, opts ...internal.RequestOption) (*LoginResult, error) {
	result, err := internal.Post[*LoginResult](ctx, client, u, reqBody, opts)
	if err != nil {
		return nil, err
	}
	if !result.Valid() {
		return nil, errors.New("playfab: invalid *LoginResult result")
	}
	return result, nil
}

// LoginInfoRequest is a set of requested parameters included in LoginInfo, which can be retrieved
// through [LoginResult.InfoResult]. Users may set LoginInfoRequest as a part of LoginRequest to include
// additional parameters while signing in to PlayFab.
type LoginInfoRequest struct {
	// TODO: Add more fields when it is actually required
}

// LoginResult a session identity that can subsequently be used for API which requires an authentication.
// It is generally returned as a result of [IdentityProvider.Login] or [LoginRequest.Login].
type LoginResult struct {
	// EntityToken is an [entity.Token] of [entity.TypeTitlePlayerAccount].
	// API requests will mostly require an [entity.Token] of [entity.TypeMasterPlayerAccount]
	// so you may exchange it with [entity.Token.Exchange] with PlayFabID.
	EntityToken *entity.Token
	// InfoResult is the additional data requested by the [LoginRequest.InfoParameters].
	InfoResult LoginInfo `json:"InfoResultPayload"`
	// LastLoginTime is the time of previous login. If there was no previous login, it is zero [time.Time].
	LastLoginTime time.Time
	// NewlyCreated is true if the account was newly created on login.
	NewlyCreated bool
	// PlayFabID is the unique ID of player. It can be used for exchanging an [entity.Token]
	// of [entity.TypeMasterPlayerAccount] with [entity.Token.Exchange].
	PlayFabID string `json:"PlayFabId"`
	// SessionTicket is a unique token authorizing the user and game at server level, for the
	// current session. In Minecraft, it is used for authorizing with franchise API using a PlayFab token.
	SessionTicket string
}

// Valid returns whether the LoginResult is valid.
func (r *LoginResult) Valid() bool {
	return r != nil && r.EntityToken != nil && r.SessionTicket != "" && r.PlayFabID != ""
}

// LoginInfo represents the additional data requested by LoginInfoRequest.
type LoginInfo struct {
	// TODO: Add more fields when it is actually required
}
