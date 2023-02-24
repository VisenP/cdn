package database

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"time"
)

type StoredFile struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Ext       string `json:"ext"`
	Public    bool   `json:"public"`
	Owner     string `json:"owner"`
	Encrypted bool   `json:"encrypted"`
}

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Admin    bool   `json:"admin"`
}

var FileData []StoredFile
var UserData []User

func save() {

	fileDataJson, _ := json.Marshal(FileData)
	userDataJson, _ := json.Marshal(UserData)

	err := os.WriteFile("./files/fileData.json", fileDataJson, 0644)
	if err != nil {
		log.Fatal("Error saving file data: " + err.Error())
	}
	err = os.WriteFile("./files/userData.json", userDataJson, 0644)
	if err != nil {
		log.Fatal("Error saving user data: " + err.Error())
	}

}

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

func InitializeDataStorage() {
	err := os.MkdirAll("./files", os.ModePerm)
	if err != nil {
		log.Fatal("Could not create directory!")
	}

	createJsonIfNotExists("fileData.json")
	createJsonIfNotExists("userData.json")

	err = json.Unmarshal(readFileData("fileData.json"), &FileData)
	if err != nil {
		log.Fatal("Error reading file data!")
	}
	err = json.Unmarshal(readFileData("userData.json"), &UserData)
	if err != nil {
		log.Fatal("Error reading user data!")
	}

	log.Println("Data loaded successfully!")

	go func() {
		for true {
			save()
			time.Sleep(10 * time.Second)
		}
	}()
}
