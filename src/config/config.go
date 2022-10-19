package config

import (
	"log"

	"github.com/joho/godotenv"
)

func ReadEnvfile() error {
	//Read the .env file
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Loading .env file failed")
		return err
	}
	if err == nil {
		log.Println("Loading .env file success")
	}
	return nil
}
