package store

import (
	"context"
	"snp_scrapper/internal/models"
)

// Store defines the interface for storage operations
type Store interface {
	SaveStockList(ctx context.Context, stockList *models.StockList) error
	GetStockList(ctx context.Context) (*models.StockList, error)
	AddSubscriber(ctx context.Context, email string) error
	GetSubscribers(ctx context.Context) ([]string, error)
} 