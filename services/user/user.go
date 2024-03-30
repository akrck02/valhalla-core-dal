package userdal

import (
	"context"

	"github.com/akrck02/valhalla-core-dal/configuration"
	"github.com/akrck02/valhalla-core-dal/database"
	devicedal "github.com/akrck02/valhalla-core-dal/services/device"
	"github.com/akrck02/valhalla-core-sdk/error"
	"github.com/akrck02/valhalla-core-sdk/http"
	"github.com/akrck02/valhalla-core-sdk/log"
	"github.com/akrck02/valhalla-core-sdk/models"
	"github.com/akrck02/valhalla-core-sdk/utils"

	"github.com/golang-jwt/jwt/v5"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type EmailChangeRequest struct {
	Email    string `json:"email"`
	NewEmail string `json:"new_email"`
}

// Register user logic
//
// [param] conn | context.Context: connection to the database
// [param] client | *mongo.Client: client to the database
// [param] user | *models.User: user to register
//
// [return] *models.Error: error if any
func Register(conn context.Context, client *mongo.Client, user *models.User) *models.Error {

	if utils.IsEmpty(user.Email) {
		return &models.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   error.EMPTY_EMAIL,
			Message: "Email cannot be empty",
		}
	}

	if utils.IsEmpty(user.Password) {
		return &models.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   error.EMPTY_PASSWORD,
			Message: "Password cannot be empty",
		}
	}

	if utils.IsEmpty(user.Username) {
		return &models.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   error.EMPTY_USERNAME,
			Message: "Username cannot be empty",
		}
	}

	var checkedPass = utils.ValidatePassword(user.Password)

	if checkedPass.Response != 200 {
		return &models.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   int(checkedPass.Response),
			Message: checkedPass.Message,
		}
	}

	checkedPass = utils.ValidateEmail(user.Email)

	if checkedPass.Response != 200 {
		return &models.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   int(checkedPass.Response),
			Message: checkedPass.Message,
		}
	}

	coll := client.Database(database.CurrentDatabase).Collection(database.USER)
	found := mailExists(user.Email, conn, coll)

	if found != nil {

		return &models.Error{
			Status:  http.HTTP_STATUS_CONFLICT,
			Error:   error.USER_ALREADY_EXISTS,
			Message: "User already exists",
		}
	}

	code, err := utils.GenerateValidationCode(user.Email)

	if err != nil {
		return &models.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   error.CANNOT_CREATE_VALIDATION_CODE,
			Message: "User not created",
		}
	}

	userToInsert := user.Clone()
	userToInsert.Password = utils.EncryptSha256(user.Clone().Password)
	userToInsert.ValidationCode = code

	// register user on database
	_, err = coll.InsertOne(conn, userToInsert)

	if err != nil {
		return &models.Error{
			Status:  http.HTTP_STATUS_CONFLICT,
			Error:   error.USER_ALREADY_EXISTS,
			Message: "User already exists",
		}
	}

	return nil
}

// Login user logic
//
// [param] conn | context.Context: connection to the database
// [param] client | *mongo.Client: client to the database
// [param] user | models.User: user to login
// [param] ip | string: ip address of the user
// [param] address | string: user agent of the user
//
// [return] string: auth token --> *models.Error: error if any
func Login(conn context.Context, client *mongo.Client, user *models.User, ip string, address string) (string, *models.Error) {

	coll := client.Database(database.CurrentDatabase).Collection(database.USER)
	log.Info("Password: " + user.Password)
	found := authorizationOk(user.Email, user.Clone().Password, conn, coll)

	if found == nil {
		return "", &models.Error{
			Status:  http.HTTP_STATUS_FORBIDDEN,
			Error:   error.USER_NOT_AUTHORIZED,
			Message: "Invalid credentials",
		}
	}

	device := &models.Device{Address: ip, UserAgent: address}
	token, err := devicedal.AddUserDevice(conn, client, found, device)

	if err != nil {
		return "", &models.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   error.UNEXPECTED_ERROR,
			Message: "Cannot generate your auth token",
		}
	}

	return token, nil
}

