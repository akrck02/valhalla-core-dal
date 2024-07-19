package tests

import (
	"testing"

	"github.com/akrck02/valhalla-core-dal/database"
	"github.com/akrck02/valhalla-core-dal/mock"
	devicedal "github.com/akrck02/valhalla-core-dal/services/device"
	devicemodels "github.com/akrck02/valhalla-core-sdk/models/device"
)

func TestDeviceExists(t *testing.T) {

	conn := database.Connect()

	// Create a new user
	user := RegisterMockTestUser(conn, t)
	device := AddMockDeviceToUser(conn, t, user)

	// check if device exists
	err := devicedal.DeviceExists(conn, device)
	if err != nil {
		t.Error("Device not found")
	}

	// delete device
	err = devicedal.DeleteDevice(conn, device)
	if err != nil {
		t.Error(err)
	}
}

func TestDeviceNotExists(t *testing.T) {

	conn := database.Connect()
	user := RegisterMockTestUser(conn, t)
	device := &devicemodels.Device{
		Token:     mock.Token(),
		User:      user.Email,
		Address:   mock.Ip(),
		UserAgent: mock.Platform(),
	}
	// check if device exists
	err := devicedal.DeviceExists(conn, device)

	if err == nil {
		t.Error("Device found")
	}

}
