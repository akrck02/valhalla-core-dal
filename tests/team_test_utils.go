package tests

import (
	"testing"

	"github.com/akrck02/valhalla-core-dal/mock"
	teamdal "github.com/akrck02/valhalla-core-dal/services/team"
	"github.com/akrck02/valhalla-core-sdk/models"
)

// CreateMockTestTeam creates a team for testing purposes
// with the given user as owner
//
// [param] t | *testing.T: Testing object
//
// [return] models.Team: Created team
func CreateMockTestTeamWithOwner(t *testing.T, owner *models.User) *models.Team {

	team := &models.Team{
		Name:        mock.TeamName(),
		Description: mock.TeamDescription(),
		Owner:       owner.ID,
	}

	return CreateTestTeam(t, team)
}

// Creates a team for testing purposes
//
// [param] t | *testing.T: Testing object
// [param] team | *models.Team: Team to create
//
// [return] *models.Team: Created team
func CreateTestTeam(t *testing.T, team *models.Team) *models.Team {

	err := teamdal.CreateTeam(team)

	if err != nil {
		t.Errorf("Error creating team: %v", err)
		return nil
	}

	t.Log("Team created successfully")
	return team
}

// Gets a team for testing purposes
//
// [param] t | *testing.T: Testing object
// [param] team | *models.Team: Team to get
//
// [return] *models.Team: Got team
func GetTestTeam(t *testing.T, team *models.Team) *models.Team {

	getTeam, err := teamdal.GetTeam(team)

	if err != nil {
		t.Errorf("Error getting team: %v", err)
		return nil
	}

	return getTeam
}

// Edits a team for testing purposes
//
// [param] t | *testing.T: Testing object
// [param] team | *models.Team: Team to edit
//
// [return] team | *models.Team: Edited team
func EditTestTeam(t *testing.T, team *models.Team) *models.Team {

	team.Name = mock.TeamNameLong()
	team.Description = mock.TeamDescriptionLong()

	err := teamdal.EditTeam(team)

	if err != nil {
		t.Errorf("Error editing team: %v", err)
		return nil
	}

	t.Log("Team edited successfully")
	return team
}

// Deletes a team for testing purposes
//
// [param] t | *testing.T: Testing object
// [param] team | *models.Team: Team to delete
func DeleteTestTeam(t *testing.T, team *models.Team) {

	err := teamdal.DeleteTeam(team)

	if err != nil {
		t.Errorf("Error deleting team: %v", err)
		return
	}

	t.Log("Team deleted successfully")
}

// Adds a team member for testing purposes
//
// [param] t | *testing.T: Testing object
// [param] team | *models.Team: Team to add member to
// [param] user | *models.User: User to add to team
//
// [return] *models.Team: Team with added member
func AddTestTeamMember(t *testing.T, team *models.Team, user *models.User) *models.Team {

	err := teamdal.AddMember(
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
