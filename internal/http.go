package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func Post[T any](u *url.URL, r any, hooks ...func(req *http.Request)) (zero T, err error) {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(r); err != nil {
		return zero, fmt.Errorf("encode: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, u.String(), buf)
	if err != nil {
		return zero, fmt.Errorf("make request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for _, hook := range hooks {
		hook(req)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return zero, fmt.Errorf("POST %s: %w", u, err)
	}
	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		var body Result[T]
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			return zero, fmt.Errorf("decode: %w", err)
		}
		return body.Data, nil
	default:
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return zero, fmt.Errorf("POST %s: %s", u, resp.Status)
		}
		var body Error
		if err := json.Unmarshal(b, &body); err != nil {
			return zero, fmt.Errorf("POST %s: %s: %s (%w)", u, resp.Status, b, err)
		}
		return zero, body
	}
}
