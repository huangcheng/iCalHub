package main

import "github.com/gin-gonic/gin"

import "github.com/huangcheng/icalhub/controllers"

func holidays(g *gin.RouterGroup) {
	holidays := g.Group("/holidays")

	controller := new(controllers.HolidaysController)

	holidays.GET("/china", controller.China)
}

func movies(g *gin.RouterGroup) {
	movies := g.Group("/movies")

	controller := new(controllers.MoviesController)

	movies.GET("/douban", controller.Douban)
	movies.GET("/imdb/*region", controller.IMDb)
}

func astronomy(g *gin.RouterGroup) {
	astronomy := g.Group("/astronomy")

	controller := new(controllers.AstronomyController)

	astronomy.GET("/moon/*year", controller.Moon)
}

func games(g *gin.RouterGroup) {
	games := g.Group("/games")

	controller := new(controllers.GamesController)

	games.GET("/steam/:type/*language", controller.Steam)
}

func index(g *gin.RouterGroup) {
	html := `
		<!DOCTYPE html>
		<html lang="en">
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
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
						<li><a href="/movies/imdb">Upcoming releases - IMDb</a></li>
						<li><a href="/movies/douban">Upcoming Movies - Douban</a></li>
					</ul>
				</details>

				<details open>
					<summary>Astronomy</summary>

					<ul>
						<li><a href="/astronomy/moon">Date and Time of the Moon Phase｜Hong Kong Observatory(HKO)</a></li>
					</ul>
				</details>

				<details open>
					<summary>Games</summary>

					<ul>
						<li><a href="/games/steam/popular/">Upcoming Releases - Steam</a></li>
					</ul>
				</details>
			</body>
		</html>
		`

	g.GET("/", func(c *gin.Context) {

		c.Data(200, "text/html; charset=utf-8", []byte(html))
	})
}

func setupRoutes(r *gin.Engine) {
	root := r.Group("/")

	movies(root)

	holidays(root)

	astronomy(root)

	games(root)

	index(root)
}
