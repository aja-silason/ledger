package http

import (
	"net/http"

	"github.com/aja-silason/ledger/internal/application"
	"github.com/aja-silason/ledger/internal/domain"
	"github.com/gin-gonic/gin"
)

type LedgerController struct {
	service       *application.LedgerService
	createAccount *application.CreateAccountService
}

func NewLedgerController(
	service *application.LedgerService,
	createAccount *application.CreateAccountService) *LedgerController {
	return &LedgerController{
		service:       service,
		createAccount: createAccount,
	}
}

func (h *LedgerController) CreateTransaction(c *gin.Context) {
	var req domain.Transaction

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload inválido"})
		return
	}

	idempotencyKey := c.GetHeader("X-Idempontency-Key")
	if idempotencyKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Header X-Idempotency-Key é obrigatória"})
		return
	}

	req.IdempotencyKey = idempotencyKey

	result, err := h.service.PostTransaction(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)

}

func (h *LedgerController) CreateAccountTransaction(c *gin.Context) {
	var req domain.Account
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload inválido"})
		return
	}

	result, err := h.createAccount.PostCreateAccount(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}
