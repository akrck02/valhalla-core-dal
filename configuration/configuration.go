package configuration

import (
	"os"
	"strings"

	"github.com/akrck02/valhalla-core-sdk/log"
	"github.com/joho/godotenv"
)

type DatabaseConfiguration struct {
	Secret    string
	Mongo     string
	MongoPort string
	User      string
	Password  string
}

var Params DatabaseConfiguration

func LoadConfiguration(envPath string) {

	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	databaseConfiguration := DatabaseConfiguration{
		Secret:    os.Getenv("SECRET"),
		Mongo:     os.Getenv("MONGO_SERVER"),
		MongoPort: os.Getenv("MONGO_PORT"),
		User:      os.Getenv("MONGO_USER"),
		Password:  os.Getenv("MONGO_PASSWORD"),
	}

	checkCompulsoryVariables(databaseConfiguration)
	Params = databaseConfiguration
}

func checkCompulsoryVariables(configuration DatabaseConfiguration) {
	log.Jump()
	log.Line()
	log.Info("Configuration variables")
	log.Line()
	log.Info("SECRET: " + strings.Repeat("*", len(configuration.Secret)))
	log.Info("MONGO: " + configuration.Mongo + ":" + configuration.MongoPort)
	log.Info("USER: " + configuration.User)
}

func IsDevelopment() bool {
	return os.Getenv("ENV") == "development"
}
