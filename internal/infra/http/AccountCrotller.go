package http

import (
	"net/http"

	"github.com/aja-silason/ledger/internal/application"
	"github.com/aja-silason/ledger/internal/domain"
	"github.com/gin-gonic/gin"
)

type AccountController struct {
	Create     *application.CreateAccountService
	GetAccount *application.GetAccountServiceFinder
}

func NewAccountController(
	Create *application.CreateAccountService,
	GetAccount *application.GetAccountServiceFinder) *AccountController {
	return &AccountController{
		Create:     Create,
		GetAccount: GetAccount}
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

	resul, err := h.GetAccount.FindAccountByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, resul)
}

func (h *AccountController) FindAll(c *gin.Context) {

	result, err := h.GetAccount.FindAllAccount(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, result)
}
