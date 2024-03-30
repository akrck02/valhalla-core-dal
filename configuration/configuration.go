package configuration

import (
	"os"
	"strings"

	"github.com/akrck02/valhalla-core-sdk/log"
	"github.com/joho/godotenv"
)

type DatabaseConfiguration struct {
	Secret string
	Mongo  string
}

var Params DatabaseConfiguration

func LoadConfiguration(envPath string) {

	err := godotenv.Load(envPath)

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var configuration = DatabaseConfiguration{
		Secret: os.Getenv("SECRET"),
		Mongo:  os.Getenv("IP_MONGODB"),
	}

	checkCompulsoryVariables(configuration)
	Params = configuration
}

func checkCompulsoryVariables(Configuration DatabaseConfiguration) {
	log.Jump()
	log.Line()
	log.Info("Configuration variables")
	log.Line()
	log.Info("SECRET: " + strings.Repeat("*", len(Configuration.Secret)))
	log.Info("MONGO: " + Configuration.Mongo)
}

func IsDevelopment() bool {
	return os.Getenv("ENV") == "development"
}
