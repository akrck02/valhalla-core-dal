package devicedal

import (
	"net/http"

	"github.com/akrck02/valhalla-core-dal/configuration"
	"github.com/akrck02/valhalla-core-dal/database"

	apierror "github.com/akrck02/valhalla-core-sdk/error"
	"github.com/akrck02/valhalla-core-sdk/log"
	apimodels "github.com/akrck02/valhalla-core-sdk/models/api"
	devicemodels "github.com/akrck02/valhalla-core-sdk/models/device"
	usersmodels "github.com/akrck02/valhalla-core-sdk/models/users"
	"github.com/akrck02/valhalla-core-sdk/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddUserDevice(conn *mongo.Client, user *usersmodels.User, device *devicemodels.Device) (*devicemodels.Device, *apimodels.Error) {

	token, err := utils.GenerateAuthToken(user, device, configuration.Params.Secret)
	if err != nil {
		return nil, &apimodels.Error{
			Status:  http.StatusInternalServerError,
			Error:   apierror.DatabaseError,
			Message: "Error generating token",
		}
	}

	device.User = user.Email
	coll := conn.Database(database.CurrentDatabase).Collection(database.DEVICE)

	foundDevice, _ := FindDevice(coll, device)
	device.Token = token
	updateDate := utils.GetCurrentMillis()
	if foundDevice != nil {
		log.Debug("Device already exists, updating token")
		device.LastUpdate = &updateDate
		coll.ReplaceOne(database.GetDefaultContext(), foundDevice, device)
		return device, nil
	}

	device.CreationDate = &updateDate
	device.LastUpdate = &updateDate

	log.Debug("Creating new device...")
	_, insertErr := coll.InsertOne(database.GetDefaultContext(), device)
	if insertErr != nil {
		return device, &apimodels.Error{
			Status:  http.StatusInternalServerError,
			Error:   apierror.DatabaseError,
			Message: "Error creating device",
		}
	}

	return device, nil
}

func FindDevice(coll *mongo.Collection, device *devicemodels.Device) (*devicemodels.Device, *apimodels.Error) {

	var found devicemodels.Device
	err := coll.FindOne(
		database.GetDefaultContext(),
		bson.M{"user": device.User, "address": device.Address, "useragent": device.UserAgent},
	).Decode(&found)

	if err != nil {
		return nil, &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   apierror.DeviceNotFound,
			Message: "Device doesn't exists",
		}
	}

	return &found, nil
}

func FindDeviceByAuthToken(coll *mongo.Collection, device *devicemodels.Device) (*devicemodels.Device, *apimodels.Error) {

	var found devicemodels.Device
	err := coll.FindOne(
		database.GetDefaultContext(),
		bson.M{"user": device.User, "address": device.Address, "useragent": device.UserAgent, "token": device.Token},
	).Decode(&found)

	if err != nil {
		return nil, &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   apierror.DeviceNotFound,
			Message: "No possible login devices",
		}
	}

	return &found, nil

}

func DeleteDevice(conn *mongo.Client, device *devicemodels.Device) *apimodels.Error {

	//check if device exists
	existsErr := DeviceExists(conn, device)
	if existsErr != nil {
		return existsErr
	}

	// delete device
	coll := conn.Database(database.CurrentDatabase).Collection(database.DEVICE)
	_, err := coll.DeleteOne(database.GetDefaultContext(), bson.M{"user": device.User, "address": device.Address, "useragent": device.UserAgent})

	if err != nil {
		return &apimodels.Error{
			Status:  http.StatusInternalServerError,
			Error:   apierror.DatabaseError,
			Message: "Error deleting device",
		}
	}

	return nil
}

func DeviceExists(conn *mongo.Client, device *devicemodels.Device) *apimodels.Error {

	// check if device exists
	coll := conn.Database(database.CurrentDatabase).Collection(database.DEVICE)
	obtained, err := FindDevice(coll, device)

	if err != nil {
		return &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   apierror.DeviceNotFound,
			Message: "Device doesn't exist",
		}
	}

	if obtained == nil {
		return &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   apierror.DeviceNotFound,
			Message: "Device doesn't exist",
		}
	}

	return nil
}
