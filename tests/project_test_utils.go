package tests

import (
	"context"
	"testing"

	"github.com/akrck02/valhalla-core-dal/mock"
	projectdal "github.com/akrck02/valhalla-core-dal/services/project"
	"github.com/akrck02/valhalla-core-sdk/log"
	"github.com/akrck02/valhalla-core-sdk/models"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateMockTestProject creates a project for testing purposes
//
// [param] t | *testing.T: Testing object
// [param] conn | context.Context: Database connection
// [param] client | *mongo.Client: Database client
//
// [return] models.Project: Created project
func CreateMockTestProjectWithUser(t *testing.T, conn context.Context, client *mongo.Client, user *models.User) *models.Project {

	project := &models.Project{
		Name:        mock.ProjectName(),
		Description: mock.ProjectDescription(),
		Owner:       user.Email,
	}

	return CreateTestProjectWithUser(t, conn, client, project, user)
}

// CreateTestProject creates a project for testing purposes
//
// [param] t | *testing.T: Testing object
// [param] conn | context.Context: Database connection
// [param] client | *mongo.Client: Database client
// [param] project | *models.Project: Project to create
//
// [return] *models.Project: Created project
func CreateTestProjectWithUser(t *testing.T, conn context.Context, client *mongo.Client, project *models.Project, user *models.User) *models.Project {

	log.FormattedInfo("Creating project: ${0}", project.Name)
	err := projectdal.CreateProject(conn, client, project)

	if err != nil {
		t.Errorf("Error creating project: %v", err)
		return nil
	}

	t.Log("Project created successfully")
	return project
}

// CreateTestProjectWithoutOwner creates a project without an owner for testing purposes
//
// [param] t | *testing.T: Testing object
// [param] conn | context.Context: Database connection
// [param] client | *mongo.Client: Database client
func CreateTestProjectWithError(t *testing.T, conn context.Context, client *mongo.Client, project *models.Project, status int, errorcode int) {

	log.FormattedInfo("Creating project: ${0}", project.Name)
	err := projectdal.CreateProject(conn, client, project)

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

// DeleteTestProject deletes a project for testing purposes
//
// [param] t | *testing.T: Testing object
// [param] conn | context.Context: Database connection
// [param] client | *mongo.Client: Database client
// [param] project | *models.Project: Project to delete
func DeleteTestProject(t *testing.T, conn context.Context, client *mongo.Client, project *models.Project) {

	log.FormattedInfo("Deleting project: ${0}", project.Name)
	err := projectdal.DeleteProject(conn, client, project)

	if err != nil {
		t.Errorf("Error deleting project: %v", err)
		return
	}

	t.Log("Project deleted successfully")
}
