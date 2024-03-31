package permissiondal

import "github.com/akrck02/valhalla-core-sdk/models"

func CanEditUser(author *models.User, user *models.User) bool {
	return author.Email == user.Email
}

func CanSeeUser(author *models.User, user *models.User) bool {
	return author.Email == user.Email
}
