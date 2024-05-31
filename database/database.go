package database

import (
	"context"

	"github.com/akrck02/valhalla-core-dal/configuration"
	"github.com/akrck02/valhalla-core-sdk/log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const MONGO_URL = "mongodb://"
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

	uri := MONGO_URL + configuration.Params.User + ":" + configuration.Params.Password + "@" + configuration.Params.Mongo + ":" + configuration.Params.MongoPort
	log.Info(uri)
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err.Error())
	}

	return client
}

func Connect() *mongo.Client {

	client := CreateClient()
	log.FormattedInfo("Database (${0}) connected on mongodb [${1}:${2}]", CurrentDatabase, configuration.Params.Mongo, configuration.Params.MongoPort)
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
	err := client.Ping(GetDefaultContext(), &readpref.ReadPref{})

	if err != nil {
		log.Error("cannot connect to database")
	}

	defer Disconnect(*client)
}

func GetDefaultContext() context.Context {
	return context.Background()
}
