package main

import (
	"github.com/gin-gonic/gin"
	"github.com/huangcheng/icalhub/config"
	"github.com/huangcheng/icalhub/controllers"
	"github.com/huangcheng/icalhub/middlewares"
)

func main() {
	config := config.ReadConfig()

	r := gin.Default()

	r.Use(middlewares.ConfigMiddleware(config))
	r.Use(middlewares.CacheMiddleware(config))

	calendar := r.Group("/calendar")
	{
		holidays := new(controllers.HolidaysController)

		calendar.GET("/holidays/china", holidays.China)

		movies := new(controllers.MoviesController)
		calendar.GET("/movies/douban", movies.Douban)
	}

	r.Run(":" + config.Port)
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
