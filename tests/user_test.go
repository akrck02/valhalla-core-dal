package tests

import (
	"testing"

	"github.com/akrck02/valhalla-core-dal/database"
	userdal "github.com/akrck02/valhalla-core-dal/services/user"
	"github.com/akrck02/valhalla-core-sdk/error"
	"github.com/akrck02/valhalla-core-sdk/http"
	"github.com/akrck02/valhalla-core-sdk/log"
	"github.com/akrck02/valhalla-core-sdk/models"

	"github.com/akrck02/valhalla-core-dal/mock"
)

func TestRegister(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	user := RegisterMockTestUser(t, client)
	DeleteTestUser(t, client, user)
}

func TestRegisterNotEmail(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	user := &models.User{
		Password: mock.Password(),
		Username: mock.Username(),
	}

	RegisterTestUserWithError(t, client, user, http.HTTP_STATUS_BAD_REQUEST, error.EMPTY_EMAIL)
}

func TestRegisterNotUsername(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	user := &models.User{
		Email:    mock.Email(),
		Password: mock.Password(),
	}

	RegisterTestUserWithError(t, client, user, http.HTTP_STATUS_BAD_REQUEST, error.EMPTY_USERNAME)
}

func TestRegisterNotPassword(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	user := &models.User{
		Email:    mock.Email(),
		Username: mock.Username(),
	}

	RegisterTestUserWithError(t, client, user, http.HTTP_STATUS_BAD_REQUEST, error.EMPTY_PASSWORD)
}

func TestRegisterNotDotEmail(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	user := &models.User{
		Email:    mock.EmailNotDot(),
		Password: mock.Password(),
		Username: mock.Username(),
	}

	RegisterTestUserWithError(t, client, user, http.HTTP_STATUS_BAD_REQUEST, error.NO_DOT_EMAIL)
}

func TestRegisterNotAtEmail(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	user := &models.User{
		Email:    mock.EmailNotAt(),
		Password: mock.Password(),
		Username: mock.Username(),
	}

	RegisterTestUserWithError(t, client, user, http.HTTP_STATUS_BAD_REQUEST, error.NO_AT_EMAIL)
}

func TestRegisterShortMail(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	user := &models.User{
		Email:    mock.EmailShort(),
		Password: mock.Password(),
		Username: mock.Username(),
	}

	RegisterTestUserWithError(t, client, user, http.HTTP_STATUS_BAD_REQUEST, error.SHORT_EMAIL)
}

func TestRegisterNotSpecialCharactersPassword(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	user := &models.User{
		Email:    mock.Email(),
		Password: mock.PasswordNotSpecialChar(),
		Username: mock.Username(),
	}

	RegisterTestUserWithError(t, client, user, http.HTTP_STATUS_BAD_REQUEST, error.NO_SPECIAL_CHARACTERS_PASSWORD)
}

func TestRegisterNotUpperCaseLowerCasePassword(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	user := &models.User{
		Email:    mock.Email(),
		Password: mock.PasswordNotLowerCase(),
		Username: mock.Username(),
	}

	RegisterTestUserWithError(t, client, user, http.HTTP_STATUS_BAD_REQUEST, error.NO_UPPER_LOWER_PASSWORD)

	user = &models.User{
		Email:    mock.Email(),
		Password: mock.PasswordNotUpperCase(),
		Username: mock.Username(),
	}

	RegisterTestUserWithError(t, client, user, http.HTTP_STATUS_BAD_REQUEST, error.NO_UPPER_LOWER_PASSWORD)
}

func TestRegisterShortPassword(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	user := &models.User{
		Email:    mock.Email(),
		Password: mock.PasswordShort(),
		Username: mock.Username(),
	}

	RegisterTestUserWithError(t, client, user, http.HTTP_STATUS_BAD_REQUEST, error.SHORT_PASSWORD)
}

func TestRegisterNotNumbersPassword(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	user := &models.User{
		Email:    mock.Email(),
		Password: mock.PasswordNotNumber(),
		Username: mock.Username(),
	}

	RegisterTestUserWithError(t, client, user, http.HTTP_STATUS_BAD_REQUEST, error.NO_ALPHANUMERIC_PASSWORD)
}

func TestLogin(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	user := RegisterMockTestUser(t, client)
	LoginTestUser(t, client, user, mock.Ip(), mock.Platform())
	DeleteTestUser(t, client, user)

}

func TestLoginWrongPassword(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	user := RegisterMockTestUser(t, client)
	user.Password = mock.PasswordShort()

	log.Info("Login with wrong password")
	log.FormattedInfo("Password: ${0}", user.Password)

	// login the user
	LoginTestUserWithError(t, client, user, mock.Ip(), mock.Platform(), http.HTTP_STATUS_FORBIDDEN, error.USER_NOT_AUTHORIZED)
	DeleteTestUser(t, client, user)
}

