package teamdal

import (
	"github.com/akrck02/valhalla-core-dal/database"
	"github.com/akrck02/valhalla-core-sdk/http"
	systemmodels "github.com/akrck02/valhalla-core-sdk/models/system"
	teammodels "github.com/akrck02/valhalla-core-sdk/models/team"
	usersmodels "github.com/akrck02/valhalla-core-sdk/models/users"
	"github.com/akrck02/valhalla-core-sdk/utils"
	"github.com/akrck02/valhalla-core-sdk/valerror"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MemberChangeRequest struct {
	Team string `json:"teamid"`
	User string `json:"userid"`
}

func CreateTeam(conn *mongo.Client, team *teammodels.Team) *systemmodels.Error {

	// Check if team name is empty
	if utils.IsEmpty(team.Name) {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.EMPTY_TEAM_NAME,
			Message: "Team cannot be nameless",
		}
	}

	// Check if team name is valid
	checkedName := utils.ValidateName(team.Name)

	if checkedName.Response != http.HTTP_STATUS_OK {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   checkedName.Response,
			Message: checkedName.Message,
		}
	}

	// Check if team description is empty
	if utils.IsEmpty(team.Description) {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.EMPTY_TEAM_DESCRIPTION,
			Message: "Team cannot be descriptionless",
		}
	}

	// Check if team description is valid
	checkedDescription := utils.ValidateDescription(team.Description)

	if checkedDescription.Response != http.HTTP_STATUS_OK {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   checkedDescription.Response,
			Message: checkedDescription.Message,
		}
	}

	// Check if owner exists
	err1 := userExists(conn, team.Owner)

	if err1 != nil {
		return err1
	}

	// Check if team owner is empty
	if utils.IsEmpty(team.Owner) {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.NO_OWNER,
			Message: "Team requires an owner",
		}
	}

	// Check if team already exists
	coll := conn.Database(database.CurrentDatabase).Collection(database.TEAM)
	found := teamExists(coll, team)

	if found.Name != "" {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   valerror.TEAM_ALREADY_EXISTS,
			Message: "Team already exists with name " + team.Name,
		}
	}

	// add current date to team
	creationDate := utils.GetCurrentMillis()
	team.CreationDate = &creationDate
	team.LastUpdate = team.CreationDate

	// Create team
	res, err2 := coll.InsertOne(database.GetDefaultContext(), team)

	if err2 != nil {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   valerror.TEAM_ALREADY_EXISTS,
			Message: "Team already exists",
		}
	}

	team.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func DeleteTeam(conn *mongo.Client, team *teammodels.Team) *systemmodels.Error {

	// Transform team id to object id
	// also check if team id is valid
	objID, err := utils.StringToObjectId(team.ID)

	if err != nil {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.INVALID_OBJECT_ID,
			Message: "Invalid object id",
		}
	}

	// Delete team
	coll := conn.Database(database.CurrentDatabase).Collection(database.TEAM)
	_, err = coll.DeleteOne(database.GetDefaultContext(), bson.M{"_id": objID})

	// Check if team was deleted
	if err != nil {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   valerror.TEAM_NOT_FOUND,
			Message: "Team not found",
		}
	}

	return nil
}

func EditTeam(conn *mongo.Client, team *teammodels.Team) *systemmodels.Error {

	// Transform team id to object id
	// also check if team id is valid
	objID, err := utils.StringToObjectId(team.ID)

	if err != nil {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.INVALID_OBJECT_ID,
			Message: "Invalid object id",
		}
	}

	coll := conn.Database(database.CurrentDatabase).Collection(database.TEAM)

	lastUpdate := utils.GetCurrentMillis()
	team.LastUpdate = &lastUpdate
	_, err = coll.UpdateOne(database.GetDefaultContext(), bson.M{"_id": objID}, bson.M{
		"$set": bson.M{
			"name":        team.Name,
			"description": team.Description,
			"profilepic":  team.ProfilePic,
			"updatedate":  team.LastUpdate,
		},
	})

	// Check if team was updated
	if err != nil {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.UPDATE_ERROR,
			Message: "Could not update team: " + err.Error(),
		}
	}

	return nil
}

