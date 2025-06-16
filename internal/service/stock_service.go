package service

import (
	"context"
	"time"

	"snp_scrapper/internal/models"
	"snp_scrapper/internal/store"

	"github.com/sashabaranov/go-openai"
)

// StockService handles business logic for stocks
type StockService struct {
	store       store.Store
	openaiClient *openai.Client
}

// Store defines the interface for storage operations
type Store interface {
	SaveStockList(ctx context.Context, stockList *models.StockList) error
	GetStockList(ctx context.Context) (*models.StockList, error)
	AddSubscriber(ctx context.Context, email string) error
	GetSubscribers(ctx context.Context) ([]string, error)
}

// NewStockService creates a new stock service
func NewStockService(store Store, openaiClient *openai.Client) *StockService {
	return &StockService{
		store:       store,
		openaiClient: openaiClient,
	}
}

// UpdateStockList updates the S&P 500 stock list
func (s *StockService) UpdateStockList(ctx context.Context) error {
	prompt := "List all current S&P 500 companies with their ticker symbols and market cap as of today."
	resp, err := s.openaiClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		return err
	}

	// Parse the response and create a StockList
	// Note: This is a simplified version. In reality, you'd need to parse the response
	// and convert it to proper Stock objects
	stockList := &models.StockList{
		Stocks: []models.Stock{}, // TODO: Parse response into proper Stock objects
		Date:   time.Now().Format("2006-01-02"),
	}

	// Save to storage
	return s.store.SaveStockList(ctx, stockList)
}

// GetStockList retrieves the current stock list
func (s *StockService) GetStockList(ctx context.Context) (*models.StockList, error) {
	return s.store.GetStockList(ctx)
}

// GetQualitativeStocks returns stocks that meet certain criteria
func (s *StockService) GetQualitativeStocks(ctx context.Context) ([]models.Stock, error) {
	stockList, err := s.store.GetStockList(ctx)
	if err != nil {
		return nil, err
	}

	// TODO: Implement filtering logic based on qualitative criteria
	return stockList.Stocks, nil
}

// Subscribe adds a new subscriber
func (s *StockService) Subscribe(ctx context.Context, email string) error {
	return s.store.AddSubscriber(ctx, email)
} 