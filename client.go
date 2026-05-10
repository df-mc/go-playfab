package playfab

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/df-mc/go-playfab/catalog"
	"github.com/df-mc/go-playfab/entity"
	"github.com/df-mc/go-playfab/internal"
	"github.com/df-mc/go-playfab/title"
	"golang.org/x/text/language"
)

// Login signs in to the PlayFab account using the given [IdentityProvider], and returns a new Client.
//
// The [ClientConfig]  may be used to customize the behavior of the resulting Client.
func Login(ctx context.Context, t title.Title, idp IdentityProvider, config ClientConfig) (*Client, error) {
	if config.HTTPClient == nil {
		config.HTTPClient = http.DefaultClient
	}
	if config.Logger == nil {
		config.Logger = slog.Default()
	}

	client := &Client{
		client: config.HTTPClient,
		title:  t,
		config: config,

		idp: idp,
	}
	result, err := client.login(ctx)
	if err != nil {
		return nil, err
	}
	client.newlyCreated = result.NewlyCreated
	client.ctx, client.cancel = context.WithCancelCause(context.Background())
	tokenCtx := context.WithValue(client.ctx, internal.HTTPClient, client.client)
	client.titlePlayerAccount = entity.ExchangeTokenSource(tokenCtx, t, result.EntityToken, result.EntityToken.Entity, config.Logger)
	client.masterPlayerAccount = entity.ExchangeTokenSource(tokenCtx, t, result.EntityToken, entity.Key{
		Type: entity.TypeMasterPlayerAccount,
		ID:   result.PlayFabID,
	}, config.Logger)

	// Not a smart way but we can at least check if the background task is dead.
	go client.background(client.titlePlayerAccount.Context())
	go client.background(client.masterPlayerAccount.Context())

	client.catalog = catalog.New(client.client, t, client.MasterPlayerAccount())

	return client, nil
}

// RequestOption specifies an option to be applied to an outgoing HTTP request.
//
// Callers may accept multiple RequestOptions as a variadic or slice parameter.
// Options must be applied to the request using [Apply].
//
// A RequestOption must be reusable and must not hold any per-request state.
type RequestOption = internal.RequestOption

// AcceptLanguage returns a [internal.RequestOption] that appends the given
// language tags to the 'Accept-Language' header on outgoing requests,
// preserving any tags already present in the header.
func AcceptLanguage(tags []language.Tag) RequestOption {
	return internal.AcceptLanguage(tags)
}

// RequestHeader returns a [internal.RequestOption] that sets a request header
// with the given name and value on outgoing requests.
func RequestHeader(key, value string) RequestOption {
	return internal.RequestHeader(key, value)
}

// Client implements an API client for PlayFab.
type Client struct {
	client *http.Client
	title  title.Title
	config ClientConfig

	titlePlayerAccount  entity.TokenSource
	masterPlayerAccount entity.TokenSource

	catalog *catalog.Client

	idp IdentityProvider

	loginResult *LoginResult
	loginTime   time.Time
	loginMu     sync.RWMutex

	newlyCreated bool

	ctx    context.Context
	cancel context.CancelCauseFunc
	once   sync.Once
}

func (c *Client) background(ctx context.Context) {
	select {
	case <-ctx.Done():
		c.cancel(context.Cause(ctx))
	case <-c.ctx.Done():
		return
	}
}

// TitlePlayerAccount returns an [entity.TokenSource] that supplies entity tokens for [entity.TypeTitlePlayerAccount].
func (c *Client) TitlePlayerAccount() entity.TokenSource {
	return c.titlePlayerAccount
}

// MasterPlayerAccount returns an [entity.TokenSource] that supplies entity tokens for [entity.TypeMasterPlayerAccount].
func (c *Client) MasterPlayerAccount() entity.TokenSource {
	return c.masterPlayerAccount
}

// Catalog returns an API client for PlayFab's Catalog API.
func (c *Client) Catalog() *catalog.Client {
	return c.catalog
}

