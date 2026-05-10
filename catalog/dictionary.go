package catalog

import (
	"encoding/json"
	"fmt"
	"strings"

	"golang.org/x/text/language"
)

// A Dictionary represents a localized table of values.
// Each Dictionary has a Neutral value and other values
// keyed with language tags.
type Dictionary[T any] map[string]T

// Lookup looks up for the value with the key.
func (d *Dictionary[T]) Lookup(key string) (zero T, ok bool) {
	for k, value := range *d {
		if strings.EqualFold(k, key) {
			return value, true
		}
	}
	return zero, false
}

// Neutral returns a neutral text for the Dictionary.
// When a text with the 'NEUTRAL' key was not found in
// the Dictionary, it returns the last value found in Dictionary.
func (d *Dictionary[T]) Neutral() (value T) {
	for key, v := range *d {
		if value = v; strings.EqualFold(key, "neutral") {
			break
		}
	}
	return value
}

// UnmarshalJSON ...
func (d *Dictionary[T]) UnmarshalJSON(b []byte) error {
	type Alias Dictionary[T]
	if err := json.Unmarshal(b, (*Alias)(d)); err != nil {
		return err
	}
	if *d == nil {
		return fmt.Errorf("invalid dictionary JSON: %s", b)
	}
	for key := range *d {
		if strings.EqualFold(key, "neutral") {
			continue
		}
		if _, err := language.Parse(key); err != nil {
			return fmt.Errorf("catalog: Dictionary: parse %q as language tag: %w", key, err)
		}
	}
	return nil
}
