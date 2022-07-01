package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnviromentVars(key string) string {
	value, isPresent := os.LookupEnv(key)

	if !isPresent {
		err := godotenv.Load(".env")

		if err != nil {
			log.Panicf("some error ocurred parsing .env %s", err)
		}
	} else {
		return value
	}

	return os.Getenv(key)
}
