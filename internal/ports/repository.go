package ports

import (
	"Shortxn/internal/domain"
	"time"
)

type URLRepository interface {
	Store(url *domain.URL) error
	GetByID(id string) (*domain.URL, error)
	GetByLongURL(longURL string) (*domain.URL, error)
	IncrementClicks(id string) error
}

type CacheRepository interface {
	Set(key string, value string, expiration time.Duration) error
	Get(key string) (string, error)
}
