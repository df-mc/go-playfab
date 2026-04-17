package playfab

import (
	"context"
	"fmt"
	"net/http"

	"github.com/df-mc/go-playfab/title"
	"github.com/df-mc/go-xsapi/v2"
)

func LoginWithXbox(ctx context.Context, t title.Title, client *xsapi.Client, config ClientConfig) (*Client, error) {
	return Login(ctx, t, &XBLIdentityProvider{Client: client}, config)
}

type XBLIdentityProvider struct {
	Client *xsapi.Client
}

func (i XBLIdentityProvider) Login(ctx context.Context, client *http.Client, request LoginRequest) (*LoginResult, error) {
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

type loginWithXbox struct {
	LoginRequest
	XboxToken string
}
