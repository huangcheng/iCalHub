package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"net/http"
)

type Response struct{}

func (r Response) Calendar(context *gin.Context, data []byte) {
	context.Render(http.StatusOK, render.Data{
		ContentType: "text/calendar; charset=utf-8",
		Data:        data,
	})
}
