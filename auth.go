package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
)

func generateJWT(userId string) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userId

	log.Println(TokenSecret)
	tokenString, err := token.SignedString(TokenSecret)

	if err != nil {
		log.Println(err)
		log.Fatal("Failed to generate JWT!")
	}

	return tokenString
}

func validateUser(ctx *gin.Context) {

	authHeader := ctx.GetHeader("Authorization")
	tokenString, foundPrefix := strings.CutPrefix(authHeader, "Bearer ")

	if !foundPrefix {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("sigining method error")
		}

		return TokenSecret, nil
	})

	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !(ok && token.Valid) {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx.Set("user_id", claims["user_id"].(string))

	ctx.Next()
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
