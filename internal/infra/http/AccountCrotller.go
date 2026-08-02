package http

import (
	"net/http"

	"github.com/aja-silason/ledger/internal/application"
	"github.com/aja-silason/ledger/internal/domain"
	"github.com/gin-gonic/gin"
)

type AccountController struct {
	Create         *application.CreateAccountService
	Finder         *application.GetAccountServiceFinder
	DepositIn      *application.DepositIn
	TransferMoney  *application.TransferMoney
	Withdraw       *application.WidrawCardLess
	CancelWithdraw *application.CancelWithDrawCardLess
	Drawed         *application.DrawedCardLess
}

func NewAccountController(
	Create *application.CreateAccountService,
	Finder *application.GetAccountServiceFinder,
	DepositIn *application.DepositIn,
	TransferMoney *application.TransferMoney,
	Withdraw *application.WidrawCardLess,
	CancelWithdraw *application.CancelWithDrawCardLess,
	Drawed *application.DrawedCardLess) *AccountController {
	return &AccountController{
		Create:         Create,
		Finder:         Finder,
		DepositIn:      DepositIn,
		TransferMoney:  TransferMoney,
		Withdraw:       Withdraw,
		CancelWithdraw: CancelWithdraw,
		Drawed:         Drawed}
}

func (h *AccountController) CreateAccount(c *gin.Context) {
	var req domain.Account
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload inválido"})
		return
	}

	result, err := h.Create.PostCreateAccount(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *AccountController) FindByID(c *gin.Context) {
	id := c.Param("id")

	resul, err := h.Finder.FindAccountByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, resul)
}

func (h *AccountController) FindAll(c *gin.Context) {

	result, err := h.Finder.FindAllAccount(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, result)
}

func (h *AccountController) DepositInAccount(c *gin.Context) {

	idempotencyKey := c.GetHeader("X-Idempotency-Key")
	if idempotencyKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Header X-Idempotency-Key é obrigatória"})
		return
	}

	var req application.DepositInput

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload inválido"})
		return
	}

	result, err := h.DepositIn.Deposit(c.Request.Context(), &req, idempotencyKey)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *AccountController) Transfer(c *gin.Context) {

	idempotencyKey := c.GetHeader("X-Idempotency-Key")
	if idempotencyKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Header X-Idempotency-Key é obrigatória"})
		return
	}

	var req application.TransferMoneyInput

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload inválido"})
		return
	}

	result, err := h.TransferMoney.Transfer(c.Request.Context(), &req, idempotencyKey)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *AccountController) WithdrawCardLess(c *gin.Context) {

	var req application.WithDrawInput

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload inválido"})
		return
	}

	result, err := h.Withdraw.DemandWidraw(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *AccountController) CancelWithdrawCardLess(c *gin.Context) {

	var req application.CancelWithDrawCardLessInput

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload inválido"})
		return
	}

	result, err := h.CancelWithdraw.Cancel(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *AccountController) DrawedCardLess(c *gin.Context) {

	idempotencyKey := c.GetHeader("X-Idempotency-Key")
	if idempotencyKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Header X-Idempotency-Key é obrigatória"})
		return
	}

	var req application.DrawedCardLessInput

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload inválido"})
		return
	}

	result, err := h.Drawed.Drawed(c.Request.Context(), idempotencyKey, &req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}
