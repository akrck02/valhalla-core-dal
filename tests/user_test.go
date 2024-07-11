package tests

import (
	"context"
	"testing"

	"github.com/akrck02/valhalla-core-dal/database"
	userdal "github.com/akrck02/valhalla-core-dal/services/user"
	"github.com/akrck02/valhalla-core-sdk/http"
	"github.com/akrck02/valhalla-core-sdk/log"
	usersmodels "github.com/akrck02/valhalla-core-sdk/models/users"
	"github.com/akrck02/valhalla-core-sdk/valerror"

	"github.com/akrck02/valhalla-core-dal/mock"
)

func TestRegister(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := RegisterMockTestUser(conn, t)
	DeleteTestUser(conn, t, user)
}

func TestRegisterNotEmail(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := &usersmodels.User{
		Password: mock.Password(),
		Username: mock.Username(),
	}

	RegisterTestUserWithError(conn, t, user, http.HTTP_STATUS_BAD_REQUEST, valerror.EMPTY_USER_EMAIL)
}

func TestRegisterNotUsername(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := &usersmodels.User{
		Email:    mock.Email(),
		Password: mock.Password(),
	}

	RegisterTestUserWithError(conn, t, user, http.HTTP_STATUS_BAD_REQUEST, valerror.EMPTY_USERNAME)
}

func TestRegisterNotPassword(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := &usersmodels.User{
		Email:    mock.Email(),
		Username: mock.Username(),
	}

	RegisterTestUserWithError(conn, t, user, http.HTTP_STATUS_BAD_REQUEST, valerror.EMPTY_USER_PASSWORD)
}

func TestRegisterNotDotEmail(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := &usersmodels.User{
		Email:    mock.EmailNotDot(),
		Password: mock.Password(),
		Username: mock.Username(),
	}

	RegisterTestUserWithError(conn, t, user, http.HTTP_STATUS_BAD_REQUEST, valerror.NO_DOT_EMAIL)
}

func TestRegisterNotAtEmail(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := &usersmodels.User{
		Email:    mock.EmailNotAt(),
		Password: mock.Password(),
		Username: mock.Username(),
	}

	RegisterTestUserWithError(conn, t, user, http.HTTP_STATUS_BAD_REQUEST, valerror.NO_AT_EMAIL)
}

func TestRegisterShortMail(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := &usersmodels.User{
		Email:    mock.EmailShort(),
		Password: mock.Password(),
		Username: mock.Username(),
	}

	RegisterTestUserWithError(conn, t, user, http.HTTP_STATUS_BAD_REQUEST, valerror.SHORT_EMAIL)
}

func TestRegisterNotSpecialCharactersPassword(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := &usersmodels.User{
		Email:    mock.Email(),
		Password: mock.PasswordNotSpecialChar(),
		Username: mock.Username(),
	}

	RegisterTestUserWithError(conn, t, user, http.HTTP_STATUS_BAD_REQUEST, valerror.NO_SPECIAL_CHARACTERS_PASSWORD)
}

func TestRegisterNotUpperCaseLowerCasePassword(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := &usersmodels.User{
		Email:    mock.Email(),
		Password: mock.PasswordNotLowerCase(),
		Username: mock.Username(),
	}

	RegisterTestUserWithError(conn, t, user, http.HTTP_STATUS_BAD_REQUEST, valerror.NO_UPPER_LOWER_PASSWORD)

	user = &usersmodels.User{
		Email:    mock.Email(),
		Password: mock.PasswordNotUpperCase(),
		Username: mock.Username(),
	}

	RegisterTestUserWithError(conn, t, user, http.HTTP_STATUS_BAD_REQUEST, valerror.NO_UPPER_LOWER_PASSWORD)
}

func TestRegisterShortPassword(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := &usersmodels.User{
		Email:    mock.Email(),
		Password: mock.PasswordShort(),
		Username: mock.Username(),
	}

	RegisterTestUserWithError(conn, t, user, http.HTTP_STATUS_BAD_REQUEST, valerror.SHORT_PASSWORD)
}

