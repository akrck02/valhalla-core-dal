package mock

import (
	"github.com/akrck02/valhalla-core-dal/configuration"
	devicemodels "github.com/akrck02/valhalla-core-sdk/models/device"
	usersmodels "github.com/akrck02/valhalla-core-sdk/models/users"
	"github.com/akrck02/valhalla-core-sdk/utils"
)

func Ip() string {
	return "127.0.0.1"
}

func Platform() string {
	return "Firefox, Linux"
}

func Token() string {

	// Create a user
	var user = &usersmodels.User{
		Username: "#TOKENHASH#",
		Password: "#T0K3NH4SHToKeNHaSH#",
		Email:    "TokenHash@tokenHash.com",
	}

	// Create a device

	var device = &devicemodels.Device{
		UserAgent: "Firefox, Linux",
		Address:   "0.0.0.0",
	}

	// Generate a token
	token, _ := utils.GenerateAuthToken(user, device, configuration.Params.Secret)

	return token

}
