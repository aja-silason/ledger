package http

import "github.com/gin-gonic/gin"

func SetupRoutes(c *Controllers) *gin.Engine {
	r := gin.Default()

	accounts := r.Group("/accounts")
	{
		accounts.POST("", c.AccountController.CreateAccount)
	}

	return r
}
