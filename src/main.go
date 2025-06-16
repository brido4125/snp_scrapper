package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/sashabaranov/go-openai"
	"github.com/go-redis/redis/v8"
)

var (
	openaiClient *openai.Client
	redisClient  *redis.Client
)

func init() {
	// OpenAI 클라이언트 초기화
	openaiClient = openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	// Redis 클라이언트 초기화
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func main() {
	// Gin 라우터 초기화
	r := gin.Default()

	// API 라우트 설정
	r.GET("/api/sp500", getSP500Stocks)
	r.GET("/api/qualitative", getQualitativeStocks)
	r.POST("/api/subscribe", subscribeToUpdates)

	// 크론 작업 설정
	c := cron.New()
	_, err := c.AddFunc("0 0 * * *", updateSP500Stocks) // 매일 자정에 실행
	if err != nil {
		log.Fatal("Failed to schedule cron job:", err)
	}
	c.Start()

	// 서버 시작
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// S&P 500 종목 업데이트 함수
func updateSP500Stocks() {
	ctx := context.Background()
	
	// ChatGPT API를 통해 S&P 500 종목 정보 가져오기
	prompt := "List all current S&P 500 companies with their ticker symbols and market cap as of today."
	resp, err := openaiClient.CreateChatCompletion(
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
		log.Printf("Failed to get S&P 500 data: %v", err)
		return
	}

	// 이전 데이터와 비교
	prevData, err := redisClient.Get(ctx, "sp500_data").Result()
	if err != nil && err != redis.Nil {
		log.Printf("Failed to get previous data: %v", err)
		return
	}

	// 새로운 데이터 저장
	err = redisClient.Set(ctx, "sp500_data", resp.Choices[0].Message.Content, 24*time.Hour).Err()
	if err != nil {
		log.Printf("Failed to save new data: %v", err)
		return
	}

	// 변경사항이 있다면 구독자들에게 알림
	if prevData != resp.Choices[0].Message.Content {
		notifySubscribers(ctx, resp.Choices[0].Message.Content)
	}
}

// 구독자들에게 알림을 보내는 함수
func notifySubscribers(ctx context.Context, newData string) {
	subscribers, err := redisClient.SMembers(ctx, "subscribers").Result()
	if err != nil {
		log.Printf("Failed to get subscribers: %v", err)
		return
	}

	// TODO: 실제 알림 전송 로직 구현 (이메일, 웹훅 등)
	for _, subscriber := range subscribers {
		log.Printf("Notifying subscriber: %s", subscriber)
	}
}

// S&P 500 종목 목록을 반환하는 API
func getSP500Stocks(c *gin.Context) {
	ctx := context.Background()
	data, err := redisClient.Get(ctx, "sp500_data").Result()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get S&P 500 data"})
		return
	}
	c.JSON(200, gin.H{"data": data})
}

// 정량적 기준을 만족하는 종목을 반환하는 API
func getQualitativeStocks(c *gin.Context) {
	// TODO: 정량적 기준에 따른 필터링 로직 구현
	c.JSON(200, gin.H{"message": "Not implemented yet"})
}

// 구독 API
func subscribeToUpdates(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	ctx := context.Background()
	err := redisClient.SAdd(ctx, "subscribers", req.Email).Err()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to subscribe"})
		return
	}

	c.JSON(200, gin.H{"message": "Successfully subscribed"})
}
