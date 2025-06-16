package main

import (
	"context"
	"log"
	"os"

	"snp_scrapper/internal/api"
	"snp_scrapper/internal/config"
	"snp_scrapper/internal/service"
	"snp_scrapper/internal/store"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/sashabaranov/go-openai"
)

func main() {
	// Load configuration
	cfg := config.New()

	// Initialize OpenAI client
	openaiClient := openai.NewClient(cfg.OpenAIAPIKey)

	// Initialize AWS store
	awsStore, err := store.NewAWSStore(
		cfg.AWSConfig.Region,
		cfg.AWSConfig.S3Bucket,
		cfg.AWSConfig.SNSTopicARN,
	)
	if err != nil {
		log.Fatal("Failed to initialize AWS store:", err)
	}

	// Initialize services
	stockService := service.NewStockService(awsStore, openaiClient)

	// Initialize handlers
	handler := api.NewHandler(stockService)

	// Initialize router
	r := gin.Default()

	// Register routes
	r.GET("/api/sp500", handler.GetSP500Stocks)
	r.GET("/api/qualitative", handler.GetQualitativeStocks)
	r.POST("/api/subscribe", handler.Subscribe)

	// Initialize cron job
	c := cron.New()
	_, err = c.AddFunc("0 0 * * *", func() {
		ctx := context.Background()
		if err := stockService.UpdateStockList(ctx); err != nil {
			log.Printf("Failed to update stock list: %v", err)
		}
	})
	if err != nil {
		log.Fatal("Failed to schedule cron job:", err)
	}
	c.Start()

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
} 