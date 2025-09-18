package retrymanager

import (
	"context"
	"time"
)

type TryFunc func() bool

type Opts struct {
	NumTries   int
	RetryAfter time.Duration
}

type RetryManager struct {
	NumTries   int
	RetryAfter time.Duration
}

func New(opts Opts) *RetryManager {
	return &RetryManager{
		NumTries:   opts.NumTries,
		RetryAfter: opts.RetryAfter,
	}
}

func (rm *RetryManager) Do(f TryFunc) {
	retries := 0
	for retries < rm.NumTries {
		retries++
		if f() {
			return
		}
		time.Sleep(rm.RetryAfter)
	}
}

func (rm *RetryManager) DoWithCtx(ctx context.Context, f TryFunc) error {
	for {
		if f() {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(rm.RetryAfter):
		}
	}
}