// Login auth logic
//
// [param] conn | context.Context: connection to the database
// [param] client | *mongo.Client: client to the database
// [param] auth | models.AuthLogin: auth to login
func LoginAuth(conn context.Context, client *mongo.Client, auth *models.AuthLogin, ip string, userAgent string) *models.Error {

	found, err := GetUser(conn, client, &models.User{Email: auth.Email}, false)

	if err != nil {
		return err
	}

	// Search a user device with the same ip and user agent that has the token
	var filter = models.Device{
		User:      found.Email,
		UserAgent: userAgent,
		Address:   ip,
		Token:     auth.AuthToken,
	}

	var devices = client.Database(database.CurrentDatabase).Collection(database.DEVICE)
	device, deviceFindingError := devicedal.FindDeviceByAuthToken(conn, devices, &filter)

	if deviceFindingError != nil || device == nil {
		return &models.Error{
			Status:  http.HTTP_STATUS_NOT_FOUND,
			Error:   error.USER_NOT_AUTHORIZED,
			Message: "No possible login devices",
		}
	}

	return nil

}

// Edit user logic
//
// [param] conn | context.Context: connection to the database
// [param] client | *mongo.Client: client to the database
// [param] user | models.User: user to edit
//
// [return] *models.Error: error if any
func EditUser(conn context.Context, client *mongo.Client, user *models.User) *models.Error {

	users := client.Database(database.CurrentDatabase).Collection(database.USER)

	// validate email
	if user.Email != "" {
		checkedPass := utils.ValidateEmail(user.Email)

		if checkedPass.Response != 200 {
			return &models.Error{
				Status:  http.HTTP_STATUS_BAD_REQUEST,
				Error:   int(checkedPass.Response),
				Message: checkedPass.Message,
			}
		}
	}

	// validate password
	if user.Password != "" {

		checkedPass := utils.ValidatePassword(user.Password)

		if checkedPass.Response != 200 {
			return &models.Error{
				Status:  http.HTTP_STATUS_BAD_REQUEST,
				Error:   int(checkedPass.Response),
				Message: checkedPass.Message,
			}
		}
	}

	toUpdate := bson.M{"$set": bson.M{}}

	if user.Username != "" {
		toUpdate["$set"].(bson.M)["username"] = user.Username
	}

	if user.Password != "" {
		log.Info("password: " + user.Password)
		encryptedPass := user.Password
		toUpdate["$set"].(bson.M)["password"] = utils.EncryptSha256(encryptedPass)
		log.Info("encrypted password: " + encryptedPass)
	}

	if user.ProfilePic != "" {
		toUpdate["$set"].(bson.M)["profilePic"] = user.ProfilePic
	}

	// update user on database
	res, err := users.UpdateOne(conn, bson.M{"email": user.Email}, toUpdate)

	if err != nil {
		return &models.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   error.USER_NOT_UPDATED,
			Message: "User not updated",
		}
	}

	if res.MatchedCount == 0 && res.ModifiedCount == 0 {
		return &models.Error{
			Status:  http.HTTP_STATUS_NOT_FOUND,
			Error:   error.USER_NOT_FOUND,
			Message: "Users not found",
		}
	}

	return nil
}

