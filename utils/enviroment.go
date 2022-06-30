package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnviromentVars(key string) string {
	_, isPresent := os.LookupEnv(key)

	if !isPresent {
		err := godotenv.Load(".env")

		if err != nil {
			log.Panicf("some error ocurred parsing .env %s", err)
		}
	}

	return os.Getenv(key)
}
