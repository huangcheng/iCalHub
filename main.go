package main

import "github.com/gin-gonic/gin"

import (
	"github.com/huangcheng/icalhub/config"
	"github.com/huangcheng/icalhub/controllers"
	"github.com/huangcheng/icalhub/middlewares"
)

func main() {
	config := config.ReadConfig()

	r := gin.Default()

	r.Use(middlewares.ConfigMiddleware(config))
	r.Use(middlewares.CacheMiddleware(config))

	root := r.Group("/")
	{

		holidays := root.Group("/holidays")
		{
			controller := new(controllers.HolidaysController)

			holidays.GET("/china", controller.China)
		}

		movies := root.Group("/movies")
		{
			controller := new(controllers.MoviesController)

			movies.GET("/douban", controller.Douban)
		}

		root.GET("/", func(c *gin.Context) {
			html := `
				<!DOCTYPE html>
				<html lang="en">
				<meta charset="UTF-8">
					<head>
						<title>iCalHub</title>
					</head>
					<body>
						<h1 align="center">iCalHub</h1>

						<details open>
							<summary>Holidays</summary>
							
							<ul>
								<li><a href="/holidays/china">China Public Holidays</a></li>
							</ul>
						</details>

						<details open>
							<summary>Movies</summary>
	
							<ul>
								<li><a href="/movies/douban">Douban Coming Movies</a></li>
							</ul>
						</details>
					</body>
				</html>
			`

			c.Data(200, "text/html; charset=utf-8", []byte(html))
		})
	}

	r.Run(":" + config.Port)
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
