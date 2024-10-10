package controllers

import (
	"net/http"
	"strings"
)

import "github.com/gin-gonic/gin"

import (
	"github.com/huangcheng/icalhub/config"
	"github.com/huangcheng/icalhub/handlers"
)

type GamesController struct {
}

func (controller GamesController) Steam(c *gin.Context) {
	t := c.Param("type")

	l := c.Param("language")

	l = strings.ReplaceAll(l, "/", "")

	if len(l) == 0 {
		l = "zh_CN"
	}

	var handler = new(handlers.Steam)

	conf, exists := c.Get("config")

	if !exists {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	handler.UserAgent = conf.(config.Config).UserAgent

	handler.Type = t
	handler.Language = l

	content, err := handler.Run()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	response.Calendar(c, []byte(content))
}
