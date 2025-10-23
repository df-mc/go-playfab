package playfab

import (
	"errors"
	"fmt"
	"github.com/df-mc/go-xsapi"
)

// XBLIdentityProvider implements IdentityProvider for signing in to PlayFab using the path
// '/Client/LoginWithXbox' with a [xsapi.Token] that relies on the party 'http://playfab.xboxlive.com/'.
type XBLIdentityProvider struct {
	// TokenSource is used for obtaining a [xsapi.Token] that relies on the party defined in
	// the constant below (see RelyingParty).
	TokenSource xsapi.TokenSource
}

// Login signs in to PlayFab using a [xsapi.Token] obtained from the TokenSource that relies on the
// party 'http://playfab.xboxlive.com/'. It returns an Identity by calling the [LoginConfig.Login] method
// with an additional field named 'XboxToken' that specifies a string obtained from [xsapi.Token.String]
// and the path '/Client/LoginWithXbox'.
func (prov XBLIdentityProvider) Login(config LoginConfig) (*Identity, error) {
	if prov.TokenSource == nil {
		return nil, errors.New("playfab: XBLIdentityProvider: TokenSource is nil")
	}

	tok, err := prov.TokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("request xbox live token: %w", err)
	}

	type loginConfig struct {
		LoginConfig
		XboxToken string `json:"XboxToken"`
	}
	return config.Login("/Client/LoginWithXbox", loginConfig{
		LoginConfig: config,
		XboxToken:   tok.String(),
	})
}

// RelyingParty is the party that a [xsapi.Token] obtained from [XBLIdentityProvider.TokenSource] should rely
// on. Using a [xsapi.Token] that relies on other party may cause an error related to "decrypting token body".
const RelyingParty = "http://playfab.xboxlive.com/"
