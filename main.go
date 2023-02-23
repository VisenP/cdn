package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

type returnUser struct {
	Username string `json:"username"`
}

func getUsers(c *gin.Context) {
	users := make([]returnUser, 0)

	for _, user := range userData {
		users = append(users, returnUser{Username: user.Username})
	}

	c.IndentedJSON(http.StatusOK, users)
}

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Could not load .env!")
	}

	log.Println("Env file loaded!")
}

var TOKEN_SECRET string

func main() {
	initializeDataStorage()

	TOKEN_SECRET = os.Getenv("TOKEN_SECRET")

	router := gin.Default()

	router.GET("/users", getUsers)

	err := router.Run("localhost:8080")
	if err != nil {
		log.Fatal("Failed to start router!")
	}
}