// Change email logic
//
// [param] conn | context.Context: connection to the database
// [param] client | *mongo.Client: client to the database
// [param] user | models.User: user to change email
//
// [return] *models.Error: error if any
func EditUserEmail(conn context.Context, client *mongo.Client, mail *EmailChangeRequest) *models.Error {

	if utils.IsEmpty(mail.Email) || utils.IsEmpty(mail.NewEmail) {
		return &models.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   error.EMPTY_EMAIL,
			Message: "Email cannot be empty",
		}
	}

	// Equal emails
	if mail.Email == mail.NewEmail {
		return &models.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   error.EMAILS_EQUAL,
			Message: "The new email is the same as the old one",
		}
	}

	// validate email
	var checkedPass = utils.ValidateEmail(mail.Email)
	if checkedPass.Response != 200 {
		return &models.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   int(checkedPass.Response),
			Message: checkedPass.Message,
		}
	}

	// Check if user exists
	users := client.Database(database.CurrentDatabase).Collection(database.USER)
	found := mailExists(mail.NewEmail, conn, users)

	if found != nil {
		return &models.Error{
			Status:  http.HTTP_STATUS_CONFLICT,
			Error:   error.USER_ALREADY_EXISTS,
			Message: "That email is already in use",
		}
	}

	// update user on database
	var checkedEmail = utils.ValidateEmail(mail.NewEmail)
	if checkedEmail.Response != 200 {
		return &models.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   int(checkedEmail.Response),
			Message: checkedEmail.Message,
		}

	}

	updateStatus, err := users.UpdateOne(conn, bson.M{"email": mail.Email}, bson.M{"$set": bson.M{"email": mail.NewEmail}})

	if err != nil {
		return &models.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   error.USER_NOT_UPDATED,
			Message: "User not updated" + err.Error(),
		}
	}

	if updateStatus.MatchedCount == 0 {
		return &models.Error{
			Status:  http.HTTP_STATUS_NOT_FOUND,
			Error:   error.USER_NOT_FOUND,
			Message: "User not found",
		}
	}

	if updateStatus.ModifiedCount == 0 {
		return &models.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   error.USER_NOT_UPDATED,
			Message: "User not updated",
		}
	}

	// update user devices on database
	devices := client.Database(database.CurrentDatabase).Collection(database.DEVICE)

	updateStatus, err = devices.UpdateMany(conn, bson.M{"user": mail.Email}, bson.M{"$set": bson.M{"user": mail.NewEmail}})

	if err != nil {
		return &models.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   error.USER_NOT_UPDATED,
			Message: "User devices not updated",
		}
	}

	if updateStatus.MatchedCount != 0 && updateStatus.ModifiedCount == 0 {
		return &models.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   error.USER_NOT_UPDATED,
			Message: "User devices not updated",
		}
	}

	return nil
}

// Change profile picture logic
//
// [param] conn | context.Context: connection to the database
// [param] client | *mongo.Client: client to the database
// [param] user | models.User: user to change email
// [param] picture | []byte: picture to change
//
// [return] *models.Error: error if any
func EditUserProfilePicture(conn context.Context, client *mongo.Client, user *models.User, picture []byte) *models.Error {

	if utils.IsEmpty(user.Email) {
		return &models.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   error.EMPTY_EMAIL,
			Message: "Email cannot be empty",
		}
	}

	//if the path base profile pic does not exist, create it
	var profilePathDir = utils.GetProfilePicturePath("", configuration.PROFILE_PICTURES_PATH)
	if !utils.ExistsDir(profilePathDir) {
		err := utils.CreateDir(profilePathDir)

		if err != nil {
			return &models.Error{
				Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
				Error:   error.USER_NOT_UPDATED,
				Message: "User not updated, image not saved :" + err.Error(),
			}
		}
	}

	// save the picture
	var profilePicPath = utils.GetProfilePicturePath(user.Email, configuration.PROFILE_PICTURES_PATH)
	err := utils.SaveFile(profilePicPath, picture)

	if err != nil {
		return &models.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   error.USER_NOT_UPDATED,
			Message: "User not updated, image not saved :" + err.Error(),
		}
	}

	user.ProfilePic = profilePicPath
	editErr := EditUser(conn, client, user)

	if editErr != nil {
		return &models.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   error.USER_NOT_UPDATED,
			Message: "User not updated",
		}
	}

	return nil
}

