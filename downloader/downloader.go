package downloader

import (
	"context"
	"errors"
)

var (
	ErrDurationLimitExceeded = errors.New("duration limit exceeded")
)

type Service interface {
	IsValidURL(url string) bool
	Download(ctx context.Context, url string) (string, error)
}
