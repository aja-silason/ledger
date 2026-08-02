package http

import "github.com/gin-gonic/gin"

func SetupRoutes(c *Controllers) *gin.Engine {
	r := gin.Default()

	accounts := r.Group("/accounts")
	{
		accounts.POST("", c.AccountController.CreateAccount)
		accounts.GET("", c.AccountController.FindAll)
		accounts.GET("/:id", c.AccountController.FindByID)
		accounts.POST("/deposit", c.AccountController.DepositInAccount)
		accounts.POST("/transfer", c.AccountController.Transfer)
		accounts.POST("/withdraw", c.AccountController.WithdrawCardLess)
		accounts.PATCH("/withdraw/cancel", c.AccountController.CancelWithdrawCardLess)
		accounts.PATCH("/withdraw/drawed", c.AccountController.DrawedCardLess)
	}

	return r
}
