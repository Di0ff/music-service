package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"music-service/internal/api"
	"music-service/internal/db"
	"music-service/internal/logging"

	_ "music-service/docs"
)

func main() {
	logging.Init()
	logging.Logger.Info("Launching the app")

	db.Init()

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api.Routes(r)

	err := r.Run(":8080")
	if err != nil {
		logging.Logger.Fatal("Server startup error: ", err)
	}
}
