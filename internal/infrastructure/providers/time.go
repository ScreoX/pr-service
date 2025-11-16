package providers

import (
	"time"

	"pr-service/internal/app"
)

type currentTimeProvider struct{}

func NewCurrentTime() app.TimeProvider {
	return &currentTimeProvider{}
}

func (r *currentTimeProvider) Now() time.Time {
	return time.Now().UTC()
}
