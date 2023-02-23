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

type userLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func login(c *gin.Context) {

	var userLogin userLogin

	err := c.BindJSON(&userLogin)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := findFirst(&userData, func(u user) bool {
		return u.Username == userLogin.Username
	})

	if user == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if !checkPasswordHash(userLogin.Password, (*user).Password) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.IndentedJSON(http.StatusOK, generateJWT((*user).Id))
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

var TokenSecret []byte

func main() {
	initializeDataStorage()

	TokenSecret = []byte(os.Getenv("TOKEN_SECRET"))

	router := gin.Default()
	router.POST("/login", login)

	router.Use(validateUser)
	router.GET("/users", getUsers)

	err := router.Run("localhost:8080")
	if err != nil {
		log.Fatal("Failed to start router!")
	}
}