func EditTeamOwner(conn *mongo.Client, team *teammodels.Team) *systemmodels.Error {

	// Check if team owner is empty
	if utils.IsEmpty(team.Owner) {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.NO_OWNER,
			Message: "Team requires an owner",
		}
	}

	// Transform team id to object id
	// also check if team id is valid
	objID, err1 := utils.StringToObjectId(team.ID)

	if err1 != nil {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.INVALID_OBJECT_ID,
			Message: "Invalid object id",
		}
	}

	// Check if owner exists
	err2 := userExists(conn, team.Owner)

	if err2 != nil {
		return err2
	}

	// Update owner
	coll := conn.Database(database.CurrentDatabase).Collection(database.TEAM)

	result := coll.FindOneAndUpdate(database.GetDefaultContext(), bson.M{"_id": objID}, bson.M{
		"$set": bson.M{
			"owner": team.Owner,
		},
		"$currentDate": bson.M{
			"lastupdate": true,
		},
	})

	// Check if team was updated
	err3 := result.Err()

	if err3 != nil {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.UPDATE_ERROR,
			Message: "Could not change owner",
		}
	}

	return nil
}

func AddMember(conn *mongo.Client, member *MemberChangeRequest) *systemmodels.Error {

	// Check if member is empty
	if utils.IsEmpty(member.User) {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.MEMBER_PARAMETER_EMPTY,
			Message: "Adding a member requires a member",
		}
	}

	// Check if team is empty
	if utils.IsEmpty(member.Team) {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.TEAM_PARAMETER_EMPTY,
			Message: "Adding a member requires a team",
		}
	}

	// Check if member exists
	err := userExists(conn, member.User)

	if err != nil {
		return err
	}

	// Transform team id to object id
	// also check if team id is valid
	objID, parseErr := utils.StringToObjectId(member.Team)

	if parseErr != nil {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.INVALID_OBJECT_ID,
			Message: "Invalid object id",
		}
	}

	// Check if member is already in team or is owner
	coll := conn.Database(database.CurrentDatabase).Collection(database.TEAM)

	if isUserMemberOrOwner(conn, member) {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.USER_ALREADY_MEMBER,
			Message: "User is already a member of the team",
		}
	}

	// Add member to team FIXME: Not working
	result, parseErr := coll.UpdateByID(database.GetDefaultContext(), objID, bson.M{"$push": bson.M{"members": member.User}})

	if parseErr != nil {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.UPDATE_ERROR,
			Message: "Could not add member",
		}
	}

	// Check if member was added
	if result.MatchedCount == 0 {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_NO_CHANGE,
			Error:   valerror.TEAM_NOT_FOUND,
			Message: "Team member not added",
		}
	}

	return nil
}

func RemoveMember(conn *mongo.Client, member *MemberChangeRequest) *systemmodels.Error {

	// Check if member is empty
	if utils.IsEmpty(member.User) {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.MEMBER_PARAMETER_EMPTY,
			Message: "Adding a member requires a member",
		}
	}

	// Check if team is empty
	if utils.IsEmpty(member.Team) {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.TEAM_PARAMETER_EMPTY,
			Message: "Adding a member requires a team",
		}
	}

	// Transform team id to object id
	// also check if team id is valid
	objID, parseErr := utils.StringToObjectId(member.Team)

	if parseErr != nil {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.INVALID_OBJECT_ID,
			Message: "Invalid object id",
		}
	}

	// Check if member exists
	err := userExists(conn, member.User)

	if err != nil {
		return err
	}

	// deleting team owner is not allowed
	if isOwner(conn, member) {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.CANNOT_DELETE_TEAM_OWNER,
			Message: "User is owner of the team",
		}
	}

	// Check if member is in team
	if !isMember(conn, member) {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.USER_NOT_MEMBER_OF_THE_TEAM,
			Message: "User is not a member of the team",
		}
	}

	// Remove member from team
	result, parseErr := conn.Database(database.CurrentDatabase).Collection(database.TEAM).UpdateByID(database.GetDefaultContext(), objID, bson.M{"$pull": bson.M{"members": member.User}})
	if parseErr != nil {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.UPDATE_ERROR,
			Message: "Could not remove member",
		}
	}

	// Check if member was removed
	if result.MatchedCount == 0 {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_NO_CHANGE,
			Error:   valerror.TEAM_NOT_FOUND,
			Message: "Team member not removed",
		}
	}

	return nil
}

