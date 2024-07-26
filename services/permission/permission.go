package permissiondal

import (
	projectmodels "github.com/akrck02/valhalla-core-sdk/models/project"
	usersmodels "github.com/akrck02/valhalla-core-sdk/models/users"
)

// Users

func CanEditUser(author *usersmodels.User, user *usersmodels.User) bool {
	return author.Id == user.Id
}

func CanSeeUser(author *usersmodels.User, user *usersmodels.User) bool {
	return author.Id == user.Id
}
func CanDeleteProject(logedUser *usersmodels.User, project *projectmodels.Project) bool {
	return logedUser.Id == project.Owner
}

//Projects

func CanEditProject(user *usersmodels.User, project *projectmodels.Project) bool {
	return user.Id == project.Owner
}

func CanSeeProject(user *usersmodels.User, project *projectmodels.Project) bool {
	return user.Id == project.Owner
}
