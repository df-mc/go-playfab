package playfab

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// CustomIDIdentityProvider implements an [IdentityProvider] that logs in to
// PlayFab with an anonymous custom ID, without an external platform token.
type CustomIDIdentityProvider struct {
	// CustomID is the device-bound identifier of the anonymous account. It
	// must not be empty.
	CustomID string

	// PlayerSecret is an optional secret used to verify API request
	// signatures. It is an Enterprise-only feature, and is ignored by titles
	// that do not have one configured. It is carried inside the encrypted
	// request when Encrypt is set, so a title that requires encryption does not
	// force the secret to be sent in clear text.
	//
	// It is a field on the provider rather than on ClientConfig because
	// ClientConfig builds the LoginRequest without one, so a caller would
	// otherwise have no way to supply it.
	PlayerSecret string

	// Encrypt optionally encrypts the request payload with the title's public
	// RSA key. Titles may require encrypted requests, in which case a plain
	// CustomID is rejected with NotAuthorizedByTitle. When nil, the CustomID is
	// sent in plain text.
	//
	// The function receives the JSON encoded payload containing the fields the
	// encrypted request carries, currently CustomId and PlayerSecret, and must
	// return the encrypted bytes to send as the EncryptedRequest field.
	Encrypt func(payload []byte) ([]byte, error)
}

// Login implements [IdentityProvider.Login].
func (i CustomIDIdentityProvider) Login(ctx context.Context, client *http.Client, request LoginRequest) (*LoginResult, error) {
	if i.CustomID == "" {
		return nil, errors.New("playfab: CustomIDIdentityProvider.CustomID cannot be empty")
	}
	requestURL := request.Title.URL().JoinPath("/Client/LoginWithCustomID")

	payload := loginWithCustomID{
		LoginRequest: request,
		CustomID:     i.CustomID,
	}
	if i.PlayerSecret != "" {
		payload.PlayerSecret = i.PlayerSecret
	}
	if i.Encrypt != nil {
		// The encryption hook takes the payload PlayFab would otherwise receive
		// in clear text, so build it first and only then encrypt. Everything the
		// title treats as sensitive goes inside it, which for LoginWithCustomID
		// means the custom ID and the player secret.
		encryptedBody := struct {
			CustomID     string `json:"CustomId"`
			PlayerSecret string `json:"PlayerSecret,omitempty"`
		}{CustomID: i.CustomID, PlayerSecret: i.PlayerSecret}
		b, err := json.Marshal(encryptedBody)
		if err != nil {
			return nil, fmt.Errorf("playfab: encode custom ID payload: %w", err)
		}
		encrypted, err := i.Encrypt(b)
		if err != nil {
			return nil, fmt.Errorf("playfab: encrypt custom ID payload: %w", err)
		}
		// A title that requires encryption rejects a plain CustomID, so the
		// clear-text fields are dropped once they are encrypted. Leaving the
		// secret in the clear would defeat the point of sending it encrypted.
		payload = loginWithCustomID{LoginRequest: request}
		payload.EncryptedRequest = encrypted
	}
	return request.Login(ctx, client, requestURL, payload)
}

// loginWithCustomID is a payload for logging in to a PlayFab account with an
// anonymous custom ID. EncryptedRequest and CustomID are mutually exclusive:
// PlayFab ignores CustomID when EncryptedRequest is present.
type loginWithCustomID struct {
	LoginRequest

	// CustomID is the device-bound identifier of the anonymous account. It is
	// omitted when the payload is sent as an EncryptedRequest.
	CustomID string `json:"CustomId,omitempty"`

	// EncryptedRequest is the base64-encoded payload encrypted with the
	// title's public RSA key.
	EncryptedRequest []byte `json:",omitempty"`
}
