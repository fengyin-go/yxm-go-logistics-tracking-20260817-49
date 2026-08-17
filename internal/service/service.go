package service

import (
	"strings"

	"logistics/internal/config"
	"logistics/internal/model"
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

func normalizeWaybillFilter(filter model.WaybillFilter) model.WaybillFilter {
	filter.Status = strings.TrimSpace(filter.Status)
	filter.Keyword = strings.TrimSpace(filter.Keyword)
	return filter
}
