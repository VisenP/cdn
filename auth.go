package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"strings"
	"time"
)

func generateJWT(userId string) string {
	token := jwt.New(jwt.SigningMethodEdDSA)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 12)
	claims["user_id"] = userId

	tokenString, err := token.SignedString(TOKEN_SECRET)

	if err != nil {
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

		_, ok := token.Method.(*jwt.SigningMethodECDSA)

		if !ok {
			return nil, errors.New("sigining method error")
		}

		return TOKEN_SECRET, nil
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