func TestRegisterNotNumbersPassword(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := &usersmodels.User{
		Email:    mock.Email(),
		Password: mock.PasswordNotNumber(),
		Username: mock.Username(),
	}

	RegisterTestUserWithError(conn, t, user, http.HTTP_STATUS_BAD_REQUEST, valerror.NO_ALPHANUMERIC_PASSWORD)
}

func TestLogin(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := RegisterMockTestUser(conn, t)
	LoginTestUser(conn, t, user, mock.Ip(), mock.Platform())
	DeleteTestUser(conn, t, user)

}

func TestLoginWrongPassword(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := RegisterMockTestUser(conn, t)
	user.Password = mock.PasswordShort()

	log.Info("Login with wrong password")
	log.FormattedInfo("Password: ${0}", user.Password)

	// login the user
	LoginTestUserWithError(conn, t, user, mock.Ip(), mock.Platform(), http.HTTP_STATUS_FORBIDDEN, valerror.USER_NOT_AUTHORIZED)
	DeleteTestUser(conn, t, user)
}

func TestLoginWrongEmail(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := RegisterMockTestUser(conn, t)
	realEmail := user.Email
	user.Email = "wrong" + mock.Email()
	LoginTestUserWithError(conn, t, user, mock.Ip(), mock.Platform(), http.HTTP_STATUS_FORBIDDEN, valerror.USER_NOT_AUTHORIZED)

	user.Email = realEmail
	DeleteTestUser(conn, t, user)
}

func TestLoginAuth(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := RegisterMockTestUser(conn, t)
	token := LoginTestUser(conn, t, user, mock.Ip(), mock.Platform())
	LoginAuthTestUser(conn, t, user.Email, token)
	DeleteTestUser(conn, t, user)

}

func TestDeleteUser(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := RegisterMockTestUser(conn, t)
	LoginTestUser(conn, t, user, mock.Ip(), mock.Platform())
	DeleteTestUser(conn, t, user)

}

func TestDeleteUserNoEmail(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := &usersmodels.User{
		Password: mock.Password(),
		Username: mock.Username(),
	}

	DeleteTestUserWithError(conn, t, user, http.HTTP_STATUS_BAD_REQUEST, valerror.EMPTY_USER_EMAIL)
}

func TestDeleteUserNotFound(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := &usersmodels.User{
		Email:    mock.Email(),
		Password: mock.Password(),
		Username: mock.Username(),
	}

	DeleteTestUserWithError(conn, t, user, http.HTTP_STATUS_NOT_FOUND, valerror.USER_NOT_FOUND)
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

	RegisterTestUser(conn, t, user)
	LoginTestUser(conn, t, user, mock.Ip(), mock.Platform())

	// Change the user email
	log.Info("Changing user email")
	emailChangeRequest := userdal.EmailChangeRequest{
		Email:    email,
		NewEmail: newEmail,
	}

	EditTestUserEmail(conn, t, &emailChangeRequest)
	user.Email = newEmail

	// delete the user
	DeleteTestUser(conn, t, user)
}

func TestEditUserEmailNoEmail(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	emailChangeRequest := userdal.EmailChangeRequest{}
	EditTestUserEmailWithError(conn, t, &emailChangeRequest, http.HTTP_STATUS_BAD_REQUEST, valerror.EMPTY_USER_EMAIL)

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

	EditTestUserEmailWithError(conn, t, &emailChangeRequest, http.HTTP_STATUS_BAD_REQUEST, valerror.NO_DOT_EMAIL)

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

	EditTestUserEmailWithError(conn, t, &emailChangeRequest, http.HTTP_STATUS_BAD_REQUEST, valerror.NO_AT_EMAIL)

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

	EditTestUserEmailWithError(conn, t, &emailChangeRequest, http.HTTP_STATUS_BAD_REQUEST, valerror.SHORT_EMAIL)

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

	EditTestUserEmailWithError(conn, t, &emailChangeRequest, http.HTTP_STATUS_NOT_FOUND, valerror.USER_NOT_FOUND)
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

	EditTestUserEmailWithError(conn, t, &emailChangeRequest, http.HTTP_STATUS_CONFLICT, valerror.USER_ALREADY_EXISTS)
	log.Jump()

	// Delete the users
	DeleteTestUser(conn, t, user)
	DeleteTestUser(conn, t, newUser)
}

