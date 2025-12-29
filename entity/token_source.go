package entity

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/df-mc/go-playfab/title"
)

// A TokenSource continuously supplies an entity token.
type TokenSource interface {
	Token() (*Token, error)
}

// ExchangeTokenSource returns a TokenSource that continuously exchanges the entity token
// in background.
func ExchangeTokenSource(ctx context.Context, tok *Token, t title.Title, masterID string) TokenSource {
	src := &exchangeTokenSource{
		tok: tok,

		title:    t,
		masterID: masterID,
	}
	go src.background(ctx)
	return src
}

// exchangeTokenSource continuously supplies a fresh token by exchanging them at
// certain interval in background.
type exchangeTokenSource struct {
	tok *Token
	err error

	mux      sync.Mutex
	title    title.Title
	masterID string
}

// background continuously exchanges the entity token at 15 minutes interval.
func (src *exchangeTokenSource) background(ctx context.Context) {
	t := time.NewTicker(time.Minute * 15)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			src.mux.Lock()
			src.tok, src.err = src.tok.Exchange(src.title, src.masterID)
			if src.err != nil {
				src.mux.Unlock()
				return
			}
			src.mux.Unlock()
		case <-ctx.Done():
			return
		}
	}
}

// Token supplies a fresh token that is exchanged at certain interval.
func (src *exchangeTokenSource) Token() (tok *Token, err error) {
	src.mux.Lock()
	defer src.mux.Unlock()
	if src.err != nil {
		return nil, fmt.Errorf("exchange token in background: %w", err)
	}

	if src.tok.Expired() || src.tok.Entity.Type != TypeMasterPlayerAccount {
		tok, err = src.tok.Exchange(src.title, src.masterID)
		if err != nil {
			return nil, fmt.Errorf("exchange: %w", err)
		}
		src.tok = tok
	}
	return src.tok, nil
}
