package title

import (
	"net/url"
	"strings"
)

type Title string

func (t Title) URL() *url.URL {
	return &url.URL{
		Scheme: "https",
		Host:   strings.ToLower(string(t)) + ".playfabapi.com",
	}
}
