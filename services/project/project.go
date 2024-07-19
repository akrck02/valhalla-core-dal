package projectdal

import (
	"github.com/akrck02/valhalla-core-dal/database"
	userdal "github.com/akrck02/valhalla-core-dal/services/user"
	"github.com/akrck02/valhalla-core-sdk/http"
	projectmodels "github.com/akrck02/valhalla-core-sdk/models/project"
	systemmodels "github.com/akrck02/valhalla-core-sdk/models/system"
	usersmodels "github.com/akrck02/valhalla-core-sdk/models/users"
	"github.com/akrck02/valhalla-core-sdk/utils"
	"github.com/akrck02/valhalla-core-sdk/valerror"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateProject(conn *mongo.Client, project *projectmodels.Project) *systemmodels.Error {

	// if the project name is empty, return an error
	if utils.IsEmpty(project.Name) {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.EMPTY_PROJECT_NAME,
			Message: "Project name cannot be empty",
		}
	}

	// if the project description is empty, return an error
	if utils.IsEmpty(project.Description) {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.EMPTY_PROJECT_DESCRIPTION,
			Message: "Project description cannot be empty",
		}
	}

	// if the project owner is empty, return an error
	if utils.IsEmpty(project.Owner) {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.EMPTY_PROJECT_OWNER,
			Message: "Owner cannot be empty",
		}
	}

	// Check if the user exists and is valid
	owner, userGetError := userdal.GetUser(conn, &usersmodels.User{ID: project.Owner}, true)
	if userGetError != nil {
		return userGetError
	}

	// convert the owner ID to an object ID
	_, parsingError := utils.StringToObjectId(owner.ID)
	if parsingError != nil {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.INVALID_OBJECT_ID,
			Message: "Invalid ID",
		}
	}

	// Check if the project already exists
	// it is not possible to have two projects with the same name and owner
	coll := conn.Database(database.CurrentDatabase).Collection(database.PROJECT)
	found := projectNameExists(coll, project.Name, owner.ID)

	if found {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_CONFLICT,
			Error:   valerror.PROJECT_ALREADY_EXISTS,
			Message: "Project already exists",
		}
	}

	// Insert project into database
	_, insertError := coll.InsertOne(database.GetDefaultContext(), project)

	if insertError != nil {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   valerror.PROJECT_ALREADY_EXISTS,
			Message: "Project already exists",
		}
	}

	return nil
}

func EditProject(conn *mongo.Client, project *projectmodels.Project) *systemmodels.Error {

	return nil
}

func DeleteProject(conn *mongo.Client, project *projectmodels.Project) *systemmodels.Error {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	if utils.IsEmpty(project.Name) {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.EMPTY_PROJECT_NAME,
			Message: "Project name cannot be empty",
		}
	}

	// Check if the project ID is valid
	id, parsingError := utils.StringToObjectId(project.ID)
	if parsingError != nil {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.INVALID_OBJECT_ID,
			Message: "Invalid ID",
		}
	}

	// Delete project
	projects := client.Database(database.CurrentDatabase).Collection(database.PROJECT)
	deleteResult, err := projects.DeleteOne(database.GetDefaultContext(), bson.M{"ID": id})

	// If an error occurs, return the error
	if err != nil {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   valerror.PROJECT_NOT_DELETED,
			Message: "Project not deleted",
		}
	}

	// If the project is not found, return an error
	if deleteResult.DeletedCount == 0 {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_NOT_FOUND,
			Error:   valerror.PROJECT_NOT_FOUND,
			Message: "Project not found",
		}
	}

	return nil
}

func GetProject(conn *mongo.Client, project *projectmodels.Project) (*projectmodels.Project, *systemmodels.Error) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	// Check if the project ID is valid
	projectIdObject, parsingError := utils.StringToObjectId(project.ID)
	if parsingError != nil {
		return nil, &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.INVALID_OBJECT_ID,
			Message: "Invalid ID",
		}
	}

	// Get project by id
	projects := client.Database(database.CurrentDatabase).Collection(database.PROJECT)
	found := projectmodels.Project{}
	err := projects.FindOne(database.GetDefaultContext(), bson.M{"_id": projectIdObject}).Decode(&found)

	// If an error occurs, return the error
	if err != nil {
		return nil, &systemmodels.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   valerror.UNEXPECTED_ERROR,
			Message: err.Error(),
		}
	}

	// If the project is not found, return an error
	if utils.IsEmpty(found.ID) {
		return nil, &systemmodels.Error{
			Status:  http.HTTP_STATUS_NOT_FOUND,
			Error:   valerror.PROJECT_NOT_FOUND,
			Message: "Project not found",
		}
	}

	return &found, nil
}

func GetUserProjects(conn *mongo.Client, ownerId string) []projectmodels.Project {

	// Check if the owner ID is valid
	_, parsingError := utils.StringToObjectId(ownerId)
	if parsingError != nil {
		return nil
	}

	// Check the projects belonging to the owner
	projects := conn.Database(database.CurrentDatabase).Collection(database.PROJECT)
	filter := bson.M{"owner": ownerId}
	cursor, err := projects.Find(database.GetDefaultContext(), filter)
	if err != nil {
		return nil
	}

	var result []projectmodels.Project
	cursor.All(database.GetDefaultContext(), &result)

	return result
}

func projectNameExists(coll *mongo.Collection, name string, owner string) bool {
	filter := bson.M{"name": name, "owner": owner}

	var result *projectmodels.Project
	coll.FindOne(database.GetDefaultContext(), filter).Decode(&result)

	return result != nil
}
