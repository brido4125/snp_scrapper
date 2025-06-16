package store

import (
	"context"
	"encoding/json"

	"snp_scrapper/internal/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// AWSStore implements the storage interface using AWS services
type AWSStore struct {
	s3Client  *s3.Client
	snsClient *sns.Client
	bucket    string
	topicARN  string
}

// NewAWSStore creates a new AWS store
func NewAWSStore(region, bucket, topicARN string) (*AWSStore, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}

	return &AWSStore{
		s3Client:  s3.NewFromConfig(cfg),
		snsClient: sns.NewFromConfig(cfg),
		bucket:    bucket,
		topicARN:  topicARN,
	}, nil
}

// SaveStockList saves the stock list to S3
func (s *AWSStore) SaveStockList(ctx context.Context, stockList *models.StockList) error {
	data, err := json.Marshal(stockList)
	if err != nil {
		return err
	}

	// Save to S3
	_, err = s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String("stocks.json"),
		Body:        aws.NewReadSeekCloser(aws.NewReader(data)),
		ContentType: aws.String("application/json"),
	})
	if err != nil {
		return err
	}

	// Publish to SNS if there are changes
	message, err := json.Marshal(map[string]interface{}{
		"type": "stock_update",
		"data": stockList,
	})
	if err != nil {
		return err
	}

	_, err = s.snsClient.Publish(ctx, &sns.PublishInput{
		TopicArn: aws.String(s.topicARN),
		Message:  aws.String(string(message)),
	})

	return err
}

// GetStockList retrieves the stock list from S3
func (s *AWSStore) GetStockList(ctx context.Context) (*models.StockList, error) {
	result, err := s.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String("stocks.json"),
	})
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	var stockList models.StockList
	if err := json.NewDecoder(result.Body).Decode(&stockList); err != nil {
		return nil, err
	}

	return &stockList, nil
}

// AddSubscriber adds a new subscriber to SNS topic
func (s *AWSStore) AddSubscriber(ctx context.Context, email string) error {
	_, err := s.snsClient.Subscribe(ctx, &sns.SubscribeInput{
		TopicArn: aws.String(s.topicARN),
		Protocol: aws.String("email"),
		Endpoint: aws.String(email),
	})
	return err
}

// GetSubscribers retrieves all subscribers from SNS topic
func (s *AWSStore) GetSubscribers(ctx context.Context) ([]string, error) {
	result, err := s.snsClient.ListSubscriptionsByTopic(ctx, &sns.ListSubscriptionsByTopicInput{
		TopicArn: aws.String(s.topicARN),
	})
	if err != nil {
		return nil, err
	}

	subscribers := make([]string, 0, len(result.Subscriptions))
	for _, sub := range result.Subscriptions {
		if sub.Protocol == "email" {
			subscribers = append(subscribers, *sub.Endpoint)
		}
	}

	return subscribers, nil
} 