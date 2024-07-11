package permissiondal

import (
	projectmodels "github.com/akrck02/valhalla-core-sdk/models/project"
	usersmodels "github.com/akrck02/valhalla-core-sdk/models/users"
)

func CanEditUser(author *usersmodels.User, user *usersmodels.User) bool {
	return author.Email == user.Email
}

func CanSeeUser(author *usersmodels.User, user *usersmodels.User) bool {
	return author.Email == user.Email
}
func CanDeleteProject(logedUser *usersmodels.User, project *projectmodels.Project) bool {
	return logedUser.ID == project.Owner
}

//Projects

// func CanEditProject(author *models.Project, user *models.Project) bool {

// 	return author.Email == user.Email
// }

// func CanSeeProject(author *models.Project, user *models.Project) bool {
// 	return author.Email == user.Email
// }
