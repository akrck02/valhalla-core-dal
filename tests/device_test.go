package tests

import (
	"testing"

	"github.com/akrck02/valhalla-core-dal/mock"
	devicedal "github.com/akrck02/valhalla-core-dal/services/device"
	devicemodels "github.com/akrck02/valhalla-core-sdk/models/device"
)

func TestDeviceExists(t *testing.T) {

	// Create a new user
	user := RegisterMockTestUser(t)
	device := AddMockDeviceToUser(t, user)

	// check if device exists
	err := devicedal.DeviceExists(device)
	if err != nil {
		t.Error("Device not found")
	}

	// delete device
	err = devicedal.DeleteDevice(device)
	if err != nil {
		t.Error(err)
	}

	// delete user
	DeleteTestUser(t, user)

}

func TestDeviceNotExists(t *testing.T) {

	user := RegisterMockTestUser(t)
	device := &devicemodels.Device{
		Token:     mock.Token(),
		User:      user.Email,
		Address:   mock.Ip(),
		UserAgent: mock.Platform(),
	}
	// check if device exists
	err := devicedal.DeviceExists(device)

	if err == nil {
		t.Error("Device found")
	}

	// delete user
	DeleteTestUser(t, user)
}
