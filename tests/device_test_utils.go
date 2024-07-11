package tests

import (
	"testing"

	"github.com/akrck02/valhalla-core-dal/mock"
	devicedal "github.com/akrck02/valhalla-core-dal/services/device"
	devicemodels "github.com/akrck02/valhalla-core-sdk/models/device"
	usersmodels "github.com/akrck02/valhalla-core-sdk/models/users"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddMockDeviceToUser(conn *mongo.Client, t *testing.T, user *usersmodels.User) *devicemodels.Device {

	var expected = devicemodels.Device{
		Token:     mock.Token(),
		User:      user.Email,
		Address:   mock.Ip(),
		UserAgent: mock.Platform(),
	}

	device, error := devicedal.AddUserDevice(conn, user, &expected)

	if error != nil {
		t.Error(error)
	}

	return device
}
