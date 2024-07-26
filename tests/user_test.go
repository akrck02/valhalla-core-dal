package tests

import (
	"context"
	"net/http"
	"testing"

	"github.com/akrck02/valhalla-core-dal/database"
	userdal "github.com/akrck02/valhalla-core-dal/services/user"
	apierror "github.com/akrck02/valhalla-core-sdk/error"

	"github.com/akrck02/valhalla-core-sdk/log"
	usersmodels "github.com/akrck02/valhalla-core-sdk/models/users"

	"github.com/akrck02/valhalla-core-dal/mock"
)

func TestRegister(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := RegisterMockTestUser(conn, t)

	if user == nil {
		t.Error("The user was not registered")
		return
	}

	if user.CreationDate == nil {
		t.Error("The user creation date is invalid")
		return
	}

	if user.LastUpdate == nil {
		t.Error("The user last update date is invalid")
		return
	}
}

func TestRegisterNotEmail(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := &usersmodels.User{
		Password: mock.Password(),
		Username: mock.Username(),
	}

	RegisterTestUserWithError(conn, t, user, http.StatusBadRequest, apierror.EmptyUserEmail)
}

func TestRegisterNotUsername(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := &usersmodels.User{
		Email:    mock.Email(),
		Password: mock.Password(),
	}

	RegisterTestUserWithError(conn, t, user, http.StatusBadRequest, apierror.EmptyUsername)
}

func TestRegisterNotPassword(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := &usersmodels.User{
		Email:    mock.Email(),
		Username: mock.Username(),
	}

	RegisterTestUserWithError(conn, t, user, http.StatusBadRequest, apierror.EmptyUserPassword)
}

func TestRegisterNotDotEmail(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := &usersmodels.User{
		Email:    mock.EmailNotDot(),
		Password: mock.Password(),
		Username: mock.Username(),
	}

	RegisterTestUserWithError(conn, t, user, http.StatusBadRequest, apierror.NoDotEmail)
}

func TestRegisterNotAtEmail(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := &usersmodels.User{
		Email:    mock.EmailNotAt(),
		Password: mock.Password(),
		Username: mock.Username(),
	}

	RegisterTestUserWithError(conn, t, user, http.StatusBadRequest, apierror.NoAtEmail)
}

func TestRegisterShortMail(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := &usersmodels.User{
		Email:    mock.EmailShort(),
		Password: mock.Password(),
		Username: mock.Username(),
	}

	RegisterTestUserWithError(conn, t, user, http.StatusBadRequest, apierror.ShortEmail)
}

func TestRegisterNotSpecialCharactersPassword(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := &usersmodels.User{
		Email:    mock.Email(),
		Password: mock.PasswordNotSpecialChar(),
		Username: mock.Username(),
	}

	RegisterTestUserWithError(conn, t, user, http.StatusBadRequest, apierror.NoSpecialCharactersPassword)
}

func TestRegisterNotUpperCaseLowerCasePassword(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := &usersmodels.User{
		Email:    mock.Email(),
		Password: mock.PasswordNotLowerCase(),
		Username: mock.Username(),
	}

	RegisterTestUserWithError(conn, t, user, http.StatusBadRequest, apierror.NoMayusMinusPassword)

	user = &usersmodels.User{
		Email:    mock.Email(),
		Password: mock.PasswordNotUpperCase(),
		Username: mock.Username(),
	}

	RegisterTestUserWithError(conn, t, user, http.StatusBadRequest, apierror.NoMayusMinusPassword)
}

func TestRegisterShortPassword(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := &usersmodels.User{
		Email:    mock.Email(),
		Password: mock.PasswordShort(),
		Username: mock.Username(),
	}

	RegisterTestUserWithError(conn, t, user, http.StatusBadRequest, apierror.ShortPassword)
}

func TestRegisterNotNumbersPassword(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := &usersmodels.User{
		Email:    mock.Email(),
		Password: mock.PasswordNotNumber(),
		Username: mock.Username(),
	}

	RegisterTestUserWithError(conn, t, user, http.StatusBadRequest, apierror.NoAlphanumericPassword)
}

func TestLogin(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := RegisterMockTestUser(conn, t)
	LoginTestUser(conn, t, user, mock.Ip(), mock.Platform())

}

func TestLoginWrongPassword(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := RegisterMockTestUser(conn, t)
	user.Password = mock.PasswordShort()

	log.Info("Login with wrong password")
	log.FormattedInfo("Password: ${0}", user.Password)
	LoginTestUserWithError(conn, t, user, mock.Ip(), mock.Platform(), http.StatusForbidden, apierror.UserNotAuthorized)

}

