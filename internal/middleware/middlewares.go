package middleware

import (
	"weather4you/config"
	"weather4you/pkg/logger"
)

// Middleware manager
type MiddlewareManager struct {
	cfg    *config.Config
	logger logger.Logger
}

// Middleware manager constructor
func NewMiddlewareManager(cfg *config.Config, logger logger.Logger) *MiddlewareManager {
	return &MiddlewareManager{cfg: cfg, logger: logger}
}
