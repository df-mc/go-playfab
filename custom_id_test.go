package playfab

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// redirectingClient returns an *http.Client that rewrites every request aimed
// at '<title>.playfabapi.com' to the given test server.
func redirectingClient(t *testing.T, server *httptest.Server) *http.Client {
	t.Helper()

	target, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("parse server URL: %v", err)
	}
	return &http.Client{Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		req = req.Clone(req.Context())
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		return http.DefaultTransport.RoundTrip(req)
	})}
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

const okLoginResult = `{"code":200,"status":"OK","data":{"SessionTicket":"ticket","PlayFabId":"ID",` +
	`"EntityToken":{"EntityToken":"entity","Entity":{"Id":"ID","Type":"title_player_account",` +
	`"TypeString":"title_player_account"},"TokenExpiration":"3000-01-01T00:00:00Z"}}}`

func TestCustomIDIdentityProviderLogin(t *testing.T) {
	t.Parallel()

	t.Run("PlainCustomID", func(t *testing.T) {
		t.Parallel()

		var got map[string]any
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/Client/LoginWithCustomID" {
				t.Errorf("unexpected path: %s", r.URL.Path)
			}
			b, _ := io.ReadAll(r.Body)
			if err := json.Unmarshal(b, &got); err != nil {
				t.Errorf("decode request body: %v", err)
			}
			_, _ = io.WriteString(w, okLoginResult)
		}))
		defer server.Close()

		provider := CustomIDIdentityProvider{CustomID: "anonymous-device"}
		result, err := provider.Login(context.Background(), redirectingClient(t, server), LoginRequest{Title: "20CA2"})
		if err != nil {
			t.Fatalf("login: %v", err)
		}
		if !result.Valid() {
			t.Fatal("result is not valid")
		}
		if got["CustomId"] != "anonymous-device" {
			t.Errorf("CustomId = %v, want anonymous-device", got["CustomId"])
		}
		if _, ok := got["EncryptedRequest"]; ok {
			t.Error("EncryptedRequest should be absent for a plain custom ID")
		}
	})

	t.Run("PlayerSecret", func(t *testing.T) {
		t.Parallel()

		var got map[string]any
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			b, _ := io.ReadAll(r.Body)
			if err := json.Unmarshal(b, &got); err != nil {
				t.Errorf("decode request body: %v", err)
			}
			_, _ = io.WriteString(w, okLoginResult)
		}))
		defer server.Close()

		// ClientConfig builds the LoginRequest without one, so the provider is
		// the only way a caller can supply a secret.
		provider := CustomIDIdentityProvider{CustomID: "anonymous-device", PlayerSecret: "top-secret"}
		result, err := provider.Login(context.Background(), redirectingClient(t, server), LoginRequest{Title: "20CA2"})
		if err != nil {
			t.Fatalf("login: %v", err)
		}
		if !result.Valid() {
			t.Fatal("result is not valid")
		}
		if got["PlayerSecret"] != "top-secret" {
			t.Errorf("PlayerSecret = %v, want top-secret", got["PlayerSecret"])
		}
	})

	t.Run("NoPlayerSecret", func(t *testing.T) {
		t.Parallel()

		var got map[string]any
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			b, _ := io.ReadAll(r.Body)
			if err := json.Unmarshal(b, &got); err != nil {
				t.Errorf("decode request body: %v", err)
			}
			_, _ = io.WriteString(w, okLoginResult)
		}))
		defer server.Close()

		provider := CustomIDIdentityProvider{CustomID: "anonymous-device"}
		if _, err := provider.Login(context.Background(), redirectingClient(t, server), LoginRequest{Title: "20CA2"}); err != nil {
			t.Fatalf("login: %v", err)
		}
		if _, ok := got["PlayerSecret"]; ok {
			t.Errorf("PlayerSecret was sent for a provider without one: %v", got["PlayerSecret"])
		}
	})

	t.Run("EncryptedRequest", func(t *testing.T) {
		t.Parallel()

		var encrypted string
		var plain, outer map[string]any
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			b, _ := io.ReadAll(r.Body)
			if err := json.Unmarshal(b, &outer); err != nil {
				t.Errorf("decode request body: %v", err)
			}
			encrypted, _ = outer["EncryptedRequest"].(string)
			_, _ = io.WriteString(w, okLoginResult)
		}))
		defer server.Close()

		provider := CustomIDIdentityProvider{
			CustomID: "anonymous-device",
			Encrypt: func(payload []byte) ([]byte, error) {
				if err := json.Unmarshal(payload, &plain); err != nil {
					return nil, err
				}
				// Stand-in for the title RSA key: reverse the payload so the
				// test can assert it reached the hook untouched.
				out := make([]byte, len(payload))
				for i := range payload {
					out[len(payload)-1-i] = payload[i]
				}
				return out, nil
			},
		}
		result, err := provider.Login(context.Background(), redirectingClient(t, server), LoginRequest{Title: "20CA2"})
		if err != nil {
			t.Fatalf("login: %v", err)
		}
		if !result.Valid() {
			t.Fatal("result is not valid")
		}
		if encrypted == "" {
			t.Fatal("EncryptedRequest was not sent")
		}
		if plain["CustomId"] != "anonymous-device" {
			t.Errorf("encrypted payload CustomId = %v, want anonymous-device", plain["CustomId"])
		}
		if _, ok := plain["PlayerSecret"]; ok {
			t.Errorf("encrypted payload carried PlayerSecret = %v for a provider without one", plain["PlayerSecret"])
		}
		// A title that requires encryption must not also be handed the fields
		// in clear text.
		if _, ok := outer["CustomId"]; ok {
			t.Errorf("CustomId was also sent in clear text: %v", outer["CustomId"])
		}
		if _, ok := outer["PlayerSecret"]; ok {
			t.Errorf("PlayerSecret was also sent in clear text: %v", outer["PlayerSecret"])
		}
	})

	t.Run("EncryptedRequestCarriesPlayerSecret", func(t *testing.T) {
		t.Parallel()

		var encrypted string
		var plain, outer map[string]any
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			b, _ := io.ReadAll(r.Body)
			if err := json.Unmarshal(b, &outer); err != nil {
				t.Errorf("decode request body: %v", err)
			}
			encrypted, _ = outer["EncryptedRequest"].(string)
			_, _ = io.WriteString(w, okLoginResult)
		}))
		defer server.Close()

		provider := CustomIDIdentityProvider{
			CustomID:     "anonymous-device",
			PlayerSecret: "top-secret",
			Encrypt: func(payload []byte) ([]byte, error) {
				if err := json.Unmarshal(payload, &plain); err != nil {
					return nil, err
				}
				out := make([]byte, len(payload))
				for i := range payload {
					out[len(payload)-1-i] = payload[i]
				}
				return out, nil
			},
		}
		result, err := provider.Login(context.Background(), redirectingClient(t, server), LoginRequest{Title: "20CA2"})
		if err != nil {
			t.Fatalf("login: %v", err)
		}
		if !result.Valid() {
			t.Fatal("result is not valid")
		}
		if encrypted == "" {
			t.Fatal("EncryptedRequest was not sent")
		}
		// The secret must travel inside the encrypted payload, not beside it.
		if plain["CustomId"] != "anonymous-device" {
			t.Errorf("encrypted payload CustomId = %v, want anonymous-device", plain["CustomId"])
		}
		if plain["PlayerSecret"] != "top-secret" {
			t.Errorf("encrypted payload PlayerSecret = %v, want top-secret", plain["PlayerSecret"])
		}
		if _, ok := outer["PlayerSecret"]; ok {
			t.Errorf("PlayerSecret was also sent in clear text: %v", outer["PlayerSecret"])
		}
	})

	t.Run("EmptyCustomID", func(t *testing.T) {
		t.Parallel()

		provider := CustomIDIdentityProvider{}
		if _, err := provider.Login(context.Background(), http.DefaultClient, LoginRequest{Title: "20CA2"}); err == nil {
			t.Fatal("expected an error for an empty custom ID")
		} else if !strings.Contains(err.Error(), "CustomID") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("EncryptError", func(t *testing.T) {
		t.Parallel()

		want := errors.New("no key")
		provider := CustomIDIdentityProvider{
			CustomID: "anonymous-device",
			Encrypt:  func([]byte) ([]byte, error) { return nil, want },
		}
		_, err := provider.Login(context.Background(), http.DefaultClient, LoginRequest{Title: "20CA2"})
		if !errors.Is(err, want) {
			t.Fatalf("err = %v, want %v", err, want)
		}
	})
}
