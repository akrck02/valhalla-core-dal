package tests

import (
	"testing"

	"github.com/akrck02/valhalla-core-dal/mock"
)

func TestCreateTeam(t *testing.T) {

	user := RegisterMockTestUser(t)
	_ = CreateMockTestTeamWithOwner(t, user)
	DeleteTestUser(t, user)

}

func TestEditTeam(t *testing.T) {

	user := RegisterMockTestUser(t)
	team := CreateMockTestTeamWithOwner(t, user)

	team.Name = mock.TeamNameLong()
	team.Description = mock.TeamDescriptionLong()
	team = EditTestTeam(t, team)
	obtainedTeam := GetTestTeam(t, team)

	if obtainedTeam.Name != team.Name || obtainedTeam.Description != team.Description {
		t.Errorf("Team changes not reflected in database")
		return
	}

	t.Log("Team edited successfully")
	DeleteTestUser(t, user)
}

func TestDeleteTeam(t *testing.T) {

	user := RegisterMockTestUser(t)
	team := CreateMockTestTeamWithOwner(t, user)
	DeleteTestTeam(t, team)
	DeleteTestUser(t, user)
}

func TestAddTeamMember(t *testing.T) {

	user := RegisterMockTestUser(t)
	member := RegisterMockTestUser(t)
	team := CreateMockTestTeamWithOwner(t, user)
	_ = AddTestTeamMember(t, team, member)
	DeleteTestTeam(t, team)
	DeleteTestUser(t, user)

}

func TestRemoveTeamMember(t *testing.T) {

	RegisterMockTestUser(t)

}
