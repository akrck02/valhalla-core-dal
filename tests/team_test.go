package tests

import (
	"testing"

	"github.com/akrck02/valhalla-core-dal/mock"
	teamdal "github.com/akrck02/valhalla-core-dal/services/team"
	"github.com/akrck02/valhalla-core-sdk/models"
)

func TestCreateTeam(t *testing.T) {

	user := RegisterMockTestUser(t)

	err := teamdal.CreateTeam(&models.Team{
		Name:        mock.Name(),
		Description: mock.Description(),
		Owner:       user.ID,
	})

	if err != nil {
		t.Error("error on create", err)
	}

	DeleteTestUser(t, user)

}

func TestEditTeam(t *testing.T) {

	RegisterMockTestUser(t)

}

func TestDeleteTeam(t *testing.T) {

	RegisterMockTestUser(t)

}

func TestAddTeamMember(t *testing.T) {

	RegisterMockTestUser(t)
}

func TestRemoveTeamMember(t *testing.T) {

	RegisterMockTestUser(t)

}
