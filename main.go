package main

import (
	"config"
	"fmt"
	"handler"
	"logger"
	"utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func HandlerStart() {
	gin.SetMode(viper.GetString("mode"))

	r := gin.Default()
	r.LoadHTMLGlob("template/*")
	Cors := cors.DefaultConfig()

	allowOrigins := []string{}
	port := viper.GetString("port")
	for _, origin := range viper.GetStringSlice("allowOrigins") {
		allowOrigins = append(allowOrigins, origin+port)
	}

	Cors.AllowOrigins = allowOrigins

	api := r.Group("/api")
	{
		api.Use(cors.New(Cors))
		v1 := api.Group("/v1")
		{
			static := v1.Group("/static")
			{
				verify := static.Group("")
				{
					verify.Use(handler.GetFidList)
					verify.Use(handler.VerifyFidList)
					verify.DELETE("/delete", handler.RemoveImages)
					verify.GET("/:controller/get", handler.GetImages)
				}
				static.POST("/post", handler.AddImages)
				static.GET("/fidList", handler.GetFidListPagingQuery)
				static.POST("/postZip", handler.AddZippedImages)
			}
		}
	}
	page := r.Group("/page")
	{
		page.GET("/index", handler.IndexPage)
		page.GET("/upload", handler.UploadPage)
	}

	fmt.Println("Running in \"" + viper.GetString("mode") + "\" mode.")
	r.Run(viper.GetString("port"))
}

func main() {
	logger.LoggerInit()
	config.ConfigInit()
	utils.UtilsInit()
	HandlerStart()
}
