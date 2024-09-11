package catalog

import (
	"fmt"
	"github.com/df-mc/go-playfab/entity"
	"github.com/df-mc/go-playfab/internal"
	"github.com/df-mc/go-playfab/title"
	"golang.org/x/text/language"
	"net/http"
)

type Filter struct {
	// Count is the number of returned items included in the SearchResult.
	// The maximum value is 50, and defaulted to 10 by the service-side.
	Count int `json:"Count,omitempty"`
	// ContinuationToken is the opaque token used for continuing the Search,
	// if any are available. It is normally filled from [SearchResult.ContinuationToken].
	ContinuationToken string `json:"ContinuationToken,omitempty"`
	// CustomTags is the optional properties associated with the request.
	CustomTags map[string]any `json:"CustomTags,omitempty"`
	// Entity is an [entity.Key] to perform any actions using the Filter.
	// If left as nil and an [entity.Token] has been provided to [Filter.Search],
	// it will be filled from [entity.Token.Key].
	Entity *entity.Key `json:"Entity,omitempty"`
	// Filter is an OData query for filtering the items included in SearchResult.
	// For example, "<Field of Item> eq '<Value of Item>'".
	Filter string `json:"Filter,omitempty"`
	// Language is the locale to be included in the dictionary of Items returned in
	// the SearchResult. It is also used as an 'Accept-Language' header of the request
	// sent from [SearchResult.Search].
	Language language.Tag `json:"Language,omitempty"`
	// OrderBy is an OData sort query for sorting the index of SearchResult. Defaulted to relevance.
	OrderBy string `json:"OrderBy,omitempty"`
	// Term is the string terms to be searched.
	Term string `json:"Search,omitempty"`
	// Select is an OData selection query for filtering the fields of returned items included in the SearchResult.
	Select string `json:"Select,omitempty"`
	// Store ...
	Store *StoreReference `json:"Store,omitempty"`
}

// Search performs a search against the public catalog using the Filter and returns a set of
// paginated SearchResult. It uses a cache of the catalog with item updates taking up to few
// minutes to propagate. If trying to immediately retrieve recent Item updates, a Query should
// be used. More information about the Search API can be found here:
// https://learn.microsoft.com/en-us/gaming/playfab/features/economy-v2/catalog/search
func (f Filter) Search(t title.Title, tok *entity.Token) (*SearchResult, error) {
	if f.Count > 50 {
		return nil, fmt.Errorf("playfab/catalog: Filter: count must be <= 50, got %d", f.Count)
	}
	if f.Entity == nil && tok != nil {
		f.Entity = &tok.Entity
	}

	return internal.Post[*SearchResult](t.URL().JoinPath("/Catalog/SearchItems"), f, func(req *http.Request) {
		if tok != nil {
			tok.SetAuthHeader(req)
		}
		if f.Language != language.Und {
			req.Header.Set("Accept-Language", f.Language.String())
		}
	})
}

type SearchResult struct {
	// ContinuationToken provides an opaque token for continuing the next page of
	// SearchResult by specifying it to [Filter.ContinuationToken], if any are available.
	ContinuationToken string `json:"ContinuationToken,omitempty"`
	// Items is a paginated set of Item for the search query.
	Items []Item `json:"Items,omitempty"`
}
