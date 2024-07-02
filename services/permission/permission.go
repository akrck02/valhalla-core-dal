package permissiondal

import (
	usersmodels "github.com/akrck02/valhalla-core-sdk/models/users"
)

func CanEditUser(author *usersmodels.User, user *usersmodels.User) bool {
	return author.Email == user.Email
}

func CanSeeUser(author *usersmodels.User, user *usersmodels.User) bool {
	return author.Email == user.Email
}
