package title

import (
	"net/url"
	"strings"
)

// Title represents a PlayFab title. The string itself is a hexadecimal ID of the title.
type Title string

// URL returns an [url.URL] of the Title. It is generally called for sending a request to
// the API for the Title. It follows the format 'https://XXX.playfabapi.com', where X is
// the lowercase ID of title.
func (t Title) URL() *url.URL {
	return &url.URL{
		Scheme: "https",
		Host:   strings.ToLower(string(t)) + ".playfabapi.com",
	}
}
