package entity

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/df-mc/go-playfab/internal"
	"github.com/df-mc/go-playfab/title"
)

type TokenSource interface {
	EntityToken(ctx context.Context) (*Token, error)
}

func ExchangeTokenSource(ctx context.Context, title title.Title, token *Token, key Key, log *slog.Logger) TokenSource {
	if token == nil {
		panic("entity: ReuseTokenSource: *entity.Token cannot be nil")
	}
	if log == nil {
		log = slog.Default()
	}

	r := &reuseTokenSource{
		title: title,
		key:   key,

		log: log,

		t: token,
	}
	var cancel context.CancelCauseFunc
	r.ctx, cancel = context.WithCancelCause(context.WithValue(ctx, internal.HTTPClient, internal.ContextClient(ctx)))
	go r.background(cancel)
	return r
}

type reuseTokenSource struct {
	title title.Title
	key   Key

	log *slog.Logger

	ctx context.Context

	t  *Token
	mu sync.Mutex
}

func (r *reuseTokenSource) background(cancel context.CancelCauseFunc) {
	r.mu.Lock()
	exp := r.t.Expiration
	r.mu.Unlock()

	for {
		select {
		case <-time.After(time.Until(exp.Add(-time.Minute * 20))):
			r.mu.Lock()
			token, err := r.t.Exchange(r.ctx, r.title, r.key)
			if err != nil {
				r.mu.Unlock()
				r.log.Error("error exchanging token", slog.Any("error", err))
				cancel(fmt.Errorf("exchange token in background: %w", err))
				return
			}
			r.t = token
			exp = token.Expiration
			r.mu.Unlock()
			r.log.Debug("exchanged entity token in background", slog.Any("entity", r.key))
		case <-r.ctx.Done():
			return
		}
	}
}

func (r *reuseTokenSource) EntityToken(ctx context.Context) (*Token, error) {
	if r.ctx.Err() != nil {
		return nil, context.Cause(r.ctx)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if r.t.Entity == r.key && r.t.Valid() {
		return r.t, nil
	}

	token, err := r.t.Exchange(ctx, r.title, r.key)
	if err != nil {
		return nil, fmt.Errorf("exchange: %w", err)
	}
	r.t = token
	return token, nil
}
