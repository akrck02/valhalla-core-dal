package permissiondal

import (
	projectmodels "github.com/akrck02/valhalla-core-sdk/models/project"
	usersmodels "github.com/akrck02/valhalla-core-sdk/models/users"
)

// Users

func CanEditUser(author *usersmodels.User, user *usersmodels.User) bool {
	return author.ID == user.ID
}

func CanSeeUser(author *usersmodels.User, user *usersmodels.User) bool {
	return author.ID == user.ID
}
func CanDeleteProject(logedUser *usersmodels.User, project *projectmodels.Project) bool {
	return logedUser.ID == project.Owner
}

//Projects

func CanEditProject(user *usersmodels.User, project *projectmodels.Project) bool {
	return user.ID == project.Owner
}

func CanSeeProject(user *usersmodels.User, project *projectmodels.Project) bool {
	return user.ID == project.Owner
}
