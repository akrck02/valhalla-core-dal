package projectdal

import (
	"log"
	"net/http"

	"github.com/akrck02/valhalla-core-dal/database"
	userdal "github.com/akrck02/valhalla-core-dal/services/user"

	apierror "github.com/akrck02/valhalla-core-sdk/error"
	apimodels "github.com/akrck02/valhalla-core-sdk/models/api"
	projectmodels "github.com/akrck02/valhalla-core-sdk/models/project"
	usersmodels "github.com/akrck02/valhalla-core-sdk/models/users"
	"github.com/akrck02/valhalla-core-sdk/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateProject(conn *mongo.Client, project *projectmodels.Project) (*projectmodels.Project, *apimodels.Error) {

	// if the project name is empty, return an error
	if utils.IsEmpty(project.Name) {
		return nil, &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   apierror.EmptyProjectName,
			Message: "Project name cannot be empty",
		}
	}

	// if the project description is empty, return an error
	if utils.IsEmpty(project.Description) {
		return nil, &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   apierror.EmptyProjectDescription,
			Message: "Project description cannot be empty",
		}
	}

	// if the project owner is empty, return an error
	if utils.IsEmpty(project.Owner) {
		return nil, &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   apierror.EmptyProjectOwner,
			Message: "Owner cannot be empty",
		}
	}

	// Check if the user exists and is valid
	owner, userGetError := userdal.GetUser(conn, &usersmodels.User{ID: project.Owner}, true)
	if userGetError != nil {
		return nil, userGetError
	}

	// convert the owner ID to an object ID
	_, parsingError := utils.StringToObjectId(owner.ID)
	if parsingError != nil {
		return nil, &apimodels.Error{
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
		return nil, &apimodels.Error{
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
		return nil, &apimodels.Error{
			Status:  http.StatusInternalServerError,
			Error:   apierror.UnexpectedError,
			Message: insertError.Error(),
		}
	}

	if insertResult.InsertedID == nil {
		return nil, &apimodels.Error{
			Status:  http.StatusInternalServerError,
			Error:   apierror.UnexpectedError,
			Message: "Project not created",
		}
	}

	// Get the project id and return the created project
	project.ID = insertResult.InsertedID.(primitive.ObjectID).Hex()
	return project, nil
}

func EditProject(conn *mongo.Client, project *projectmodels.Project) *apimodels.Error {

	// Transform team id to object id
	// also check if team id is valid
	objID, parsingError := utils.StringToObjectId(project.ID)

	if parsingError != nil {
		return &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   apierror.InvalidObjectId,
			Message: "Invalid object id [" + project.ID + "] : " + parsingError.Error(),
		}
	}

	coll := conn.Database(database.CurrentDatabase).Collection(database.PROJECT)
	updateResult, updateError := coll.UpdateByID(database.GetDefaultContext(), objID, bson.M{
		"$set": bson.M{
			"name":        project.Name,
			"description": project.Description,
			"updatedate":  project.LastUpdate,
		},
	})

	// Check if team was updated
	if updateError != nil {
		return &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   apierror.DatabaseError,
			Message: "Could not update team: " + updateError.Error(),
		}
	}

	o := objID.Hex()
	log.Default().Print(o)

	// Check if team was found
	if updateResult.MatchedCount == 0 && updateResult.ModifiedCount == 0 {
		return &apimodels.Error{
			Status:  http.StatusNotFound,
			Error:   apierror.TeamNotFound,
			Message: "Team not found",
		}
	}

	return nil
}

func DeleteProject(conn *mongo.Client, project *projectmodels.Project) *apimodels.Error {

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
			Message: "Invalid object id [" + project.ID + "] : " + parsingError.Error(),
		}
	}

	// Delete project
	projects := conn.Database(database.CurrentDatabase).Collection(database.PROJECT)
	deleteResult, err := projects.DeleteOne(database.GetDefaultContext(), bson.M{"_id": id})

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

	// Check if the project ID is valid
	projectIdObject, parsingError := utils.StringToObjectId(project.ID)
	if parsingError != nil {
		return nil, &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   apierror.InvalidObjectId,
			Message: "Invalid object id [" + project.ID + "] : " + parsingError.Error(),
		}
	}

	// Get project by id
	projects := conn.Database(database.CurrentDatabase).Collection(database.PROJECT)
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
