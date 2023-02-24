package user

import (
	"cdn/auth"
	"cdn/database"
	"cdn/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

type userRegister struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Admin    bool   `json:"admin"`
}

func Register(c *gin.Context) {
	var newUser userRegister

	user := auth.ExtractUser(c)

	if user == nil || !(*user).Admin {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err := c.BindJSON(&newUser)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.HashPassword(newUser.Password)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	database.UserData = append(database.UserData, database.User{Id: uuid.New().String(), Username: newUser.Username, Password: hashedPassword, Admin: newUser.Admin})

	c.IndentedJSON(http.StatusOK, newUser)
}

func GetUsers(c *gin.Context) {
	users := make([]returnUser, 0)

	for _, user := range database.UserData {
		users = append(users, returnUser{Username: user.Username})
	}

	c.IndentedJSON(http.StatusOK, users)
}
