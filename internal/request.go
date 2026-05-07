package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/text/language"
)

// RequestOption specifies an option to be applied to an outgoing HTTP request.
//
// Callers may accept multiple RequestOptions as a variadic or slice parameter.
// Options must be applied to the request using [Apply].
//
// A RequestOption must be reusable and must not hold any per-request state.
type RequestOption func(req *http.Request) error

// AcceptLanguage returns a [internal.RequestOption] that appends the given
// language tags to the 'Accept-Language' header on outgoing requests,
// preserving any tags already present in the header.
func AcceptLanguage(tags []language.Tag) RequestOption {
	s := make([]string, len(tags))
	for i, tag := range tags {
		s[i] = tag.String()
	}
	return func(req *http.Request) error {
		req.Header.Add("Accept-Language", strings.Join(s, ", "))
		return nil
	}
}

// RequestHeader returns a [internal.RequestOption] that sets a request header
// with the given name and value on outgoing requests.
func RequestHeader(key, value string) RequestOption {
	return func(req *http.Request) error {
		req.Header.Set(key, value)
		return nil
	}
}

// Post issues a POST request to the endpoint.
func Post[T any](ctx context.Context, client *http.Client, u *url.URL, reqBody any, opts []RequestOption) (value T, err error) {
	var r io.Reader
	if reqBody != nil {
		buf := &bytes.Buffer{}
		defer buf.Reset()
		if err := json.NewEncoder(buf).Encode(reqBody); err != nil {
			return value, fmt.Errorf("encode request body: %w", err)
		}
		r = buf
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), r)
	if err != nil {
		return value, fmt.Errorf("make request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if err := Apply(req, opts); err != nil {
		return value, fmt.Errorf("apply request options: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return value, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		buf := &bytes.Buffer{}
		if _, err := buf.ReadFrom(resp.Body); err != nil {
			return value, fmt.Errorf("read response body: %w", err)
		}
		resp.Body = io.NopCloser(buf)

		var result Result[T]
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return value, fmt.Errorf("decode response body: %w", err)
		}
		return result.Data, nil
	default:
		b, err := io.ReadAll(io.LimitReader(resp.Body, 1024))
		if err != nil {
			return value, fmt.Errorf("read response body: %w", err)
		}
		var e Error
		if err := json.Unmarshal(b, &e); err != nil {
			return value, fmt.Errorf("%s %s: %s (%s)", req.Method, req.URL, resp.Status, b)
		}
		return value, &e
	}
}

// contextKey is the type used for defining a context key.
type contextKey struct{}

// HTTPClient is the context key used to specify the HTTP client
// used to issue the request.
var HTTPClient contextKey

// DefaultLanguage is the default language applied by default.
var DefaultLanguage = []language.Tag{
	language.AmericanEnglish,
	language.English,
}

// ContextClient returns an HTTP client based on the given [context.Context].
func ContextClient(ctx context.Context) *http.Client {
	if hc, ok := ctx.Value(HTTPClient).(*http.Client); ok {
		return hc
	}
	return http.DefaultClient
}

// UnexpectedStatusCode returns an error describing an unexpected HTTP status code,
// including the request method and URL for context.
// The resp must be a client response because [http.Response.Request] is only
// populated on responses received by the client.
func UnexpectedStatusCode(resp *http.Response) error {
	return fmt.Errorf("%s %s: %s", resp.Request.Method, resp.Request.URL, resp.Status)
}

// Apply applies the given RequestOptions to the request in order.
// Caller-provided opts take precedence over any defaults appended after them.
// For example, append caller opts before defaults like append(opts, internal.DefaultLanguage)
// so that the caller's preferences are evaluated first.
func Apply(req *http.Request, opts []RequestOption) error {
	for _, opt := range opts {
		if err := opt(req); err != nil {
			return err
		}
	}
	return nil
}
