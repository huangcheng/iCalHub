package middlewares

import "github.com/gin-gonic/gin"

import "github.com/huangcheng/icalhub/config"

func ConfigMiddleware(config config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("config", config)

		c.Next()
	}
}
