package configuration

var (
	BASE_PATH             = ""
	ENV_PATH              = ""
	RESOURCES_PATH        = ""
	IMAGES_PATH           = ""
	PROFILE_PICTURES_PATH = ""
)

func SetBasePath(path string) int {

	BASE_PATH = path
	ENV_PATH = BASE_PATH + "/.env"
	RESOURCES_PATH = BASE_PATH + "/resources/"
	IMAGES_PATH = RESOURCES_PATH + "images/"
	PROFILE_PICTURES_PATH = RESOURCES_PATH + "profile_pictures/"

	return 0
}
