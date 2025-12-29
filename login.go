package playfab

import (
	"encoding/json"
	"errors"

	"github.com/df-mc/go-playfab/internal"
	"github.com/df-mc/go-playfab/title"
)

// LoginConfig represents a base structure used for signing in to the PlayFab through a method.
// An implementation of IdentityProvider may embed LoginConfig with additional fields required
// to fill in.
type LoginConfig struct {
	// Title indicates the title ID specific to the PlayFab.
	Title title.Title `json:"TitleId,omitempty"`
	// CreateAccount specifies to automatically create a PlayFab account if
	// one is not currently linked to the ID.
	CreateAccount bool `json:"CreateAccount,omitempty"`
	// CustomTags is the optional custom tags associated with the request.
	CustomTags map[string]any `json:"CustomTags,omitempty"`
	// EncryptedRequest is a base64-encoded body that is encrypted with the
	// Title's public RSA key (Enterprise Only).
	EncryptedRequest []byte `json:"EncryptedRequest,omitempty"`
	// RequestParameters can be used to request additional data in the sign-in process.
	// If the sign-in was successful, the relevant data requested through this field
	// will be available in [Identity.ResponseParameters].
	RequestParameters json.RawMessage `json:"InfoRequestParameters,omitempty"`
	// PlayerSecret that is used to verify API request signatures (Enterprise Only).
	PlayerSecret string `json:"PlayerSecret,omitempty"`
}

// IdentityProvider implements a Login method, which signs in to the PlayFab using [LoginConfig.Login] with a LoginConfig
// with additional fields required to fill in, through the path '/Client/LoginWithXXX', where X is normally the method to
// sign in. IdentityProvider is implemented by several platforms which supports signing in to the PlayFab with their identity.
type IdentityProvider interface {
	Identity() (*Identity, error)
}

// Login signs in to PlayFab using the request body and the path. The path normally follows the format
// 'Client/LoginWithXXX' where X typically is the method for signing in. The request body is generally
// a structure that may embed LoginConfig with additional fields required for the specific path, such
// as a token of IdentityProvider.
func (l LoginConfig) Login(path string, body any) (*Identity, error) {
	if l.Title == "" {
		return nil, errors.New("playfab: LoginConfig: Title not set")
	}
	return internal.Post[*Identity](l.Title.URL().JoinPath(path), body)
}
