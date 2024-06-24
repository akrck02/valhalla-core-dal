package tests

import (
	"testing"

	"github.com/akrck02/valhalla-core-dal/mock"
	projectdal "github.com/akrck02/valhalla-core-dal/services/project"
	"github.com/akrck02/valhalla-core-sdk/http"
	"github.com/akrck02/valhalla-core-sdk/models"
	"github.com/akrck02/valhalla-core-sdk/valerror"
)

func TestCreateProject(t *testing.T) {

	user := RegisterMockTestUser(t)
	CreateMockTestProjectWithUser(t, user)
	DeleteTestUser(t, user)
}

func TestCreateProjectWithoutOwner(t *testing.T) {

	project := models.Project{
		Name:        "Test Project",
		Description: "Test Description",
	}

	CreateTestProjectWithError(t, &project, http.HTTP_STATUS_BAD_REQUEST, valerror.EMPTY_PROJECT_OWNER)
}

func TestCreateProjectWithoutName(t *testing.T) {

	project := &models.Project{
		Description: mock.ProjectDescription(),
		Owner:       mock.Email(),
	}

	CreateTestProjectWithError(t, project, http.HTTP_STATUS_BAD_REQUEST, valerror.EMPTY_PROJECT_NAME)
}
func TestCreateProjectWithoutDescription(t *testing.T) {

	project := &models.Project{
		Name:  mock.ProjectName(),
		Owner: mock.Email(),
	}

	CreateTestProjectWithError(t, project, http.HTTP_STATUS_BAD_REQUEST, valerror.EMPTY_PROJECT_DESCRIPTION)
}

func TestCreateProjectThatAlreadyExists(t *testing.T) {

	user := RegisterMockTestUser(t)
	project := CreateMockTestProjectWithUser(t, user)

	CreateTestProjectWithError(t, project, http.HTTP_STATUS_CONFLICT, valerror.PROJECT_ALREADY_EXISTS)
	DeleteTestUser(t, user)
}

func TestGetUserProjects(t *testing.T) {

	user := RegisterMockTestUser(t)
	project := CreateMockTestProjectWithUser(t, user)
	project2 := CreateMockTestProjectWithUser(t, user)

	projects := projectdal.GetUserProjects(user.Email)

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

	DeleteTestUser(t, user)
}
