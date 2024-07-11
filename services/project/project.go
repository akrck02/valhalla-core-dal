package projectdal

import (
	"github.com/akrck02/valhalla-core-dal/database"
	"github.com/akrck02/valhalla-core-sdk/http"
	projectmodels "github.com/akrck02/valhalla-core-sdk/models/project"
	systemmodels "github.com/akrck02/valhalla-core-sdk/models/system"
	"github.com/akrck02/valhalla-core-sdk/utils"
	"github.com/akrck02/valhalla-core-sdk/valerror"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateProject(conn *mongo.Client, project *projectmodels.Project) *systemmodels.Error {

	if utils.IsEmpty(project.Name) {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.EMPTY_PROJECT_NAME,
			Message: "Project name cannot be empty",
		}
	}

	if utils.IsEmpty(project.Description) {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.EMPTY_PROJECT_DESCRIPTION,
			Message: "Project description cannot be empty",
		}
	}

	if utils.IsEmpty(project.Owner) {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.EMPTY_PROJECT_OWNER,
			Message: "Owner cannot be empty",
		}
	}

	coll := conn.Database(database.CurrentDatabase).Collection(database.PROJECT)
	found := nameExists(coll, project.Name, project.Owner)

	if found {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_CONFLICT,
			Error:   valerror.PROJECT_ALREADY_EXISTS,
			Message: "Project already exists",
		}
	}

	_, err := coll.InsertOne(database.GetDefaultContext(), project)

	if err != nil {
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

	// delete user devices
	devices := client.Database(database.CurrentDatabase).Collection(database.DEVICE)
	_, err := devices.DeleteMany(database.GetDefaultContext(), bson.M{"project": project.Name})

	if err != nil {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   valerror.PROJECT_NOT_DELETED,
			Message: "Project not deleted",
		}
	}

	projects := client.Database(database.CurrentDatabase).Collection(database.PROJECT)

	var deleteResult *mongo.DeleteResult
	deleteResult, err = projects.DeleteOne(database.GetDefaultContext(), bson.M{"name": project.Name})

	if err != nil {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   valerror.PROJECT_NOT_DELETED,
			Message: "Project not deleted",
		}
	}

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

	projects := client.Database(database.CurrentDatabase).Collection(database.PROJECT)

	found := projectmodels.Project{}
	err := projects.FindOne(database.GetDefaultContext(), bson.M{"name": project.Name}).Decode(&found)

	if err != nil {
		return nil, &systemmodels.Error{
			Status:  http.HTTP_STATUS_NOT_FOUND,
			Error:   valerror.PROJECT_NOT_FOUND,
			Message: "Project not found",
		}
	}

	return &found, nil
}

func GetUserProjects(conn *mongo.Client, email string) []projectmodels.Project {

	projects := conn.Database(database.CurrentDatabase).Collection(database.PROJECT)
	filter := bson.M{"owner": email}
	cursor, err := projects.Find(database.GetDefaultContext(), filter)
	if err != nil {
		return nil
	}

	var result []projectmodels.Project
	cursor.All(database.GetDefaultContext(), &result)

	return result
}

func nameExists(coll *mongo.Collection, name string, owner string) bool {
	filter := bson.M{"name": name, "owner": owner}

	var result *projectmodels.Project
	coll.FindOne(database.GetDefaultContext(), filter).Decode(&result)

	return result != nil
}
