package tests

import (
	"testing"

	"github.com/akrck02/valhalla-core-dal/mock"
	teamdal "github.com/akrck02/valhalla-core-dal/services/team"
	teammodels "github.com/akrck02/valhalla-core-sdk/models/team"
	usersmodels "github.com/akrck02/valhalla-core-sdk/models/users"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateMockTestTeamWithOwner(conn *mongo.Client, t *testing.T, owner *usersmodels.User) *teammodels.Team {

	team := &teammodels.Team{
		Name:        mock.TeamName(),
		Description: mock.TeamDescription(),
		Owner:       owner.ID,
	}

	return CreateTestTeam(conn, t, team)
}

func CreateTestTeam(conn *mongo.Client, t *testing.T, team *teammodels.Team) *teammodels.Team {

	err := teamdal.CreateTeam(conn, team)

	if err != nil {
		t.Errorf("Error creating team: %v", err)
		return nil
	}

	t.Log("Team created successfully")
	return team
}

func GetTestTeam(conn *mongo.Client, t *testing.T, team *teammodels.Team) *teammodels.Team {

	getTeam, err := teamdal.GetTeam(conn, team)

	if err != nil {
		t.Errorf("Error getting team: %v", err)
		return nil
	}

	return getTeam
}

func EditTestTeam(conn *mongo.Client, t *testing.T, team *teammodels.Team) *teammodels.Team {

	team.Name = mock.TeamNameLong()
	team.Description = mock.TeamDescriptionLong()

	err := teamdal.EditTeam(conn, team)

	if err != nil {
		t.Errorf("Error editing team: %v", err)
		return nil
	}

	t.Log("Team edited successfully")
	return team
}

func DeleteTestTeam(conn *mongo.Client, t *testing.T, team *teammodels.Team) {

	err := teamdal.DeleteTeam(conn, team)

	if err != nil {
		t.Errorf("Error deleting team: %v", err)
		return
	}

	t.Log("Team deleted successfully")
}

func AddTestTeamMember(conn *mongo.Client, t *testing.T, team *teammodels.Team, user *usersmodels.User) *teammodels.Team {

	err := teamdal.AddMember(
		conn,
		&teamdal.MemberChangeRequest{
			Team: team.ID,
			User: user.ID,
		},
	)

	if err != nil {
		t.Errorf("Error adding team member: %v", err)
		return nil
	}

	t.Log("Team member added successfully")
	return team
}
