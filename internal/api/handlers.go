package api

import (
	"net/http"

	"snp_scrapper/internal/service"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests
type Handler struct {
	stockService *service.StockService
}

// NewHandler creates a new handler
func NewHandler(stockService *service.StockService) *Handler {
	return &Handler{
		stockService: stockService,
	}
}

// GetSP500Stocks handles GET /api/sp500
func (h *Handler) GetSP500Stocks(c *gin.Context) {
	stockList, err := h.stockService.GetStockList(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get S&P 500 data"})
		return
	}
	c.JSON(http.StatusOK, stockList)
}

// GetQualitativeStocks handles GET /api/qualitative
func (h *Handler) GetQualitativeStocks(c *gin.Context) {
	stocks, err := h.stockService.GetQualitativeStocks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get qualitative stocks"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"stocks": stocks})
}

// Subscribe handles POST /api/subscribe
func (h *Handler) Subscribe(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.stockService.Subscribe(c.Request.Context(), req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to subscribe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully subscribed"})
} 