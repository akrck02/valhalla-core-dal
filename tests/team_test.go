package tests

import (
	"context"
	"testing"

	"github.com/akrck02/valhalla-core-dal/database"
	"github.com/akrck02/valhalla-core-dal/mock"
)

func TestCreateTeam(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := RegisterMockTestUser(conn, t)
	_ = CreateMockTestTeamWithOwner(conn, t, user)
	DeleteTestUser(conn, t, user)

}

func TestEditTeam(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := RegisterMockTestUser(conn, t)
	team := CreateMockTestTeamWithOwner(conn, t, user)

	team.Name = mock.TeamNameLong()
	team.Description = mock.TeamDescriptionLong()
	team = EditTestTeam(conn, t, team)
	obtainedTeam := GetTestTeam(conn, t, team)

	if obtainedTeam.Name != team.Name || obtainedTeam.Description != team.Description {
		t.Errorf("Team changes not reflected in database")
		return
	}

	t.Log("Team edited successfully")
	DeleteTestUser(conn, t, user)
}

func TestDeleteTeam(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := RegisterMockTestUser(conn, t)
	team := CreateMockTestTeamWithOwner(conn, t, user)
	DeleteTestTeam(conn, t, team)
	DeleteTestUser(conn, t, user)
}

func TestAddTeamMember(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := RegisterMockTestUser(conn, t)
	member := RegisterMockTestUser(conn, t)
	team := CreateMockTestTeamWithOwner(conn, t, user)
	_ = AddTestTeamMember(conn, t, team, member)
	DeleteTestTeam(conn, t, team)
	DeleteTestUser(conn, t, user)
	DeleteTestUser(conn, t, member)

}

func TestRemoveTeamMember(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	RegisterMockTestUser(conn, t)

}
