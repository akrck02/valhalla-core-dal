package projectdal

import (
	"context"

	"github.com/akrck02/valhalla-core-dal/database"
	"github.com/akrck02/valhalla-core-sdk/error"
	"github.com/akrck02/valhalla-core-sdk/http"
	"github.com/akrck02/valhalla-core-sdk/models"
	"github.com/akrck02/valhalla-core-sdk/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Create project logic
//
// [param] conn | context.Context : The connection context
// [param] client | *mongo.Client : The mongo client
// [param] project | models.Project : The project to create
//
// [return] *models.Error : The error
func CreateProject(conn context.Context, client *mongo.Client, project *models.Project) *models.Error {

	if utils.IsEmpty(project.Name) {
		return &models.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   int(error.EMPTY_PROJECT_NAME),
			Message: "Project name cannot be empty",
		}
	}

	if utils.IsEmpty(project.Description) {
		return &models.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   int(error.EMPTY_PROJECT_DESCRIPTION),
			Message: "Project description cannot be empty",
		}
	}

	if utils.IsEmpty(project.Owner) {
		return &models.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   int(error.EMPTY_PROJECT_OWNER),
			Message: "Owner cannot be empty",
		}
	}

	coll := client.Database(database.CurrentDatabase).Collection(database.PROJECT)
	found := nameExists(conn, coll, project.Name, project.Owner)

	if found {
		return &models.Error{
			Status:  http.HTTP_STATUS_CONFLICT,
			Error:   int(error.PROJECT_ALREADY_EXISTS),
			Message: "Project already exists",
		}
	}

	_, err := coll.InsertOne(conn, project)

	if err != nil {
		return &models.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   int(error.PROJECT_ALREADY_EXISTS),
			Message: "Project already exists",
		}
	}

	return nil
}

// Edit project logic
//
// [param] conn | context.Context : The connection context
// [param] client | *mongo.Client : The mongo client
// [param] project | models.Project : The project to edit
//
// [return] *models.Error : The error
func EditProject(conn context.Context, client *mongo.Client, project *models.Project) *models.Error {

	return nil
}

// Delete project logic
//
// [param] conn | context.Context : The connection context
// [param] client | *mongo.Client : The mongo client
// [param] project | models.Project : The project to delete
//
// [return] *models.Error : The error
func DeleteProject(conn context.Context, client *mongo.Client, project *models.Project) *models.Error {

	if utils.IsEmpty(project.Name) {
		return &models.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   int(error.EMPTY_PROJECT_NAME),
			Message: "Project name cannot be empty",
		}
	}

	// delete user devices
	devices := client.Database(database.CurrentDatabase).Collection(database.DEVICE)
	_, err := devices.DeleteMany(conn, bson.M{"project": project.Name})

	if err != nil {
		return &models.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   int(error.PROJECT_NOT_DELETED),
			Message: "Project not deleted",
		}
	}

	projects := client.Database(database.CurrentDatabase).Collection(database.PROJECT)

	var deleteResult *mongo.DeleteResult
	deleteResult, err = projects.DeleteOne(conn, bson.M{"name": project.Name})

	if err != nil {
		return &models.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   int(error.PROJECT_NOT_DELETED),
			Message: "Project not deleted",
		}
	}

	if deleteResult.DeletedCount == 0 {
		return &models.Error{
			Status:  http.HTTP_STATUS_NOT_FOUND,
			Error:   int(error.PROJECT_NOT_FOUND),
			Message: "Project not found",
		}
	}

	return nil
}

// Get project logic
//
// [param] conn | context.Context : The connection context
// [param] client | *mongo.Client : The mongo client
// [param] project | models.Project : The project to get
//
// [return] *models.Error : The error
func GetProject(conn context.Context, client *mongo.Client, project *models.Project) (*models.Project, *models.Error) { // get project from database

	projects := client.Database(database.CurrentDatabase).Collection(database.PROJECT)

	found := models.Project{}
	err := projects.FindOne(conn, bson.M{"name": project.Name}).Decode(&found)

	if err != nil {
		return nil, &models.Error{
			Status:  http.HTTP_STATUS_NOT_FOUND,
			Error:   int(error.PROJECT_NOT_FOUND),
			Message: "Project not found",
		}
	}

	return &found, nil
}

// Get all projects by user logic
//
// [param] conn | context.Context : The connection context
// [param] client | *mongo.Client : The mongo client
// [param] email | string : The email of the user
//
// [return] []models.Project : The projects of the user
func GetUserProjects(conn context.Context, client *mongo.Client, email string) []models.Project {

	projects := client.Database(database.CurrentDatabase).Collection(database.PROJECT)

	filter := bson.M{"owner": email}

	cursor, err := projects.Find(conn, filter)

	if err != nil {
		return nil
	}

	var result []models.Project
	cursor.All(conn, &result)

	return result
}

// Check if project name exists
//
// [param] conn | context.Context : The connection context
// [param] coll | *mongo.Collection : The mongo collection
// [param] name | string : The name of the project
// [param] owner | string : The owner of the project
//
// [return] models.Project : The project found
func nameExists(conn context.Context, coll *mongo.Collection, name string, owner string) bool {
	filter := bson.M{"name": name, "owner": owner}

	var result *models.Project
	coll.FindOne(conn, filter).Decode(&result)

	return result != nil
}
