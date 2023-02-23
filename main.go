package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Could not load .env!")
	}

	log.Println("Env file loaded!")
}

var TokenSecret []byte

var port string

func main() {
	initializeDataStorage()

	TokenSecret = []byte(os.Getenv("TOKEN_SECRET"))
	port = os.Getenv("PORT")

	router := gin.Default()
	router.POST("/login", login)
	router.GET("/users", getUsers)

	router.GET("/files", getAllFiles)
	router.GET("/files/:id", getFile(false))
	router.GET("/view/:id", getFile(true))
	router.POST("/upload", uploadFile)

	err := router.Run(":" + port)
	if err != nil {
		log.Fatal("Failed to start router!")
	}
}
