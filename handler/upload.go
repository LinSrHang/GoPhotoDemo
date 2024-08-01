package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func UploadPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "upload.tpl", gin.H{
		"serverUrl": viper.GetString("serverUrl"),
	})
}
