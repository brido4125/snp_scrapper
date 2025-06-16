package models

// Stock represents a company in S&P 500
type Stock struct {
	Ticker    string  `json:"ticker"`
	Name      string  `json:"name"`
	MarketCap float64 `json:"market_cap"`
}

// StockList represents the list of S&P 500 stocks
type StockList struct {
	Stocks []Stock `json:"stocks"`
	Date   string  `json:"date"`
}

// Subscriber represents a user who subscribed to updates
type Subscriber struct {
	Email string `json:"email"`
} 