func TestLoginWrongEmail(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := RegisterMockTestUser(conn, t)
	user.Email = "wrong" + mock.Email()
	LoginTestUserWithError(conn, t, user, mock.Ip(), mock.Platform(), http.StatusForbidden, apierror.UserNotAuthorized)
}

func TestLoginAuth(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := RegisterMockTestUser(conn, t)
	token := LoginTestUser(conn, t, user, mock.Ip(), mock.Platform())
	LoginAuthTestUser(conn, t, user.Email, token)
}

func TestDeleteUser(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := RegisterMockTestUser(conn, t)
	LoginTestUser(conn, t, user, mock.Ip(), mock.Platform())
}

func TestDeleteUserNoEmail(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := &usersmodels.User{
		Password: mock.Password(),
		Username: mock.Username(),
	}

	DeleteTestUserWithError(conn, t, user, http.StatusBadRequest, apierror.EmptyUserEmail)
}

func TestDeleteUserNotFound(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := &usersmodels.User{
		Email:    mock.Email(),
		Password: mock.Password(),
		Username: mock.Username(),
	}

	DeleteTestUserWithError(conn, t, user, http.StatusNotFound, apierror.UserNotFound)
}

func TestEditUserEmail(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	email := mock.Email()
	newEmail := "xXx" + mock.Email()

	user := &usersmodels.User{
		Email:    email,
		Password: mock.Password(),
		Username: mock.Username(),
	}

	newUser := RegisterTestUser(conn, t, user)
	newUser.Password = user.Password

	LoginTestUser(conn, t, user, mock.Ip(), mock.Platform())

	// Change the user email
	log.Info("Changing user email")
	emailChangeRequest := userdal.EmailChangeRequest{
		Email:    email,
		NewEmail: newEmail,
	}

	EditTestUserEmail(conn, t, &emailChangeRequest)
	newUser.Email = newEmail

	newUser, err := userdal.GetUser(conn, newUser, true)

	if err != nil {
		t.Error("The user was not found", err)
		return
	}

	if newUser.LastUpdate == user.LastUpdate {
		t.Error("The user last update date was not changed")
		return
	}
}

func TestEditUserEmailNoEmail(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	emailChangeRequest := userdal.EmailChangeRequest{}
	EditTestUserEmailWithError(conn, t, &emailChangeRequest, http.StatusBadRequest, apierror.EmptyUserEmail)

}

func TestEditUserEmailNoDotEmail(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	email := mock.Email()
	newEmail := mock.EmailNotDot()

	emailChangeRequest := userdal.EmailChangeRequest{
		Email:    email,
		NewEmail: newEmail,
	}

	EditTestUserEmailWithError(conn, t, &emailChangeRequest, http.StatusBadRequest, apierror.NoDotEmail)

}

func TestEditUserEmailNoAtEmail(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	email := mock.Email()
	newEmail := mock.EmailNotAt()

	emailChangeRequest := userdal.EmailChangeRequest{
		Email:    email,
		NewEmail: newEmail,
	}

	EditTestUserEmailWithError(conn, t, &emailChangeRequest, http.StatusBadRequest, apierror.NoAtEmail)

}

func TestEditUserEmailShortEmail(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	email := mock.Email()
	newEmail := mock.EmailShort()

	emailChangeRequest := userdal.EmailChangeRequest{
		Email:    email,
		NewEmail: newEmail,
	}

	EditTestUserEmailWithError(conn, t, &emailChangeRequest, http.StatusBadRequest, apierror.ShortEmail)

}

func TestEditUserEmailNotFound(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	email := mock.Email()
	newEmail := "xXx" + mock.Email()

	emailChangeRequest := userdal.EmailChangeRequest{
		Email:    email,
		NewEmail: newEmail,
	}

	EditTestUserEmailWithError(conn, t, &emailChangeRequest, http.StatusNotFound, apierror.UserNotFound)
}

func TestEditUserEmailExists(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	email := mock.Email()
	newEmail := mock.Email() + "xXx"

	// Create a user
	user := &usersmodels.User{
		Email:    email,
		Password: mock.Password(),
		Username: mock.Username(),
	}

	RegisterTestUser(conn, t, user)
	log.Jump()

	// Create a new user with the new email
	newUser := &usersmodels.User{
		Email:    newEmail,
		Password: mock.Password(),
		Username: mock.Username(),
	}

	RegisterTestUser(conn, t, newUser)
	log.Jump()

	// Change the email
	emailChangeRequest := userdal.EmailChangeRequest{
		Email:    email,
		NewEmail: newEmail,
	}

	EditTestUserEmailWithError(conn, t, &emailChangeRequest, http.StatusConflict, apierror.UserAlreadyExists)
}

func TestEditUserSameEmail(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	email := mock.Email()
	emailChangeRequest := userdal.EmailChangeRequest{
		Email:    email,
		NewEmail: email,
	}

	EditTestUserEmailWithError(conn, t, &emailChangeRequest, http.StatusBadRequest, apierror.UpdateParametersEqual)
}

