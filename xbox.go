package playfab

import (
	"context"
	"fmt"
	"net/http"

	"github.com/df-mc/go-playfab/title"
	"github.com/df-mc/go-xsapi/v2"
)

// LoginWithXbox logs in to PlayFab account in the specified title ID with an Xbox Live account.
func LoginWithXbox(ctx context.Context, t title.Title, client *xsapi.Client, config ClientConfig) (*Client, error) {
	return Login(ctx, t, &XBLIdentityProvider{Client: client}, config)
}

// XBLIdentityProvider implements an [IdentityProvider] that logs in to PlayFab account
// with an Xbox Live account using the underlying Client.
type XBLIdentityProvider struct {
	// Client is the Xbox Live API Client used to log in to Xbox Live services.
	Client *xsapi.Client
}

// Login ...
func (i XBLIdentityProvider) Login(ctx context.Context, client *http.Client, request LoginRequest) (*LoginResult, error) {
	if i.Client == nil {
		panic("playfab: XBLIdentityProvider.Client cannot be nil")
	}
	requestURL := request.Title.URL().JoinPath("/Client/LoginWithXbox")
	token, _, err := i.Client.TokenAndSignature(ctx, requestURL)
	if err != nil {
		return nil, fmt.Errorf("request XSTS token and signature: %w", err)
	}
	return request.Login(ctx, client, requestURL, loginWithXbox{
		LoginRequest: request,
		XboxToken:    token.String(),
	})
}

// loginWithXbox is a payload for logging in to PlayFab account with an Xbox Live account.
type loginWithXbox struct {
	LoginRequest
	// XboxToken is the XSTS token in string form.
	XboxToken string
}
