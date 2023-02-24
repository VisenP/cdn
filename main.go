package main

import (
	"cdn/database"
	"cdn/global"
	"cdn/routes/file"
	"cdn/routes/user"
	"github.com/gin-gonic/gin"
	"log"
)

func init() {
	global.InitGlobals()
}

func main() {
	database.InitializeDataStorage()

	router := gin.Default()
	router.POST("/login", user.Login)
	router.GET("/user", user.GetUsers)

	router.GET("/file", file.GetAllFiles)
	router.GET("/file/:id", file.GetFile(false))
	router.GET("/view/:id", file.GetFile(true))
	router.POST("/upload", file.UploadFile)

	err := router.Run(":" + global.Port)
	if err != nil {
		log.Fatal("Failed to start router!")
	}
}