// Delete user logic
//
// [param] conn | context.Context: connection to the database
// [param] client | *mongo.Client: client to the database
// [param] user | models.User: user to delete
//
// [return] *models.Error: error if any
func DeleteUser(conn context.Context, client *mongo.Client, user *models.User) *models.Error {

	if utils.IsEmpty(user.Email) {
		return &models.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   error.EMPTY_EMAIL,
			Message: "Email cannot be empty",
		}
	}

	// delete user projects
	projects := client.Database(database.CurrentDatabase).Collection(database.PROJECT)
	_, err := projects.DeleteMany(conn, bson.M{"owner": user.Email})

	if err != nil {
		return &models.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   error.USER_NOT_DELETED,
			Message: "User not deleted",
		}
	}

	// delete user devices
	devices := client.Database(database.CurrentDatabase).Collection(database.DEVICE)
	_, err = devices.DeleteMany(conn, bson.M{"user": user.Email})

	if err != nil {
		return &models.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   error.USER_NOT_DELETED,
			Message: "User not deleted",
		}
	}

	// delete user on database
	users := client.Database(database.CurrentDatabase).Collection(database.USER)

	var deleteResult *mongo.DeleteResult
	deleteResult, err = users.DeleteOne(conn, bson.M{"email": user.Email})

	if err != nil {
		return &models.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   error.USER_NOT_DELETED,
			Message: "User not deleted",
		}
	}

	if deleteResult.DeletedCount == 0 {
		return &models.Error{
			Status:  http.HTTP_STATUS_NOT_FOUND,
			Error:   error.USER_NOT_FOUND,
			Message: "User not found",
		}
	}

	return nil
}

// Get user logic
func GetUser(conn context.Context, client *mongo.Client, user *models.User, secure bool) (*models.User, *models.Error) { // get user from database

	users := client.Database(database.CurrentDatabase).Collection(database.USER)
	var found models.User
	err := users.FindOne(conn, bson.M{"email": user.Email}).Decode(&found)

	if err != nil {
		return nil, &models.Error{
			Status:  http.HTTP_STATUS_NOT_FOUND,
			Error:   error.USER_NOT_FOUND,
			Message: "User not found",
		}
	}

	if secure {
		found.Password = "****************"
	}

	return &found, nil
}

// Validate user logic
//
// [param] conn | context.Context: connection to the database
// [param] client | *mongo.Client: client to the database
// [param] code | string: code to validate
//
// [return] *models.Error: error if any
func ValidateUser(conn context.Context, client *mongo.Client, code string) *models.Error {

	if utils.IsEmpty(code) {
		return &models.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   error.INVALID_VALIDATION_CODE,
			Message: "Code cannot be empty",
		}
	}

	var user = &models.User{
		ValidationCode: code,
	}
	coll := client.Database(database.CurrentDatabase).Collection(database.USER)
	err := coll.FindOne(conn, user).Decode(user)

	log.FormattedInfo("user: ${0}", user.Email)
	log.FormattedInfo("code: ${0}", code)

	if err != nil {
		return &models.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   error.INVALID_VALIDATION_CODE,
			Message: "Invalid validation code",
		}
	}

	if user.Validated {
		return &models.Error{
			Status:  http.HTTP_STATUS_OK,
			Error:   error.USER_ALREADY_VALIDATED,
			Message: "User already validated",
		}
	}

	if user.ValidationCode != code {
		return &models.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   error.INVALID_VALIDATION_CODE,
			Message: "Invalid validation code",
		}
	}

	user.ValidationCode = ""
	user.Validated = true

	// update user on database
	result, editerr := coll.UpdateOne(conn, bson.M{"email": user.Email}, bson.M{"$set": bson.M{"validation_code": "", "validated": true}})

	if result.MatchedCount == 0 {
		return &models.Error{
			Status:  http.HTTP_STATUS_NOT_FOUND,
			Error:   error.USER_NOT_FOUND,
			Message: "User not found",
		}
	}

	if editerr != nil {
		return &models.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   error.USER_NOT_UPDATED,
			Message: "User not validated: " + editerr.Error(),
		}
	}

	return nil
}

