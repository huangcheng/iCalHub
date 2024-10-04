package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/huangcheng/icalhub/config"

	"io"
	"net/http"
)
import "github.com/huangcheng/icalhub/utils"
import "github.com/huangcheng/icalhub/spiders"

type HolidaysController struct{}

type MoviesController struct{}

var response = new(utils.Response)

func (controller HolidaysController) China(c *gin.Context) {
	resp, err := http.Get("https://calendars.icloud.com/holidays/cn_zh.ics")

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	response.Calendar(c, body)
}

func (controller MoviesController) Douban(c *gin.Context) {
	var spider = new(spiders.Douban)

	conf, exists := c.Get("config")

	if !exists {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	content, err := spider.Run(conf.(config.Config).UserAgent)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	response.Calendar(c, []byte(content))
}
