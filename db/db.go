package db

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func DBConnection() map[string]interface{} {

	Connection := make(map[string]interface{})
	errr := godotenv.Load()
	if errr != nil {
		log.Fatal(errr)
	}
	mongoConn := Mongo{
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Cluster:  os.Getenv("DB_CLUSTER"),
	}

	mongoConnection := mongoConn

	Connection["mongodb"] = mongoConnection

	return Connection
}
