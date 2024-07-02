package devicedal

import (
	"github.com/akrck02/valhalla-core-dal/configuration"
	"github.com/akrck02/valhalla-core-dal/database"
	"github.com/akrck02/valhalla-core-sdk/http"
	"github.com/akrck02/valhalla-core-sdk/log"
	devicemodels "github.com/akrck02/valhalla-core-sdk/models/device"
	systemmodels "github.com/akrck02/valhalla-core-sdk/models/system"
	usersmodels "github.com/akrck02/valhalla-core-sdk/models/users"
	"github.com/akrck02/valhalla-core-sdk/utils"
	"github.com/akrck02/valhalla-core-sdk/valerror"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddUserDevice(user *usersmodels.User, device *devicemodels.Device) (*devicemodels.Device, *systemmodels.Error) {

	token, err := utils.GenerateAuthToken(user, device, configuration.Params.Secret)

	if err != nil {
		return nil, &systemmodels.Error{
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
		return device, &systemmodels.Error{
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
func FindDevice(coll *mongo.Collection, device *devicemodels.Device) (*devicemodels.Device, *systemmodels.Error) {

	var found devicemodels.Device
	err := coll.FindOne(
		database.GetDefaultContext(),
		bson.M{"user": device.User, "address": device.Address, "useragent": device.UserAgent},
	).Decode(&found)

	if err != nil {
		return nil, &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.DEVICE_DOES_NOT_EXIST,
			Message: "Device doesn't exists",
		}
	}

	return &found, nil
}

func FindDeviceByAuthToken(coll *mongo.Collection, device *devicemodels.Device) (*devicemodels.Device, *systemmodels.Error) {

	var found devicemodels.Device
	err := coll.FindOne(
		database.GetDefaultContext(),
		bson.M{"user": device.User, "address": device.Address, "useragent": device.UserAgent, "token": device.Token},
	).Decode(&found)

	if err != nil {
		return nil, &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.DEVICE_DOES_NOT_EXIST,
			Message: "Device doesn't exists",
		}
	}

	return &found, nil

}

func DeleteDevice(device *devicemodels.Device) *systemmodels.Error {

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
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   valerror.DATABASE_ERROR,
			Message: "Error deleting device",
		}
	}

	return nil
}

func DeviceExists(device *devicemodels.Device) *systemmodels.Error {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	// check if device exists
	coll := client.Database(database.CurrentDatabase).Collection(database.DEVICE)
	obtained, err := FindDevice(coll, device)

	if err != nil {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.DEVICE_DOES_NOT_EXIST,
			Message: "Device doesn't exist",
		}
	}

	if obtained == nil {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.DEVICE_DOES_NOT_EXIST,
			Message: "Device doesn't exist",
		}
	}

	return nil
}
