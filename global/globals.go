package global

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var TokenSecret []byte
var Port string

func InitGlobals() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Could not load .env!")
	}

	log.Println("Env file loaded!")

	TokenSecret = []byte(os.Getenv("TOKEN_SECRET"))
	Port = os.Getenv("PORT")
}