func TestLoginWrongEmail(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	user := RegisterMockTestUser(t, client)
	realEmail := user.Email
	user.Email = "wrong" + mock.Email()
	LoginTestUserWithError(t, client, user, mock.Ip(), mock.Platform(), http.HTTP_STATUS_FORBIDDEN, error.USER_NOT_AUTHORIZED)

	user.Email = realEmail
	userdal.DeleteUser(client, user)
}

func TestLoginAuth(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	user := RegisterMockTestUser(t, client)
	token := LoginTestUser(t, client, user, mock.Ip(), mock.Platform())
	LoginAuthTestUser(t, client, user.Email, token)
	DeleteTestUser(t, client, user)

}

func TestDeleteUser(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	user := RegisterMockTestUser(t, client)
	LoginTestUser(t, client, user, mock.Ip(), mock.Platform())
	DeleteTestUser(t, client, user)

}

func TestDeleteUserNoEmail(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	user := &models.User{
		Password: mock.Password(),
		Username: mock.Username(),
	}

	DeleteTestUserWithError(t, client, user, http.HTTP_STATUS_BAD_REQUEST, error.EMPTY_EMAIL)
}

func TestDeleteUserNotFound(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	user := &models.User{
		Email:    mock.Email(),
		Password: mock.Password(),
		Username: mock.Username(),
	}

	DeleteTestUserWithError(t, client, user, http.HTTP_STATUS_NOT_FOUND, error.USER_NOT_FOUND)
}

func TestEditUserEmail(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	email := mock.Email()
	newEmail := "xXx" + mock.Email()

	user := &models.User{
		Email:    email,
		Password: mock.Password(),
		Username: mock.Username(),
	}

	RegisterTestUser(t, client, user)
	LoginTestUser(t, client, user, mock.Ip(), mock.Platform())

	// Change the user email
	log.Info("Changing user email")
	emailChangeRequest := userdal.EmailChangeRequest{
		Email:    email,
		NewEmail: newEmail,
	}

	EditTestUserEmail(t, client, &emailChangeRequest)
	user.Email = newEmail

	// delete the user
	DeleteTestUser(t, client, user)
}

func TestEditUserEmailNoEmail(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	emailChangeRequest := userdal.EmailChangeRequest{}
	EditTestUserEmailWithError(t, client, &emailChangeRequest, http.HTTP_STATUS_BAD_REQUEST, error.EMPTY_EMAIL)

}

func TestEditUserEmailNoDotEmail(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	email := mock.Email()
	newEmail := mock.EmailNotDot()

	emailChangeRequest := userdal.EmailChangeRequest{
		Email:    email,
		NewEmail: newEmail,
	}

	EditTestUserEmailWithError(t, client, &emailChangeRequest, http.HTTP_STATUS_BAD_REQUEST, error.NO_DOT_EMAIL)

}

func TestEditUserEmailNoAtEmail(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	email := mock.Email()
	newEmail := mock.EmailNotAt()

	emailChangeRequest := userdal.EmailChangeRequest{
		Email:    email,
		NewEmail: newEmail,
	}

	EditTestUserEmailWithError(t, client, &emailChangeRequest, http.HTTP_STATUS_BAD_REQUEST, error.NO_AT_EMAIL)

}

func TestEditUserEmailShortEmail(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	email := mock.Email()
	newEmail := mock.EmailShort()

	emailChangeRequest := userdal.EmailChangeRequest{
		Email:    email,
		NewEmail: newEmail,
	}

	EditTestUserEmailWithError(t, client, &emailChangeRequest, http.HTTP_STATUS_BAD_REQUEST, error.SHORT_EMAIL)

}

func TestEditUserEmailNotFound(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	email := mock.Email()
	newEmail := "xXx" + mock.Email()

	emailChangeRequest := userdal.EmailChangeRequest{
		Email:    email,
		NewEmail: newEmail,
	}

	EditTestUserEmailWithError(t, client, &emailChangeRequest, http.HTTP_STATUS_NOT_FOUND, error.USER_NOT_FOUND)
}

func TestEditUserEmailExists(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	email := mock.Email()
	newEmail := mock.Email() + "xXx"

	// Create a user
	user := &models.User{
		Email:    email,
		Password: mock.Password(),
		Username: mock.Username(),
	}

	RegisterTestUser(t, client, user)
	log.Jump()

	// Create a new user with the new email
	newUser := &models.User{
		Email:    newEmail,
		Password: mock.Password(),
		Username: mock.Username(),
	}

	RegisterTestUser(t, client, newUser)
	log.Jump()

	// Change the email
	emailChangeRequest := userdal.EmailChangeRequest{
		Email:    email,
		NewEmail: newEmail,
	}

	EditTestUserEmailWithError(t, client, &emailChangeRequest, http.HTTP_STATUS_CONFLICT, error.USER_ALREADY_EXISTS)
	log.Jump()

	// Delete the users
	DeleteTestUser(t, client, user)
	DeleteTestUser(t, client, newUser)
}

