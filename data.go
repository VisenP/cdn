package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
)

type storedFile struct {
	Name      string `json:"name"`
	Id        string `json:"id"`
	Owner     string `json:"owner"`
	Encrypted bool   `json:"encrypted"`
}

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var fileData []storedFile
var userData []user

func createJsonIfNotExists(name string) {
	_, err := os.Stat("./files/" + name)
	if err == nil {
		return
	} else if errors.Is(err, os.ErrNotExist) {
		log.Println("Creating file: " + name)
		err = os.WriteFile("./files/"+name, []byte("[]"), 0644)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("Unable to create file: " + name)
	}
}

func readFileData(name string) []byte {
	log.Println("Reading file: " + name)
	file, err := os.Open("./files/" + name)
	if err != nil {
		log.Fatal("Unable to open file: " + name)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal("Error closing file: " + name)
		}
	}(file)
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("Unable to read file data: " + name)
	}
	return fileBytes
}

func initializeDataStorage() {
	err := os.MkdirAll("./files", os.ModePerm)
	if err != nil {
		log.Fatal("Could not create directory!")
	}

	createJsonIfNotExists("fileData.json")
	createJsonIfNotExists("userData.json")

	err = json.Unmarshal(readFileData("fileData.json"), &fileData)
	if err != nil {
		log.Fatal("Error reading file data!")
	}
	err = json.Unmarshal(readFileData("userData.json"), &userData)
	if err != nil {
		log.Fatal("Error reading user data!")
	}

	log.Println("Data loaded successfully!")
}
