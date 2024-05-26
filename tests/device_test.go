package tests

import (
	"log"
	"testing"

	"github.com/akrck02/valhalla-core-dal/database"
	"github.com/akrck02/valhalla-core-dal/mock"
	devicedal "github.com/akrck02/valhalla-core-dal/services/device"
	userdal "github.com/akrck02/valhalla-core-dal/services/user"
	"github.com/akrck02/valhalla-core-sdk/models"
)

func TestDeviceExists(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	// Create user
	user := models.User{
		Email:    mock.Email(),
		Username: mock.Username(),
		Password: mock.Password(),
	}

	err := userdal.Register(client, &user)

	if err != nil {
		t.Error(err)
	}

	// add device to user
	var expected = models.Device{
		Token:     mock.Token(),
		User:      user.Email,
		Address:   mock.Ip(),
		UserAgent: mock.Platform(),
	}

	_, error := devicedal.AddUserDevice(client, &user, &expected)

	if error != nil {
		t.Error(err)
	}

	// check if device exists
	coll := client.Database(database.CurrentDatabase).Collection(database.DEVICE)
	obtained, error := devicedal.FindDevice(coll, &expected)

	if error != nil {
		t.Error(err)
	}

	if obtained == nil {
		t.Error("Device not found")
	}

	log.Print("Device expected: ", expected)
	log.Print("Device found: ", obtained)

	// delete device
	error = devicedal.DeleteDevice(client, &expected)

	if error != nil {
		t.Error(err)
	}

	// delete user
	err = userdal.DeleteUser(client, &user)

	if err != nil {
		t.Error(err)
	}

}

func TestDeviceNotExists(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	// Create user
	user := models.User{
		Email:    mock.Email(),
		Username: mock.Username(),
		Password: mock.Password(),
	}

	err := userdal.Register(client, &user)

	if err != nil {
		t.Error(err)
	}

	// add device to user
	var expected = models.Device{
		Token:     mock.Token(),
		User:      user.Email,
		Address:   mock.Ip(),
		UserAgent: mock.Platform(),
	}

	_, error := devicedal.AddUserDevice(client, &user, &expected)

	if error != nil {
		t.Error(err)
	}

	// check if device exists
	coll := client.Database(database.CurrentDatabase).Collection(database.DEVICE)
	obtained, error := devicedal.FindDevice(coll, &models.Device{
		Token: mock.Token(),
	})

	if error == nil || obtained != nil {
		t.Error("Device found")
	}

	log.Print("Device expected: ", expected)
	log.Print("Device not found: ", obtained)

	// delete device
	error = devicedal.DeleteDevice(client, &expected)

	if error != nil {
		t.Error(err)
	}

	// delete user
	err = userdal.DeleteUser(client, &user)

	if err != nil {
		t.Error(err)
	}

}
