package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func IndexPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.tpl", gin.H{
		"serverUrl": viper.GetString("serverUrl"),
	})
}
