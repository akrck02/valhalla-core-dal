package tests

import (
	"testing"

	"github.com/akrck02/valhalla-core-dal/mock"
	devicedal "github.com/akrck02/valhalla-core-dal/services/device"
	"github.com/akrck02/valhalla-core-sdk/models"
)

// AddMockDeviceToUser adds a mock device to a user
//
// [param] t | *testing.T: testing object
// [param] user | *models.User: user to add the device
//
// [return] device: the device added
func AddMockDeviceToUser(t *testing.T, user *models.User) *models.Device {

	var expected = models.Device{
		Token:     mock.Token(),
		User:      user.Email,
		Address:   mock.Ip(),
		UserAgent: mock.Platform(),
	}

	device, error := devicedal.AddUserDevice(user, &expected)

	if error != nil {
		t.Error(error)
	}

	return device
}