// Check email on database
//
//	[param] email | string The email to check
//	[param] conn | context.Context The connection to the database
//
//	[return] model.User : The user found or empty
func mailExists(email string, conn context.Context, coll *mongo.Collection) *models.User {

	filter := bson.D{{Key: "email", Value: email}}

	var result models.User
	err := coll.FindOne(conn, filter).Decode(&result)

	if err != nil {
		log.FormattedError("Error: ${0}", err.Error())
		return nil
	}

	return &result
}

// Get if the given credentials are valid
//
//	[param] username | string : The username to check
//	[param] password | string : The password to check
//	[param] conn | context.Context : The connection to the database
//
//	[return] model.User : The user found or empty
func authorizationOk(email string, password string, conn context.Context, coll *mongo.Collection) *models.User {

	filter := bson.D{
		{Key: "email", Value: email},
		{Key: "password", Value: utils.EncryptSha256(password)},
	}

	var result models.User
	err := coll.FindOne(conn, filter).Decode(&result)

	if err != nil {
		log.FormattedError("Error: ${0}", err.Error())
		return nil
	}

	return &result
}

// Get user from token
//
//	[param] conn | context.Context : The connection to the database
//	[param] client | *mongo.Client : The client to the database
//	[param] token | *string : The token to check
//	[param] tokenUser | *models.User : The user found or empty --> *models.Error: error if any
func getUserFromToken(conn context.Context, client *mongo.Client, token string) (models.User, *models.Error) {

	var tokenDevice models.Device

	devices := client.Database(database.CurrentDatabase).Collection(database.DEVICE)
	err := devices.FindOne(conn, bson.M{"token": token}).Decode(&tokenDevice)

	if err != nil {
		return models.User{}, &models.Error{
			Status:  http.HTTP_STATUS_FORBIDDEN,
			Error:   error.INVALID_TOKEN,
			Message: "User not matching token",
		}
	}

	var tokenUser models.User

	users := client.Database(database.CurrentDatabase).Collection(database.USER)
	err = users.FindOne(conn, bson.M{"email": tokenDevice.User}).Decode(&tokenUser)

	if err != nil {
		return models.User{}, &models.Error{
			Status:  http.HTTP_STATUS_FORBIDDEN,
			Error:   error.INVALID_TOKEN,
			Message: "User not matching token",
		}
	}

	return tokenUser, nil
}

// Get  if token is valid
//
//	[param] token | string : The token to check
//
//	[return] bool : True if token is valid --> *models.Error: error if any
func IsTokenValid(client *mongo.Client, token string) (*models.User, *models.Error) {

	// decode token
	claims, err := utils.DecryptToken(token, configuration.Params.Secret)

	if err != nil {
		return nil, &models.Error{
			Status:  http.HTTP_STATUS_FORBIDDEN,
			Error:   error.INVALID_TOKEN,
			Message: "invalid token format",
		}
	}

	// log token claims
	log.Info("device: " + claims.Claims.(jwt.MapClaims)["device"].(string))
	log.Info("username: " + claims.Claims.(jwt.MapClaims)["username"].(string))
	log.Info("email: " + claims.Claims.(jwt.MapClaims)["email"].(string))

	email := claims.Claims.(jwt.MapClaims)["email"].(string)

	foundUser, tokenUserErr := getUserFromToken(context.Background(), client, token)

	if tokenUserErr != nil {
		return nil, &models.Error{
			Status:  http.HTTP_STATUS_FORBIDDEN,
			Error:   error.INVALID_TOKEN,
			Message: "invalid token",
		}
	}

	if foundUser.Email != email {
		return nil, &models.Error{
			Status:  http.HTTP_STATUS_FORBIDDEN,
			Error:   error.INVALID_TOKEN,
			Message: "invalid token",
		}
	}

	return &foundUser, nil
}
