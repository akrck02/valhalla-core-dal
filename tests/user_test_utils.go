package tests

import (
	"testing"

	"github.com/akrck02/valhalla-core-dal/mock"
	userdal "github.com/akrck02/valhalla-core-dal/services/user"
	apierror "github.com/akrck02/valhalla-core-sdk/error"
	"github.com/akrck02/valhalla-core-sdk/log"
	usersmodels "github.com/akrck02/valhalla-core-sdk/models/users"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterMockTestUser(conn *mongo.Client, t *testing.T) *usersmodels.User {

	var user = &usersmodels.User{
		Email:    mock.Email(),
		Password: mock.Password(),
		Username: mock.Username(),
	}

	return RegisterTestUser(conn, t, user)

}

func RegisterTestUser(conn *mongo.Client, t *testing.T, user *usersmodels.User) *usersmodels.User {

	log.FormattedInfo("Registering user: ${0}", user.Email)
	log.FormattedInfo("Password: ${0}", user.Password)
	log.FormattedInfo("Username: ${0}", user.Username)

	err := userdal.Register(conn, user)

	if err != nil {
		t.Error("The user was not registered", err)
		return nil
	}

	return user
}

func RegisterTestUserWithError(conn *mongo.Client, t *testing.T, user *usersmodels.User, status int, errorcode apierror.ApiError) {

	log.FormattedInfo("Registering user: ${0}", user.Email)
	log.FormattedInfo("Password: ${0}", user.Password)
	log.FormattedInfo("Username: ${0}", user.Username)

	err := userdal.Register(conn, user)

	if err == nil {
		t.Error("The user was registered.")
		return
	}

	if err.Status != status || err.Error != errorcode {
		t.Error("The error is not the expected [", status, " / ", errorcode, "], current: [", err.Status, " / ", err.Error, "]", err.Message)
		return
	}

	log.Info("User not registered")
	log.FormattedInfo("Error: ${0}", err.Message)

}

func LoginTestUser(conn *mongo.Client, t *testing.T, user *usersmodels.User, ip string, userAgent string) string {

	log.FormattedInfo("Logging in user: ${0}", user.Email)
	log.FormattedInfo("Password: ${0}", user.Password)

	token, err := userdal.Login(conn, user, ip, userAgent)

	if err != nil {
		t.Error("The user was not logged in", err)
		return ""
	}

	if token == "" {
		t.Error("The token is empty")
		return ""
	}

	_, err = userdal.IsTokenValid(conn, token)

	if err != nil {
		t.Error("The token was not validated", err)
		return ""
	}

	log.Info("User logged in")
	log.FormattedInfo("Token: ${0}", token)
	return token

}

func LoginTestUserWithError(conn *mongo.Client, t *testing.T, user *usersmodels.User, ip string, userAgent string, status int, errorcode apierror.ApiError) {

	log.FormattedInfo("Logging in user: ${0}", user.Email)
	log.FormattedInfo("Password: ${0}", user.Password)

	_, err := userdal.Login(conn, user, ip, userAgent)

	if err == nil {
		t.Error("The user was logged in.")
		return
	}

	if err.Status != status || err.Error != errorcode {
		t.Error("The error is not the expected [", status, " / ", errorcode, "], current: [", err.Status, " / ", err.Error, "]", err.Message)
		return
	}

	log.Info("User not logged in")
	log.FormattedInfo("Error: ${0}", err.Message)
}

func LoginAuthTestUser(conn *mongo.Client, t *testing.T, email string, token string) {

	log.FormattedInfo("Authenticating token: ${0}", token)
	log.Info("Checking if the token is valid")

	var authLogin = &usersmodels.AuthLogin{
		Email:     email,
		AuthToken: token,
	}

	err := userdal.LoginAuth(conn, authLogin, mock.Ip(), mock.Platform())

	if err != nil {
		t.Error("The token is not valid", err)
		return
	}

	log.Info("Token is valid")

}

func DeleteTestUser(conn *mongo.Client, t *testing.T, user *usersmodels.User) {

	log.FormattedInfo("Deleting user: ${0}", user.Email)

	err := userdal.DeleteUser(conn, user)

	if err != nil {
		t.Error("The user was not deleted", err)
		return
	}

	log.Info("User deleted")

}

func DeleteTestUserWithError(conn *mongo.Client, t *testing.T, user *usersmodels.User, status int, errorcode apierror.ApiError) {

	log.FormattedInfo("Deleting user: ${0}", user.Email)

	err := userdal.DeleteUser(conn, user)

	if err == nil {
		t.Error("The user was deleted.")
		return
	}

	if err.Status != status || err.Error != errorcode {
		t.Error("The error is not the expected [", status, " / ", errorcode, "], current: [", err.Status, " / ", err.Error, "]", err.Message)
		return
	}

	log.Info("User not deleted")
	log.FormattedInfo("Error: ${0}", err.Message)

}

func EditTestUserEmail(conn *mongo.Client, t *testing.T, emailChangeRequest *userdal.EmailChangeRequest) {

	log.FormattedInfo("Editing user mail: ${0}", emailChangeRequest.Email)
	log.FormattedInfo("New email: ${0}", emailChangeRequest.NewEmail)

	err := userdal.EditUserEmail(conn, emailChangeRequest)

	if err != nil {
		t.Error("The user was not edited", err)
		return
	}

	log.Info("User edited")

}

func EditTestUserEmailWithError(conn *mongo.Client, t *testing.T, emailChangeRequest *userdal.EmailChangeRequest, status int, errorcode apierror.ApiError) {

	log.FormattedInfo("Editing user mail: ${0}", emailChangeRequest.Email)
	log.FormattedInfo("New email: ${0}", emailChangeRequest.NewEmail)

	err := userdal.EditUserEmail(conn, emailChangeRequest)

	if err == nil {
		t.Error("The user was edited.")
		return
	}

	if err.Status != status || err.Error != errorcode {
		t.Error("The error is not the expected [", status, " / ", errorcode, "], current: [", err.Status, " / ", err.Error, "]", err.Message)
		return
	}

	log.Info("User not edited")
	log.FormattedInfo("Error: ${0}", err.Message)

}

func EditTestUser(conn *mongo.Client, t *testing.T, user *usersmodels.User) {

	log.FormattedInfo("Editing user: ${0}", user.Email)
	err := userdal.EditUser(conn, user)

	if err != nil {
		t.Error("The user was not edited", err)
		return
	}

	log.Info("User edited")

}

func EditTestUserWithError(conn *mongo.Client, t *testing.T, user *usersmodels.User, status int, errorcode apierror.ApiError) {

	log.FormattedInfo("Editing user: ${0}", user.Email)
	err := userdal.EditUser(conn, user)

	if err == nil {
		t.Error("The user was edited.")
		return
	}

	if err.Status != status || err.Error != errorcode {
		t.Error("The error is not the expected [", status, " / ", errorcode, "], current: [", err.Status, " / ", err.Error, "]", err.Message)
		return
	}

	log.Info("User not edited")
	log.FormattedInfo("Error: ${0}", err.Message)

}

func ValidateTestTokenWithError(conn *mongo.Client, t *testing.T, token string, status int, errorcode apierror.ApiError) {

	log.FormattedInfo("Validating token: ${0}", token)

	_, err := userdal.IsTokenValid(conn, token)

	if err == nil {
		t.Error("The token was validated.")
		return
	}

	if err.Status != status || err.Error != errorcode {
		t.Error("The error is not the expected [", status, " / ", errorcode, "], current: [", err.Status, " / ", err.Error, "]", err.Message)
		return
	}

	log.Info("Token not validated")
	log.FormattedInfo("Token not validated, error: ${0}", err.Message)

}
