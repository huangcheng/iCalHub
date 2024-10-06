package main

import "github.com/gin-gonic/gin"

import (
	"github.com/huangcheng/icalhub/config"
	"github.com/huangcheng/icalhub/middlewares"
)

func main() {
	config := config.ReadConfig()

	r := gin.Default()

	r.Use(middlewares.ConfigMiddleware(config))
	r.Use(middlewares.CacheMiddleware(config))

	setupRoutes(r)

	r.Run(":" + config.Port)
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
