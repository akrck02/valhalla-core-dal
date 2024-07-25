package tests

import (
	"testing"

	"github.com/akrck02/valhalla-core-dal/mock"
	projectdal "github.com/akrck02/valhalla-core-dal/services/project"
	apierror "github.com/akrck02/valhalla-core-sdk/error"
	"github.com/akrck02/valhalla-core-sdk/log"
	projectmodels "github.com/akrck02/valhalla-core-sdk/models/project"
	usersmodels "github.com/akrck02/valhalla-core-sdk/models/users"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateMockTestProjectWithUser(conn *mongo.Client, t *testing.T, user *usersmodels.User) *projectmodels.Project {

	project := &projectmodels.Project{
		Name:        mock.ProjectName(),
		Description: mock.ProjectDescription(),
		Owner:       user.ID,
	}

	return CreateTestProjectWithUser(conn, t, project, user)
}

func CreateTestProjectWithUser(conn *mongo.Client, t *testing.T, project *projectmodels.Project, user *usersmodels.User) *projectmodels.Project {

	log.FormattedInfo("Creating project: ${0}", project.Name)
	project, err := projectdal.CreateProject(conn, project)

	if err != nil {
		t.Errorf("Error creating project: %v", err)
		return nil
	}

	t.Log("Project created successfully")
	return project
}

func CreateTestProjectWithError(conn *mongo.Client, t *testing.T, project *projectmodels.Project, status int, errorcode apierror.ApiError) {

	log.FormattedInfo("Creating project: ${0}", project.Name)
	project, err := projectdal.CreateProject(conn, project)

	if err == nil {
		t.Error("Project created successfully")
		return
	}

	if err.Status != status || err.Error != errorcode {
		t.Errorf("Error code mismatch: %v", err)
		return
	}

	log.Info("Project creation failed as expected")
	log.FormattedInfo("Error creating project: ${0}", err.Message)
}

func DeleteTestProject(conn *mongo.Client, t *testing.T, project *projectmodels.Project) {

	log.FormattedInfo("Deleting project: ${0}", project.Name)
	err := projectdal.DeleteProject(conn, project)

	if err != nil {
		t.Errorf("Error deleting project: %v", err)
		return
	}

	t.Log("Project deleted successfully")
}