func TestEditUserSameEmail(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	email := mock.Email()
	emailChangeRequest := userdal.EmailChangeRequest{
		Email:    email,
		NewEmail: email,
	}

	EditTestUserEmailWithError(t, client, &emailChangeRequest, http.HTTP_STATUS_BAD_REQUEST, error.EMAILS_EQUAL)
}

func TestEditUserPassword(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	user := RegisterMockTestUser(t, client)

	// change the user password
	user.Password = mock.Password() + "xXx"
	EditTestUser(t, client, user)

	// check if the user can login with the new password
	LoginTestUser(t, client, user, mock.Ip(), mock.Platform())

	// delete the user
	DeleteTestUser(t, client, user)
}

func TestEditUserPasswordUserNotFound(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	user := &models.User{
		Email:    mock.Email(),
		Password: mock.Password(),
		Username: mock.Username(),
	}

	EditTestUserWithError(t, client, user, http.HTTP_STATUS_NOT_FOUND, error.USER_NOT_FOUND)
}

func TestEditUserPasswordShort(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	user := RegisterMockTestUser(t, client)

	// change the user password
	user.Password = mock.PasswordShort()
	EditTestUserWithError(t, client, user, http.HTTP_STATUS_BAD_REQUEST, error.SHORT_PASSWORD)

	// delete the user
	DeleteTestUser(t, client, user)
}

func TestEditUserPasswordNoLowercase(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	user := RegisterMockTestUser(t, client)

	// change the user password
	user.Password = mock.PasswordNotLowerCase()
	EditTestUserWithError(t, client, user, http.HTTP_STATUS_BAD_REQUEST, error.NO_UPPER_LOWER_PASSWORD)

	// delete the user
	DeleteTestUser(t, client, user)
}

func TestEditUserPasswordNoUppercase(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	user := RegisterMockTestUser(t, client)

	// change the user password
	user.Password = mock.PasswordNotUpperCase()
	EditTestUserWithError(t, client, user, http.HTTP_STATUS_BAD_REQUEST, error.NO_UPPER_LOWER_PASSWORD)

	// delete the user
	DeleteTestUser(t, client, user)

}

func TestEditUserPasswordNoNumber(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	user := RegisterMockTestUser(t, client)

	// change the user password
	user.Password = mock.PasswordNotNumber()
	EditTestUserWithError(t, client, user, http.HTTP_STATUS_BAD_REQUEST, error.NO_ALPHANUMERIC_PASSWORD)

	// delete the user
	DeleteTestUser(t, client, user)

}

// func TestEditProfilePicture(t *testing.T) {

// // Connect database
// var client = database.Connect()
// defer database.Disconnect(*client)

// 	user := RegisterMockTestUser(t, client)

// 	// Read the profile picture from the file
// 	profilePic, readErr := utils.ReadFile(mock.ProfilePicture())

// 	if readErr != nil {
// 		t.Error("The file was not read", readErr)
// 		return
// 	}

// 	err := EditUserProfilePicture(client, user, profilePic)

// 	if err != nil {
// 		t.Error("The profile picture was not changed", err)
// 		return
// 	}

// 	log.Info("Profile picture changed")

// 	// delete the user
// 	DeleteTestUser(t, client, user)

// }

func TestTokenValidation(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	user := RegisterMockTestUser(t, client)
	LoginTestUser(t, client, user, mock.Ip(), mock.Platform())
	DeleteTestUser(t, client, user)
}

func TestTokenValidationInvalidToken(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	// Create a fake token
	ValidateTestTokenWithError(t, client, mock.Token(), http.HTTP_STATUS_FORBIDDEN, error.INVALID_TOKEN)

}

func TestTokenValidationInvalidTokenFormat(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	// Create a fake token
	token := mock.Username()
	ValidateTestTokenWithError(t, client, token, http.HTTP_STATUS_FORBIDDEN, error.INVALID_TOKEN)

}

func TestTokenValidationEmptyToken(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	// Create a fake token
	ValidateTestTokenWithError(t, client, "", http.HTTP_STATUS_FORBIDDEN, error.INVALID_TOKEN)
}

func TestValidationCode(t *testing.T) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	// Create a new user
	user := RegisterMockTestUser(t, client)

	// get the user
	user, err := userdal.GetUser(client, user, true)

	if err != nil {
		t.Error("The user was not found", err)
		return
	}

	// validate the user
	err = userdal.ValidateUser(client, user.ValidationCode)

	if err != nil {
		t.Error("The user was not validated", err)
		return
	}

	// delete the user
	DeleteTestUser(t, client, user)

}
