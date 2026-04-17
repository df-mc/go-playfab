package catalog

import (
	"encoding/json"
	"fmt"
	"strings"

	"golang.org/x/text/language"
)

type Dictionary[T any] map[string]T

func (d *Dictionary[T]) Lookup(key string) (zero T, ok bool) {
	for k, value := range *d {
		if strings.EqualFold(k, key) {
			return value, true
		}
	}
	return zero, false
}

func (d *Dictionary[T]) Neutral() (value T) {
	for key, v := range *d {
		if value = v; strings.EqualFold(key, "neutral") {
			break
		}
	}
	return value
}

func (d *Dictionary[T]) UnmarshalJSON(b []byte) error {
	type Alias Dictionary[T]
	if err := json.Unmarshal(b, (*Alias)(d)); err != nil {
		return err
	}
	for key := range *d {
		if strings.EqualFold(key, "neutral") {
			continue
		}
		if _, err := language.Parse(key); err != nil {
			return fmt.Errorf("catalog: Dictionary: parse %q as language tag: %w", key, err)
		}
	}
	// if !hasNeutral {
	//	return errors.New("catalog: Dictionary must contain a value for NEUTRAL key")
	// }
	return nil
}