func TestEditUserPassword(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())
	user := RegisterMockTestUser(conn, t)

	// change the user password
	userpass := mock.Password() + "xXx"
	user.Password = userpass
	EditTestUser(conn, t, user)

	// check if the user can login with the new password
	user.Password = userpass
	LoginTestUser(conn, t, user, mock.Ip(), mock.Platform())

	newUser, err := userdal.GetUser(conn, user, true)

	if err != nil {
		t.Error("The user was not found", err)
		return
	}

	if newUser.LastUpdate == user.LastUpdate {
		t.Error("The user last update date was not changed")
		return
	}

}

func TestEditUserPassworInvalidObjectId(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := &usersmodels.User{
		Email:    mock.Email(),
		Password: mock.Password(),
		Username: mock.Username(),
	}

	newUser := RegisterTestUser(conn, t, user)
	newUser.Password = user.Password
	newUser.Id = mock.InvalidId()

	EditTestUserWithError(conn, t, newUser, http.StatusBadRequest, apierror.InvalidObjectId)
}

func TestEditUserPasswordShort(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())
	user := RegisterMockTestUser(conn, t)

	// change the user password
	user.Password = mock.PasswordShort()
	EditTestUserWithError(conn, t, user, http.StatusBadRequest, apierror.ShortPassword)
}

func TestEditUserPasswordNoLowercase(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())
	user := RegisterMockTestUser(conn, t)

	// change the user password
	user.Password = mock.PasswordNotLowerCase()
	EditTestUserWithError(conn, t, user, http.StatusBadRequest, apierror.NoMayusMinusPassword)
}

func TestEditUserPasswordNoUppercase(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())
	user := RegisterMockTestUser(conn, t)

	// change the user password
	user.Password = mock.PasswordNotUpperCase()
	EditTestUserWithError(conn, t, user, http.StatusBadRequest, apierror.NoMayusMinusPassword)
}

func TestEditUserPasswordNoNumber(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())
	user := RegisterMockTestUser(conn, t)

	// change the user password
	user.Password = mock.PasswordNotNumber()
	EditTestUserWithError(conn, t, user, http.StatusBadRequest, apierror.NoAlphanumericPassword)
}

// func TestEditProfilePicture(t *testing.T) {

// 	user := RegisterMockTestUser(t)

// 	// Read the profile picture from the file
// 	profilePic, readErr := utils.ReadFile(mock.ProfilePicture())

// 	if readErr != nil {
// 		t.Error("The file was not read", readErr)
// 		return
// 	}

// 	err := EditUserProfilePicture(user, profilePic)

// 	if err != nil {
// 		t.Error("The profile picture was not changed", err)
// 		return
// 	}

// 	log.Info("Profile picture changed")

// }

func TestTokenValidation(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := RegisterMockTestUser(conn, t)
	LoginTestUser(conn, t, user, mock.Ip(), mock.Platform())
}

func TestTokenValidationInvalidToken(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	// Create a fake token
	ValidateTestTokenWithError(conn, t, mock.Token(), http.StatusForbidden, apierror.InvalidToken)

}

func TestTokenValidationInvalidTokenFormat(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	// Create a fake token
	token := mock.Username()
	ValidateTestTokenWithError(conn, t, token, http.StatusForbidden, apierror.InvalidToken)

}

func TestTokenValidationEmptyToken(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	// Create a fake token
	ValidateTestTokenWithError(conn, t, "", http.StatusForbidden, apierror.InvalidToken)
}

func TestValidationCode(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	// Create a new user
	user := RegisterMockTestUser(conn, t)

	// get the user
	user, err := userdal.GetUser(conn, user, true)

	if err != nil {
		t.Error("The user was not found", err)
		return
	}

	// validate the user
	err = userdal.ValidateUser(conn, user.ValidationCode)

	if err != nil {
		t.Error("The user was not validated", err)
		return
	}

}

func TestGetUser(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := RegisterMockTestUser(conn, t)
	found, err := userdal.GetUser(conn, user, true)

	if err != nil {
		t.Error("The user was not found", err)
	}

	if found == nil {
		t.Error("The user was not found")
		return
	}

}

func TestGetUserNotFound(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := RegisterMockTestUser(conn, t)
	user.Id = mock.InvalidId()
	_, err := userdal.GetUser(conn, user, true)

	if err == nil {
		t.Error("The user was found")
	}

}

func TestGetUserByEmail(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := RegisterMockTestUser(conn, t)
	found, err := userdal.GetUserByEmail(conn, user.Email, true)

	if err != nil {
		t.Error("The user was not found", err)
	}

	if found == nil {
		t.Error("The user was not found")
		return
	}

}
