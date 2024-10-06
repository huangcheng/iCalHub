package controllers

import (
	"fmt"
	"net/http"
	"strings"
)

import "github.com/gin-gonic/gin"

import (
	"github.com/huangcheng/icalhub/config"
	"github.com/huangcheng/icalhub/handlers"
)

type MoviesController struct{}

func (controller MoviesController) Douban(c *gin.Context) {
	var handler = new(handlers.Douban)

	conf, exists := c.Get("config")

	if !exists {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	handler.UserAgent = conf.(config.Config).UserAgent

	content, err := handler.Run()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	response.Calendar(c, []byte(content))
}

func (controller MoviesController) IMDb(c *gin.Context) {
	region := c.Param("region")
	region = strings.ReplaceAll(region, "/", "")

	if len(region) == 0 {
		region = "CN"
	}

	url := fmt.Sprintf("https://www.imdb.com/calendar/?region=%s&type=MOVIE", region)

	conf, exists := c.Get("config")

	if !exists {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	handler := new(handlers.IMDb)

	handler.UserAgent = conf.(config.Config).UserAgent
	handler.URL = url

	content, err := handler.Run()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	response.Calendar(c, []byte(content))
}
