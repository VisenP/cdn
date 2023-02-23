package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
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
