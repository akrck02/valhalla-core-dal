package devicedal

import (
	"github.com/akrck02/valhalla-core-dal/configuration"
	"github.com/akrck02/valhalla-core-dal/database"
	"github.com/akrck02/valhalla-core-sdk/log"
	"github.com/akrck02/valhalla-core-sdk/models"
	"github.com/akrck02/valhalla-core-sdk/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// AddUserDevice adds a new device to the database
// or updates the token if the device already exists
//
// [param] client | *mongo.Client: client to the database
// [param] user | models.User: user that owns the device
// [param] device | models.Device: device to add
//
// [return] string: token of the device --> error : The error that occurred
func AddUserDevice(client *mongo.Client, user *models.User, device *models.Device) (string, error) {

	token, err := utils.GenerateAuthToken(user, device, configuration.Params.Secret)

	if err != nil {
		return "", err
	}

	coll := client.Database(database.CurrentDatabase).Collection(database.DEVICE)
	device.Token = token
	device.User = user.Email

	found, err := FindDevice(coll, device)

	if err != nil && err != mongo.ErrNoDocuments {
		return "", err
	}

	if found != nil {

		log.Debug("Device already exists, updating token")
		coll.ReplaceOne(database.GetDefaultContext(), found, device)

		return token, nil
	}

	log.Debug("Creating new device...")

	_, err = coll.InsertOne(database.GetDefaultContext(), device)

	if err != nil {
		return "", err
	}

	return token, nil
}

// FindDevice finds a device in the database
//
// [param] coll | *mongo.Collection: collection to search
// [param] device | models.Device: device to find
//
// [return] models.Device: device found --> error : The error that occurred
func FindDevice(coll *mongo.Collection, device *models.Device) (*models.Device, error) {

	var found models.Device
	err := coll.FindOne(
		database.GetDefaultContext(),
		bson.M{"user": device.User, "address": device.Address, "useragent": device.UserAgent},
	).Decode(&found)

	if err != nil {
		return nil, err
	}

	return &found, nil
}

// FindDeviceByAuthToken finds a device in the database by its token, user, address and user agent
//
// [param] client | *mongo.Client: client to the database
// [param] token | string: token of the device
//
// [return] models.Device: device found --> error : The error that occurred
func FindDeviceByAuthToken(coll *mongo.Collection, device *models.Device) (*models.Device, error) {

	var found models.Device
	err := coll.FindOne(
		database.GetDefaultContext(),
		bson.M{"user": device.User, "address": device.Address, "useragent": device.UserAgent, "token": device.Token},
	).Decode(&found)

	if err != nil {
		return nil, err
	}

	return &found, nil

}

// DeleteDevice removes a device from the database
//
// [param] client | *mongo.Client: client to the database
// [param] user | models.User: user that owns the device
// [param] device | models.Device: device to remove
//
// [return] *mongo.DeleteResult: result of the operation --> error : The error that occurred
func DeleteDevice(client *mongo.Client, device *models.Device) error {

	coll := client.Database(database.CurrentDatabase).Collection(database.DEVICE)
	_, err := coll.DeleteOne(database.GetDefaultContext(), bson.M{"user": device.User, "address": device.Address, "useragent": device.UserAgent})

	return err
}
