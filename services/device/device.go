package devicedal

import (
	"github.com/akrck02/valhalla-core-dal/configuration"
	"github.com/akrck02/valhalla-core-dal/database"
	"github.com/akrck02/valhalla-core-sdk/http"
	"github.com/akrck02/valhalla-core-sdk/log"
	"github.com/akrck02/valhalla-core-sdk/models"
	devicemodels "github.com/akrck02/valhalla-core-sdk/models/device"
	usersmodels "github.com/akrck02/valhalla-core-sdk/models/users"
	"github.com/akrck02/valhalla-core-sdk/utils"
	"github.com/akrck02/valhalla-core-sdk/valerror"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// AddUserDevice adds a new device to the database
// or updates the token if the device already exists
//
// [param] user | models.User: user that owns the device
// [param] device | models.Device: device to add
//
// [return] device: the device with new token --> error : The error that occurred
func AddUserDevice(user *usersmodels.User, device *devicemodels.Device) (*devicemodels.Device, *models.Error) {

	token, err := utils.GenerateAuthToken(user, device, configuration.Params.Secret)

	if err != nil {
		return nil, &models.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   valerror.DATABASE_ERROR,
			Message: "Error generating token",
		}
	}

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	coll := client.Database(database.CurrentDatabase).Collection(database.DEVICE)
	device.User = user.Email

	found, err := FindDevice(coll, device)

	if found != nil {

		log.Debug("Device already exists, updating token")
		coll.ReplaceOne(database.GetDefaultContext(), found, device)

		return device, nil
	}

	device.Token = token
	log.Debug("Creating new device...")

	_, insertErr := coll.InsertOne(database.GetDefaultContext(), device)

	if insertErr != nil {
		return device, &models.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   valerror.DATABASE_ERROR,
			Message: "Error creating device",
		}
	}

	return device, nil
}

// FindDevice finds a device in the database
//
// [param] coll | *mongo.Collection: collection to search
// [param] device | models.Device: device to find
//
// [return] models.Device: device found --> error : The error that occurred
func FindDevice(coll *mongo.Collection, device *devicemodels.Device) (*devicemodels.Device, *models.Error) {

	var found devicemodels.Device
	err := coll.FindOne(
		database.GetDefaultContext(),
		bson.M{"user": device.User, "address": device.Address, "useragent": device.UserAgent},
	).Decode(&found)

	if err != nil {
		return nil, &models.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.DEVICE_DOES_NOT_EXIST,
			Message: "Device doesn't exists",
		}
	}

	return &found, nil
}

// FindDeviceByAuthToken finds a device in the database by its token, user, address and user agent
//
// [param] token | string: token of the device
//
// [return] models.Device: device found --> error : The error that occurred
func FindDeviceByAuthToken(coll *mongo.Collection, device *devicemodels.Device) (*devicemodels.Device, *models.Error) {

	var found devicemodels.Device
	err := coll.FindOne(
		database.GetDefaultContext(),
		bson.M{"user": device.User, "address": device.Address, "useragent": device.UserAgent, "token": device.Token},
	).Decode(&found)

	if err != nil {
		return nil, &models.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.DEVICE_DOES_NOT_EXIST,
			Message: "Device doesn't exists",
		}
	}

	return &found, nil

}

// DeleteDevice removes a device from the database
//
// [param] user | models.User: user that owns the device
// [param] device | models.Device: device to remove
//
// [return] error: The error that occurred
func DeleteDevice(device *devicemodels.Device) *models.Error {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	//check if device exists
	existsErr := DeviceExists(device)
	if existsErr != nil {
		return existsErr
	}

	// delete device
	coll := client.Database(database.CurrentDatabase).Collection(database.DEVICE)
	_, err := coll.DeleteOne(database.GetDefaultContext(), bson.M{"user": device.User, "address": device.Address, "useragent": device.UserAgent})

	if err != nil {
		return &models.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   valerror.DATABASE_ERROR,
			Message: "Error deleting device",
		}
	}

	return nil
}

// DeviceExists checks if a device exists in the database
//
// [param] device | models.Device: device to check
//
// [return] error: The error that occurred
func DeviceExists(device *devicemodels.Device) *models.Error {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	// check if device exists
	coll := client.Database(database.CurrentDatabase).Collection(database.DEVICE)
	obtained, err := FindDevice(coll, device)

	if err != nil {
		return &models.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.DEVICE_DOES_NOT_EXIST,
			Message: "Device doesn't exist",
		}
	}

	if obtained == nil {
		return &models.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.DEVICE_DOES_NOT_EXIST,
			Message: "Device doesn't exist",
		}
	}

	return nil
}
