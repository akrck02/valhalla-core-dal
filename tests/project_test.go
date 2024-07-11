package tests

import (
	"context"
	"testing"

	"github.com/akrck02/valhalla-core-dal/database"
	"github.com/akrck02/valhalla-core-dal/mock"
	projectdal "github.com/akrck02/valhalla-core-dal/services/project"
	"github.com/akrck02/valhalla-core-sdk/http"
	projectmodels "github.com/akrck02/valhalla-core-sdk/models/project"
	"github.com/akrck02/valhalla-core-sdk/valerror"
)

func TestCreateProject(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := RegisterMockTestUser(conn, t)
	CreateMockTestProjectWithUser(conn, t, user)
	DeleteTestUser(conn, t, user)

}

func TestCreateProjectWithoutOwner(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	project := projectmodels.Project{
		Name:        "Test Project",
		Description: "Test Description",
	}

	CreateTestProjectWithError(conn, t, &project, http.HTTP_STATUS_BAD_REQUEST, valerror.EMPTY_PROJECT_OWNER)
}

func TestCreateProjectWithoutName(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	project := &projectmodels.Project{
		Description: mock.ProjectDescription(),
		Owner:       mock.Email(),
	}

	CreateTestProjectWithError(conn, t, project, http.HTTP_STATUS_BAD_REQUEST, valerror.EMPTY_PROJECT_NAME)
}
func TestCreateProjectWithoutDescription(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	project := &projectmodels.Project{
		Name:  mock.ProjectName(),
		Owner: mock.Email(),
	}

	CreateTestProjectWithError(conn, t, project, http.HTTP_STATUS_BAD_REQUEST, valerror.EMPTY_PROJECT_DESCRIPTION)
}

func TestCreateProjectThatAlreadyExists(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := RegisterMockTestUser(conn, t)
	project := CreateMockTestProjectWithUser(conn, t, user)

	CreateTestProjectWithError(conn, t, project, http.HTTP_STATUS_CONFLICT, valerror.PROJECT_ALREADY_EXISTS)
	DeleteTestUser(conn, t, user)
}

func TestGetUserProjects(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := RegisterMockTestUser(conn, t)
	project := CreateMockTestProjectWithUser(conn, t, user)
	project2 := CreateMockTestProjectWithUser(conn, t, user)

	projects := projectdal.GetUserProjects(conn, user.Email)

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

	DeleteTestUser(conn, t, user)
}
