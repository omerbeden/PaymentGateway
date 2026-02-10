package handler

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/omerbeden/paymentgateway/internal/adapter/provider"
	"github.com/omerbeden/paymentgateway/internal/usecase/webhook"
)

type WebhookHandler struct {
	weebhookUseCase *webhook.ProcessWebHookUseCase
}

func NewWebhookHandler(weebhookUseCase *webhook.ProcessWebHookUseCase) *WebhookHandler {
	return &WebhookHandler{
		weebhookUseCase: weebhookUseCase,
	}
}

func (h *WebhookHandler) HandlePaypal(c *gin.Context) {
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	webhookCtx := &provider.WebhookContext{
		Payload:   payload,
		Headers:   c.Request.Header,
		Signature: c.GetHeader("PAYPAL-TRANSMISSION-SIG"),
	}
	input := webhook.ProcessWebHookInput{
		ProviderId:     "paypal",
		WebhookContext: webhookCtx,
	}

	if err := h.weebhookUseCase.Execute(c.Request.Context(), input); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})

}
