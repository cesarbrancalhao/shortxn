package service

import (
	"Shortxn/internal/domain"
	"Shortxn/internal/ports"
	"crypto/sha256"
	"encoding/base64"
	"time"
)

type URLService struct {
	urlRepo   ports.URLRepository
	cacheRepo ports.CacheRepository
}

func NewURLService(urlRepo ports.URLRepository, cacheRepo ports.CacheRepository) *URLService {
	return &URLService{
		urlRepo:   urlRepo,
		cacheRepo: cacheRepo,
	}
}

func (s *URLService) ShortenURL(longURL string) (*domain.URL, error) {
	if existing, err := s.urlRepo.GetByLongURL(longURL); err == nil {
		return existing, nil
	}

	hash := sha256.Sum256([]byte(longURL))
	shortID := base64.URLEncoding.EncodeToString(hash[:8])

	url := &domain.URL{
		ID:        shortID,
		LongURL:   longURL,
		ShortURL:  shortID,
		CreatedAt: time.Now(),
		Clicks:    0,
	}

	if err := s.urlRepo.Store(url); err != nil {
		return nil, err
	}

	s.cacheRepo.Set(url.ShortURL, url.LongURL, 24*time.Hour)

	return url, nil
}
