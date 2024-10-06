package controllers

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

import (
	"github.com/gin-gonic/gin"
)

import (
	"github.com/huangcheng/icalhub/config"
	"github.com/huangcheng/icalhub/handlers"
)

type AstronomyController struct {
}

func (controller AstronomyController) Moon(c *gin.Context) {
	t := time.Now()
	y := t.Format("2006")

	year := c.Param("year")

	if matched, err := regexp.Match(`\d{4}`, []byte(year)); err != nil || !matched {
		year = y
	}

	year = strings.ReplaceAll(year, "/", "")

	url := fmt.Sprintf("https://www.hko.gov.hk/tc/gts/astronomy/files/MoonPhases_%s.xml", year)

	conf, exists := c.Get("config")

	if !exists {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	handler := new(handlers.HKO)

	handler.UserAgent = conf.(config.Config).UserAgent
	handler.URL = url

	content, err := handler.Run()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	response.Calendar(c, []byte(content))
}
