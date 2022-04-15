package client

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func loadLocalEnv() interface{} {
	if _, runningInContainer := os.LookupEnv("CONTAINER"); !runningInContainer {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("Load .env from env_file specified in docker compose")
	}
	return nil
}

func GetEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatal("Environment variable not found: ", key)
	}
	return value
}