package handler

import (
	"logger"
	"net/http"
	"utils"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Refresh(ctx *gin.Context) {
	ctx.Request.ParseForm()

	user := ctx.PostForm("user")
	password := ctx.PostForm("password")

	if user == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		logger.Logger.Println("Invalid request: user is required")
		ctx.Abort()
		return
	}
	if password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		logger.Logger.Println("Invalid request: password is required")
		ctx.Abort()
		return
	}

	if user != viper.GetString("user") || password != viper.GetString("password") {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		logger.Logger.Printf("Unauthorized request: user=%s, password=%s", user, password)
		ctx.Abort()
		return
	}

	utils.RefreshSum256Map()
	utils.RefreshSum256List()

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Refresh successful",
	})
	logger.Logger.Println("Refresh successful")
	ctx.Next()
}