// LoginInfo returns the supplementary information from the most recent login result.
func (c *Client) LoginInfo() LoginInfo {
	c.loginMu.RLock()
	defer c.loginMu.RUnlock()
	return c.loginResult.InfoResult
}

// NewlyCreated reports whether a PlayFab account was newly created during the initial login.
func (c *Client) NewlyCreated() bool {
	return c.newlyCreated
}

// LastLoginTime returns the time of the caller's most recent previous login.
func (c *Client) LastLoginTime() time.Time {
	c.loginMu.RLock()
	defer c.loginMu.RUnlock()
	return c.loginResult.LastLoginTime
}

// SessionTicket returns the session ticket from the current login result.
// If the cached login result has expired (24 hours after the previous login),
// the provided [context.Context] is used to re-authenticate before returning.
func (c *Client) SessionTicket(ctx context.Context) (string, error) {
	result, err := c.login(ctx)
	if err != nil {
		return "", fmt.Errorf("login: %w", err)
	}
	return result.SessionTicket, nil
}

const (
	// loginExpiration is the duration after which a LoginResult is considered expired.
	// The Client records the timestamp of each login and uses this duration together
	// with loginExpirationDelta to determine whether a cached LoginResult is still valid.
	loginExpiration = time.Hour * 24
	// loginExpirationDelta is subtracted from loginExpiration to add a safety margin
	// when deciding whether a cached login result should be refreshed.
	loginExpirationDelta = time.Minute * 15
)

// login authenticates with PlayFab using the Client's identity provider.
// Results are cached internally and reused until they expire.
func (c *Client) login(ctx context.Context) (*LoginResult, error) {
	c.loginMu.Lock()
	defer c.loginMu.Unlock()

	if c.loginResult != nil && c.loginResult.Valid() && time.Now().Before(c.loginTime.Add(loginExpiration-loginExpirationDelta)) {
		return c.loginResult, nil
	}

	result, err := c.idp.Login(ctx, c.client, c.config.login(c.title))
	if err != nil {
		return nil, fmt.Errorf("login: %w", err)
	}
	c.loginResult, c.loginTime = result, time.Now()
	return result, nil
}

// Close closes the Client. Once the Client is closed, the entity tokens are no longer
// exchanged in background as the internal context is closed.
func (c *Client) Close() (err error) {
	return c.close(net.ErrClosed)
}

// close cancels the background context with the provided cause.
func (c *Client) close(cause error) (err error) {
	c.once.Do(func() {
		if cause == nil {
			cause = net.ErrClosed
		}
		c.config.Logger.Debug("client is closing", slog.Any("cause", cause))
		c.cancel(cause)
	})
	return err
}

// IdentityProvider is the interface for logging in to a PlayFab account through
// an underlying identity provider, such as Xbox Live.
//
// Implementations must internally cache and refresh credentials obtained from
// the external identity provider so that Login is always called with a current token.
type IdentityProvider interface {
	// Login authenticates with PlayFab using the provided [LoginRequest].
	// Implementations may first authenticate with an external identity provider
	// and embed the resulting token in the request before logging in to the
	// associated PlayFab account. Implementations do not need to cache the
	// returned [LoginResult] because the Client caches it internally.
	Login(ctx context.Context, client *http.Client, request LoginRequest) (*LoginResult, error)
}

// ClientConfig contains options to configure a Client when logging in to a PlayFab account.
type ClientConfig struct {
	// HTTPClient is the HTTP client through all PlayFab API requests are sent.
	// Defaults to [http.DefaultClient] if nil.
	HTTPClient *http.Client
	// Logger receives log output at various levels during token exchange and authentication.
	// Defaults to [slog.Default] if nil.
	Logger *slog.Logger

	// CreateAccount specifies whether to create a new PlayFab account
	// if one does not already exist for the given identity.
	CreateAccount bool
}

// login makes a [LoginRequest] from the ClientConfig for the given title.
func (c ClientConfig) login(t title.Title) LoginRequest {
	return LoginRequest{
		Title:         t,
		CreateAccount: c.CreateAccount,
	}
}
