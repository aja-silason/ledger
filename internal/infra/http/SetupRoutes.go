package http

import "github.com/gin-gonic/gin"

func SetupRoutes(c *Controllers) *gin.Engine {
	r := gin.Default()

	accounts := r.Group("/accounts")
	{
		accounts.POST("", c.AccountController.CreateAccount)
		accounts.GET("/:id", c.AccountController.FindByID)
	}

	return r
}
