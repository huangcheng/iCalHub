package controllers

import (
	"github.com/huangcheng/icalhub/config"
	"net/http"
)

import "github.com/gin-gonic/gin"

import (
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
