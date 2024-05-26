package database

import (
	"context"

	"github.com/akrck02/valhalla-core-dal/configuration"
	"github.com/akrck02/valhalla-core-sdk/log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MONGO_URL = "mongodb://"
const MONGO_USER = "admin"
const MONGO_PASSWORD = "p4ssw0rd"
const MONGO_PORT = "27017"

const TEST_DATABASE_NAME = "valhalla-test"

const USER = "user"
const DEVICE = "device"
const TEAM = "team"
const PROJECT = "project"
const TASK = "task"
const NOTE = "note"
const WIKI = "wiki"
const ROLE = "role"

var CurrentDatabase = "valhalla"

func CreateClient() *mongo.Client {

	var host = configuration.Params.Mongo
	clientOptions := options.Client().ApplyURI(MONGO_URL + MONGO_USER + ":" + MONGO_PASSWORD + "@" + host + ":" + MONGO_PORT)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err.Error())
	}

	return client
}

func Connect() *mongo.Client {

	client := CreateClient()
	log.FormattedInfo("Database (${0}) connected on mongodb [${1}:${2}]", CurrentDatabase, configuration.Params.Mongo, MONGO_PORT)
	return client

}

func SetupTest() {

	CurrentDatabase = TEST_DATABASE_NAME
	var client = Connect()

	log.Info("Dropping database " + CurrentDatabase)
	err := client.Database(CurrentDatabase).Drop(context.Background())
	if err != nil {
		log.FormattedError("Error dropping database ${0} : ${1}", CurrentDatabase, err.Error())
	}

	defer Disconnect(*client)
}

func Disconnect(client mongo.Client) {
	defer client.Disconnect(context.Background())
}

func Setup() {

	var client = CreateClient()
	defer Disconnect(*client)

}

func GetDefaultContext() context.Context {
	return context.Background()
}
