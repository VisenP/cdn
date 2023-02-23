package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"log"
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

func extractUser(ctx *gin.Context) *user {

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

		return TokenSecret, nil
	})

	if err != nil {
		return nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !(ok && token.Valid) {
		return nil
	}

	return findFirst(&userData, func(u user) bool {
		return u.Id == claims["user_id"]
	})
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
