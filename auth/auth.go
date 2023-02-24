package auth

import (
	"cdn/database"
	"cdn/global"
	"cdn/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
)

func GenerateJWT(userId string) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userId

	tokenString, err := token.SignedString(global.TokenSecret)

	if err != nil {
		log.Println(err)
		log.Fatal("Failed to generate JWT!")
	}

	return tokenString
}

func ExtractUser(ctx *gin.Context) *database.User {

	authHeader := ctx.GetHeader("Authorization")
	tokenString, foundPrefix := strings.CutPrefix(authHeader, "Bearer ")

	if !foundPrefix {
		return nil
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("sigining method error")
		}

		return global.TokenSecret, nil
	})

	if err != nil {
		return nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !(ok && token.Valid) {
		return nil
	}

	return utils.FindFirst(&database.UserData, func(u database.User) bool {
		return u.Id == claims["user_id"]
	})
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
