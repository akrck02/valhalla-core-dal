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

	client, err := mongo.NewClient(options.Client().ApplyURI(MONGO_URL + MONGO_USER + ":" + MONGO_PASSWORD + "@" + host + ":" + MONGO_PORT))
	if err != nil {
		log.Fatal(err.Error())
	}
	return client
}

func Connect(client mongo.Client) context.Context {

	ctx := context.Background() //context.WithTimeout(context.Background(), 10*time.Second)
	err := client.Connect(ctx)

	if err != nil {
		log.Fatal(err.Error())
	}

	log.FormattedInfo("Database (${0}) connected on mongodb [${1}:${2}]", CurrentDatabase, configuration.Params.Mongo, MONGO_PORT)
	return ctx
}

func SetupTest() {

	CurrentDatabase = TEST_DATABASE_NAME
	var client = CreateClient()
	var ctx = Connect(*client)

	log.Info("Dropping database " + CurrentDatabase)
	err := client.Database(CurrentDatabase).Drop(ctx)
	if err != nil {
		log.FormattedError("Error dropping database ${0} : ${1}", CurrentDatabase, err.Error())
	}

	defer Disconnect(*client, ctx)
}

func Disconnect(client mongo.Client, ctx context.Context) {
	defer client.Disconnect(ctx)
}

func Setup() {

	var client = CreateClient()
	var ctx = Connect(*client)

	defer Disconnect(*client, ctx)

}
