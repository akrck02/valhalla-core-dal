package tests

import (
	"context"
	"net/http"
	"testing"

	"github.com/akrck02/valhalla-core-dal/database"
	"github.com/akrck02/valhalla-core-dal/mock"
	projectdal "github.com/akrck02/valhalla-core-dal/services/project"

	apierror "github.com/akrck02/valhalla-core-sdk/error"
	projectmodels "github.com/akrck02/valhalla-core-sdk/models/project"
)

func TestCreateProject(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := RegisterMockTestUser(conn, t)
	project := CreateMockTestProjectWithUser(conn, t, user)

	if project == nil {
		t.Errorf("Error creating project")
		return
	}

	if project.CreationDate == nil || project.LastUpdate == nil {
		t.Errorf("Error creating project, dates not set")
		return
	}

}

func TestCreateProjectWithoutOwner(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	project := projectmodels.Project{
		Name:        "Test Project",
		Description: "Test Description",
	}

	CreateTestProjectWithError(conn, t, &project, http.StatusBadRequest, apierror.EmptyProjectOwner)
}

func TestCreateProjectWithoutName(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	project := &projectmodels.Project{
		Description: mock.ProjectDescription(),
		Owner:       mock.Email(),
	}

	CreateTestProjectWithError(conn, t, project, http.StatusBadRequest, apierror.EmptyProjectName)
}
func TestCreateProjectWithoutDescription(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	project := &projectmodels.Project{
		Name:  mock.ProjectName(),
		Owner: mock.Email(),
	}

	CreateTestProjectWithError(conn, t, project, http.StatusBadRequest, apierror.EmptyProjectDescription)
}

func TestCreateProjectThatAlreadyExists(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := RegisterMockTestUser(conn, t)
	project := CreateMockTestProjectWithUser(conn, t, user)

	CreateTestProjectWithError(conn, t, project, http.StatusConflict, apierror.ProjectAlreadyExists)
	DeleteTestUser(conn, t, user)
}

func TestGetUserProjects(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := RegisterMockTestUser(conn, t)
	project := CreateMockTestProjectWithUser(conn, t, user)
	project2 := CreateMockTestProjectWithUser(conn, t, user)

	projects := projectdal.GetUserProjects(conn, user.ID)

	if len(projects) == 0 {
		t.Errorf("No projects found for user: %v", user.ID)
		return
	}

	if len(projects) != 2 {
		t.Errorf("Incorrect number of projects found for user: %v", user.ID)
		return
	}

	if projects[0].Name != project.Name {
		t.Errorf("Incorrect project found: %v", projects[0].Name)
		return
	}

	if projects[1].Name != project2.Name {
		t.Errorf("Incorrect project found: %v", projects[1].Name)
		return
	}

	DeleteTestUser(conn, t, user)
}

func TestEditProject(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := RegisterMockTestUser(conn, t)
	project := CreateMockTestProjectWithUser(conn, t, user)

	if project == nil {
		t.Errorf("Error creating project")
		return
	}

	if project.CreationDate == nil || project.LastUpdate == nil {
		t.Errorf("Error creating project, dates not set")
		return
	}

}