func TestEditUserSameEmail(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	email := mock.Email()
	emailChangeRequest := userdal.EmailChangeRequest{
		Email:    email,
		NewEmail: email,
	}

	EditTestUserEmailWithError(conn, t, &emailChangeRequest, http.HTTP_STATUS_BAD_REQUEST, valerror.USER_EDITING_EMAILS_EQUAL)
}

func TestEditUserPassword(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())
	user := RegisterMockTestUser(conn, t)

	// change the user password
	user.Password = mock.Password() + "xXx"
	EditTestUser(conn, t, user)

	// check if the user can login with the new password
	LoginTestUser(conn, t, user, mock.Ip(), mock.Platform())

	// delete the user
	DeleteTestUser(conn, t, user)
}

func TestEditUserPasswordUserNotFound(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := &usersmodels.User{
		Email:    mock.Email(),
		Password: mock.Password(),
		Username: mock.Username(),
	}

	EditTestUserWithError(conn, t, user, http.HTTP_STATUS_NOT_FOUND, valerror.USER_NOT_FOUND)
}

func TestEditUserPasswordShort(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())
	user := RegisterMockTestUser(conn, t)

	// change the user password
	user.Password = mock.PasswordShort()
	EditTestUserWithError(conn, t, user, http.HTTP_STATUS_BAD_REQUEST, valerror.SHORT_PASSWORD)

	// delete the user
	DeleteTestUser(conn, t, user)
}

func TestEditUserPasswordNoLowercase(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())
	user := RegisterMockTestUser(conn, t)

	// change the user password
	user.Password = mock.PasswordNotLowerCase()
	EditTestUserWithError(conn, t, user, http.HTTP_STATUS_BAD_REQUEST, valerror.NO_UPPER_LOWER_PASSWORD)

	// delete the user
	DeleteTestUser(conn, t, user)
}

func TestEditUserPasswordNoUppercase(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())
	user := RegisterMockTestUser(conn, t)

	// change the user password
	user.Password = mock.PasswordNotUpperCase()
	EditTestUserWithError(conn, t, user, http.HTTP_STATUS_BAD_REQUEST, valerror.NO_UPPER_LOWER_PASSWORD)

	// delete the user
	DeleteTestUser(conn, t, user)

}

func TestEditUserPasswordNoNumber(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())
	user := RegisterMockTestUser(conn, t)

	// change the user password
	user.Password = mock.PasswordNotNumber()
	EditTestUserWithError(conn, t, user, http.HTTP_STATUS_BAD_REQUEST, valerror.NO_ALPHANUMERIC_PASSWORD)

	// delete the user
	DeleteTestUser(conn, t, user)

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

// 	// delete the user
// 	DeleteTestUser(t, user)

// }

func TestTokenValidation(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	user := RegisterMockTestUser(conn, t)
	LoginTestUser(conn, t, user, mock.Ip(), mock.Platform())
	DeleteTestUser(conn, t, user)
}

func TestTokenValidationInvalidToken(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	// Create a fake token
	ValidateTestTokenWithError(conn, t, mock.Token(), http.HTTP_STATUS_FORBIDDEN, valerror.INVALID_TOKEN)

}

func TestTokenValidationInvalidTokenFormat(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	// Create a fake token
	token := mock.Username()
	ValidateTestTokenWithError(conn, t, token, http.HTTP_STATUS_FORBIDDEN, valerror.INVALID_TOKEN)

}

func TestTokenValidationEmptyToken(t *testing.T) {

	conn := database.Connect()
	defer conn.Disconnect(context.Background())

	// Create a fake token
	ValidateTestTokenWithError(conn, t, "", http.HTTP_STATUS_FORBIDDEN, valerror.INVALID_TOKEN)
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

	// delete the user
	DeleteTestUser(conn, t, user)

}
