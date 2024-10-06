package controllers

import (
	"io"
	"net/http"
)

import "github.com/gin-gonic/gin"

import "github.com/huangcheng/icalhub/utils"

type HolidaysController struct {
}

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
