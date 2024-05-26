package tests

import (
	"testing"

	"github.com/akrck02/valhalla-core-dal/database"
	"github.com/akrck02/valhalla-core-dal/mock"
	projectdal "github.com/akrck02/valhalla-core-dal/services/project"
	"github.com/akrck02/valhalla-core-sdk/error"
	"github.com/akrck02/valhalla-core-sdk/http"
	"github.com/akrck02/valhalla-core-sdk/models"
)

func TestCreateProject(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	user := RegisterMockTestUser(t, client)
	CreateMockTestProjectWithUser(t, client, user)
	DeleteTestUser(t, client, user)
}

func TestCreateProjectWithoutOwner(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	project := models.Project{
		Name:        "Test Project",
		Description: "Test Description",
	}

	CreateTestProjectWithError(t, client, &project, http.HTTP_STATUS_BAD_REQUEST, error.EMPTY_PROJECT_OWNER)
}

func TestCreateProjectWithoutName(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	project := &models.Project{
		Description: mock.ProjectDescription(),
		Owner:       mock.Email(),
	}

	CreateTestProjectWithError(t, client, project, http.HTTP_STATUS_BAD_REQUEST, error.EMPTY_PROJECT_NAME)
}
func TestCreateProjectWithoutDescription(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	project := &models.Project{
		Name:  mock.ProjectName(),
		Owner: mock.Email(),
	}

	CreateTestProjectWithError(t, client, project, http.HTTP_STATUS_BAD_REQUEST, error.EMPTY_PROJECT_DESCRIPTION)
}

func TestCreateProjectThatAlreadyExists(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	user := RegisterMockTestUser(t, client)
	project := CreateMockTestProjectWithUser(t, client, user)

	CreateTestProjectWithError(t, client, project, http.HTTP_STATUS_CONFLICT, error.PROJECT_ALREADY_EXISTS)
	DeleteTestUser(t, client, user)
}

func TestGetUserProjects(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	user := RegisterMockTestUser(t, client)
	project := CreateMockTestProjectWithUser(t, client, user)
	project2 := CreateMockTestProjectWithUser(t, client, user)

	projects := projectdal.GetUserProjects(client, user.Email)

	if len(projects) == 0 {
		t.Errorf("No projects found for user: %v", user.Email)
	}

	if len(projects) != 2 {
		t.Errorf("Incorrect number of projects found for user: %v", user.Email)
	}

	if projects[0].Name != project.Name {
		t.Errorf("Incorrect project found: %v", projects[0].Name)
	}

	if projects[1].Name != project2.Name {
		t.Errorf("Incorrect project found: %v", projects[1].Name)
	}

	DeleteTestUser(t, client, user)
}
