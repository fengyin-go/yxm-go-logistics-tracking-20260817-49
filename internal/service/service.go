package service

import (
	"logistics/internal/config"
	"logistics/internal/store"
	"logistics/pkg/logger"
)

type Service struct {
	store store.Store
	log   *logger.Logger
	cfg   *config.Config
}

func New(st store.Store, log *logger.Logger, cfg *config.Config) *Service {
	return &Service{store: st, log: log, cfg: cfg}
}

func normalizePagination(page, size, fallbackSize, maxSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if fallbackSize < 1 {
		fallbackSize = 20
	}
	if size < 1 {
		size = fallbackSize
	}
	if maxSize > 0 && size > maxSize {
		size = maxSize
	}
	return page, size
}
