package user

import (
	"cdn/auth"
	"cdn/database"
	"cdn/utils"
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

func Login(c *gin.Context) {

	var userLogin userLogin

	err := c.BindJSON(&userLogin)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := utils.FindFirst(&database.UserData, func(u database.User) bool {
		return u.Username == userLogin.Username
	})

	if user == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if !auth.CheckPasswordHash(userLogin.Password, (*user).Password) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.IndentedJSON(http.StatusOK, auth.GenerateJWT((*user).Id))
}

func GetUsers(c *gin.Context) {
	users := make([]returnUser, 0)

	for _, user := range database.UserData {
		users = append(users, returnUser{Username: user.Username})
	}

	c.IndentedJSON(http.StatusOK, users)
}