func GetTeams(conn *mongo.Client, user *usersmodels.User) ([]*teammodels.Team, *systemmodels.Error) {

	// Get the teams that the user owns
	coll := conn.Database(database.CurrentDatabase).Collection(database.TEAM)
	teamsCursor, err := coll.Find(database.GetDefaultContext(), bson.M{"owner": user.ID})

	if err != nil {
		return nil, &systemmodels.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   valerror.UNEXPECTED_ERROR,
			Message: "Cannot find teams",
		}
	}

	var teams []*teammodels.Team

	for teamsCursor.Next(database.GetDefaultContext()) {
		var team teammodels.Team
		err := teamsCursor.Decode(&team)
		if err != nil {
			return nil, &systemmodels.Error{
				Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
				Error:   valerror.UNEXPECTED_ERROR,
				Message: "Cannot find teams",
			}
		}
		teams = append(teams, &team)
	}

	// TODO: get the teams the user is member of

	return teams, nil
}

func GetTeam(conn *mongo.Client, team *teammodels.Team) (*teammodels.Team, *systemmodels.Error) {

	objID, err1 := utils.StringToObjectId(team.ID)

	if err1 != nil {
		return nil, &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.INVALID_OBJECT_ID,
			Message: "Invalid object id",
		}
	}

	coll := conn.Database(database.CurrentDatabase).Collection(database.TEAM)
	var foundTeam teammodels.Team

	err2 := coll.FindOne(database.GetDefaultContext(), bson.M{"_id": objID}).Decode(&foundTeam)

	if err2 != nil {
		return nil, &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.TEAM_NOT_FOUND,
			Message: "Team not found",
		}
	}

	return &foundTeam, nil
}

func SearchTeams(searchText *string) (*[]teammodels.Team, *systemmodels.Error) {

	foundTeams := []teammodels.Team{}
	return &foundTeams, nil

}

func userExists(conn *mongo.Client, user string) *systemmodels.Error {

	coll := conn.Database(database.CurrentDatabase).Collection(database.USER)
	var foundUser usersmodels.User

	objID, err1 := utils.StringToObjectId(user)

	if err1 != nil {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.INVALID_OBJECT_ID,
			Message: "Invalid object id",
		}
	}

	err2 := coll.FindOne(database.GetDefaultContext(), bson.M{"_id": objID}).Decode(&foundUser)

	if err2 != nil {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.OWNER_DOESNT_EXIST,
			Message: "User doesn't exists",
		}
	}

	return nil
}

func teamExists(coll *mongo.Collection, team *teammodels.Team) teammodels.Team {

	filter := bson.D{
		{Key: "name", Value: team.Name},
		{Key: "owner", Value: team.Owner},
	}
	var result teammodels.Team

	coll.FindOne(database.GetDefaultContext(), filter).Decode(&result)

	return result
}

func isUserMemberOrOwner(conn *mongo.Client, request *MemberChangeRequest) bool {

	filterMember := bson.D{
		{Key: "_id", Value: request.Team},
		{Key: "members", Value: bson.D{{Key: "$all", Value: bson.A{request.User}}}},
	}

	filterOwner := bson.D{
		{Key: "_id", Value: request.Team},
		{Key: "owner", Value: request.User},
	}

	coll := conn.Database(database.CurrentDatabase).Collection(database.TEAM)

	var result teammodels.Team

	err := coll.FindOne(database.GetDefaultContext(), filterMember).Decode(&result)

	if err == nil && result.ID != "" {
		return true
	}

	err = coll.FindOne(database.GetDefaultContext(), filterOwner).Decode(&result)

	if err == nil && result.ID != "" {
		return true
	}

	return false
}

func isMember(conn *mongo.Client, request *MemberChangeRequest) bool {

	filter := bson.D{
		{Key: "_id", Value: request.Team},
		{Key: "members", Value: bson.D{{Key: "$all", Value: bson.A{request.User}}}},
	}

	coll := conn.Database(database.CurrentDatabase).Collection(database.TEAM)

	var result teammodels.Team
	err := coll.FindOne(database.GetDefaultContext(), filter).Decode(&result)
	return err != nil
}

func isOwner(conn *mongo.Client, request *MemberChangeRequest) bool {

	filter := bson.D{
		{Key: "_id", Value: request.Team},
		{Key: "owner", Value: request.User},
	}

	coll := conn.Database(database.CurrentDatabase).Collection(database.TEAM)

	var result teammodels.Team
	err := coll.FindOne(database.GetDefaultContext(), filter).Decode(&result)
	return err != nil
}
