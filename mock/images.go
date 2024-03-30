package mock

import (
	"github.com/akrck02/valhalla-core-dal/configuration"
)

func ProfilePicture() string {
	return configuration.BASE_PATH + "resources/images/cat.jpg"
}
