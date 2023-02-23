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

func main() {
	initializeDataStorage()

	TokenSecret = []byte(os.Getenv("TOKEN_SECRET"))

	router := gin.Default()
	router.POST("/login", login)
	router.GET("/users", getUsers)

	router.GET("/files", getAllFiles)
	router.POST("/upload", uploadFile)

	err := router.Run("localhost:8080")
	if err != nil {
		log.Fatal("Failed to start router!")
	}
}
