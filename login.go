package playfab

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/df-mc/go-playfab/entity"
	"github.com/df-mc/go-playfab/internal"
	"github.com/df-mc/go-playfab/title"
)

type LoginRequest struct {
	Title            title.Title       `json:"TitleId"`
	CreateAccount    bool              `json:",omitempty"`
	CustomTags       map[string]any    `json:",omitempty"`
	EncryptedRequest []byte            `json:",omitempty"`
	InfoParameters   *LoginInfoRequest `json:"InfoRequestParameters,omitempty"`
	PlayerSecret     string            `json:",omitempty"`
}

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

type LoginInfoRequest struct {
}

type LoginResult struct {
	EntityToken   *entity.Token
	InfoResult    LoginInfo
	LastLoginTime time.Time
	NewlyCreated  bool
	PlayFabID     string `json:"PlayFabId"`
	SessionTicket string
}

func (r *LoginResult) Valid() bool {
	return r != nil && r.EntityToken != nil && r.SessionTicket != "" && r.PlayFabID != ""
}

type LoginInfo struct {
}
