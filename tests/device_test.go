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
	var client = database.CreateClient()
	var conn = database.Connect(*client)
	defer database.Disconnect(*client, conn)

	// Create user
	user := models.User{
		Email:    mock.Email(),
		Username: mock.Username(),
		Password: mock.Password(),
	}

	err := userdal.Register(conn, client, &user)

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

	_, error := devicedal.AddUserDevice(conn, client, &user, &expected)

	if error != nil {
		t.Error(err)
	}

	// check if device exists
	coll := client.Database(database.CurrentDatabase).Collection(database.DEVICE)
	obtained, error := devicedal.FindDevice(conn, coll, &expected)

	if error != nil {
		t.Error(err)
	}

	if obtained == nil {
		t.Error("Device not found")
	}

	log.Print("Device expected: ", expected)
	log.Print("Device found: ", obtained)

	// delete device
	error = devicedal.DeleteDevice(conn, client, &expected)

	if error != nil {
		t.Error(err)
	}

	// delete user
	err = userdal.DeleteUser(conn, client, &user)

	if err != nil {
		t.Error(err)
	}

}

func TestDeviceNotExists(t *testing.T) {

	// Connect database
	var client = database.CreateClient()
	var conn = database.Connect(*client)
	defer database.Disconnect(*client, conn)

	// Create user
	user := models.User{
		Email:    mock.Email(),
		Username: mock.Username(),
		Password: mock.Password(),
	}

	err := userdal.Register(conn, client, &user)

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

	_, error := devicedal.AddUserDevice(conn, client, &user, &expected)

	if error != nil {
		t.Error(err)
	}

	// check if device exists
	coll := client.Database(database.CurrentDatabase).Collection(database.DEVICE)
	obtained, error := devicedal.FindDevice(conn, coll, &models.Device{
		Token: mock.Token(),
	})

	if error == nil || obtained != nil {
		t.Error("Device found")
	}

	log.Print("Device expected: ", expected)
	log.Print("Device not found: ", obtained)

	// delete device
	error = devicedal.DeleteDevice(conn, client, &expected)

	if error != nil {
		t.Error(err)
	}

	// delete user
	err = userdal.DeleteUser(conn, client, &user)

	if err != nil {
		t.Error(err)
	}

}
