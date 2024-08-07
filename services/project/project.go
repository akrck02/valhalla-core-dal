package projectdal

import (
	"net/http"

	"github.com/akrck02/valhalla-core-dal/database"
	userdal "github.com/akrck02/valhalla-core-dal/services/user"

	apierror "github.com/akrck02/valhalla-core-sdk/error"
	apimodels "github.com/akrck02/valhalla-core-sdk/models/api"
	projectmodels "github.com/akrck02/valhalla-core-sdk/models/project"
	usersmodels "github.com/akrck02/valhalla-core-sdk/models/users"
	"github.com/akrck02/valhalla-core-sdk/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateProject(conn *mongo.Client, project *projectmodels.Project) *apimodels.Error {

	// if the project name is empty, return an error
	if utils.IsEmpty(project.Name) {
		return &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   apierror.EmptyProjectName,
			Message: "Project name cannot be empty",
		}
	}

	// if the project description is empty, return an error
	if utils.IsEmpty(project.Description) {
		return &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   apierror.EmptyProjectDescription,
			Message: "Project description cannot be empty",
		}
	}

	// if the project owner is empty, return an error
	if utils.IsEmpty(project.Owner) {
		return &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   apierror.EmptyProjectOwner,
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
		return &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   apierror.InvalidObjectId,
			Message: "Invalid ID",
		}
	}

	// Check if the project already exists
	// it is not possible to have two projects with the same name and owner
	coll := conn.Database(database.CurrentDatabase).Collection(database.PROJECT)
	found := projectNameExists(coll, project.Name, owner.ID)

	if found {
		return &apimodels.Error{
			Status:  http.StatusConflict,
			Error:   apierror.ProjectAlreadyExists,
			Message: "Project already exists",
		}
	}

	// Insert project into database
	creationDate := utils.GetCurrentMillis()
	project.CreationDate = &creationDate
	project.LastUpdate = &creationDate
	insertResult, insertError := coll.InsertOne(database.GetDefaultContext(), project)

	if insertError != nil {
		return &apimodels.Error{
			Status:  http.StatusInternalServerError,
			Error:   apierror.UnexpectedError,
			Message: insertError.Error(),
		}
	}

	if insertResult.InsertedID == nil {
		return &apimodels.Error{
			Status:  http.StatusInternalServerError,
			Error:   apierror.UnexpectedError,
			Message: "Project not created",
		}
	}

	return nil
}

func EditProject(conn *mongo.Client, project *projectmodels.Project) *apimodels.Error {

	// Transform team id to object id
	// also check if team id is valid
	objID, err := utils.StringToObjectId(project.ID)

	if err != nil {
		return &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   apierror.InvalidObjectId,
			Message: "Invalid object id",
		}
	}

	coll := conn.Database(database.CurrentDatabase).Collection(database.TEAM)

	_, err = coll.UpdateOne(database.GetDefaultContext(), bson.M{"_id": objID}, bson.M{
		"$set": bson.M{
			"name":        project.Name,
			"description": project.Description,
			"updatedate":  project.LastUpdate,
		},
	})

	// Check if team was updated
	if err != nil {
		return &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   apierror.DatabaseError,
			Message: "Could not update team: " + err.Error(),
		}
	}

	return nil
}

func DeleteProject(conn *mongo.Client, project *projectmodels.Project) *apimodels.Error {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	if utils.IsEmpty(project.Name) {
		return &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   apierror.EmptyProjectName,
			Message: "Project name cannot be empty",
		}
	}

	// Check if the project ID is valid
	id, parsingError := utils.StringToObjectId(project.ID)
	if parsingError != nil {
		return &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   apierror.InvalidObjectId,
			Message: "Invalid ID",
		}
	}

	// Delete project
	projects := client.Database(database.CurrentDatabase).Collection(database.PROJECT)
	deleteResult, err := projects.DeleteOne(database.GetDefaultContext(), bson.M{"ID": id})

	// If an error occurs, return the error
	if err != nil {
		return &apimodels.Error{
			Status:  http.StatusInternalServerError,
			Error:   apierror.DatabaseError,
			Message: "Project not deleted",
		}
	}

	// If the project is not found, return an error
	if deleteResult.DeletedCount == 0 {
		return &apimodels.Error{
			Status:  http.StatusNotFound,
			Error:   apierror.ProjectNotFound,
			Message: "Project not found",
		}
	}

	return nil
}

func GetProject(conn *mongo.Client, project *projectmodels.Project) (*projectmodels.Project, *apimodels.Error) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	// Check if the project ID is valid
	projectIdObject, parsingError := utils.StringToObjectId(project.ID)
	if parsingError != nil {
		return nil, &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   apierror.InvalidObjectId,
			Message: "Invalid ID",
		}
	}

	// Get project by id
	projects := client.Database(database.CurrentDatabase).Collection(database.PROJECT)
	found := projectmodels.Project{}
	err := projects.FindOne(database.GetDefaultContext(), bson.M{"_id": projectIdObject}).Decode(&found)

	// If an error occurs, return the error
	if err != nil {
		return nil, &apimodels.Error{
			Status:  http.StatusInternalServerError,
			Error:   apierror.UnexpectedError,
			Message: err.Error(),
		}
	}

	// If the project is not found, return an error
	if utils.IsEmpty(found.ID) {
		return nil, &apimodels.Error{
			Status:  http.StatusNotFound,
			Error:   apierror.ProjectNotFound,
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
