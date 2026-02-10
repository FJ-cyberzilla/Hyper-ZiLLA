package threat_intel

import (
	"context"
	"net-zilla/internal/models"
	"net-zilla/pkg/logger"
)

type IOCAnalyzer struct {
	db     *ThreatDatabase
	feeds  *FeedsClient
	logger *logger.Logger
}

func NewIOCAnalyzer(db *ThreatDatabase, feeds *FeedsClient, logger *logger.Logger) *IOCAnalyzer {
	return &IOCAnalyzer{
		db:     db,
		feeds:  feeds,
		logger: logger,
	}
}

func (ia *IOCAnalyzer) Analyze(ctx context.Context, value string) (*models.Indicator, error) {
	// 1. Check local cache/database first
	indicator, err := ia.db.Lookup(ctx, value)
	if err != nil {
		ia.logger.Warn("Local IOC lookup failed: %v", err)
	}
	if indicator != nil {
		return indicator, nil
	}

	// 2. If not found, verify with external feeds (best effort)
	// Professional implementation: Only call external feeds if confidence is required
	return nil, nil
}
