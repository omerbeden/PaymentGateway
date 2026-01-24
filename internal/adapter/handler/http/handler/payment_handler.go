package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/omerbeden/paymentgateway/internal/usecase/payment"
)

type PaymentHandler struct {
	createPaymentUC *payment.CreatePaymentUseCase
}

func NewPaymentHandler(createPaymentUC *payment.CreatePaymentUseCase) *PaymentHandler {
	return &PaymentHandler{
		createPaymentUC: createPaymentUC,
	}
}

type CreatePaymentRequest struct {
	Amount     float64           `json:"amount" binding:"required,gt=0"`
	Currency   string            `json:"currency" binding:"required,oneof=USD EUR TRY GBP"`
	ProviderID string            `json:"provider_id" binding:"required"`
	Metadata   map[string]string `json:"metadata"`
}

type CreatePaymentResponse struct {
	ID        string    `json:"payment_id"`
	Status    string    `json:"status"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {

	var req CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.createPaymentUC.Execute(c.Request.Context(), payment.CreatePaymentInput{
		Amount:   req.Amount,
		Currency: req.Currency,
		Metadata: req.Metadata,
	})
	if err != nil {
		//h.log.Error("Failed to create payment", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment"})
		return
	}
	resp := CreatePaymentResponse{
		ID:        payment.ID,
		Status:    string(payment.Status),
		Amount:    payment.Amount,
		Currency:  string(payment.Currency),
		CreatedAt: payment.CreatedAt,
	}
	c.JSON(http.StatusCreated, resp)
}
