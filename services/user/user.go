package userdal

import (
	"github.com/akrck02/valhalla-core-dal/configuration"
	"github.com/akrck02/valhalla-core-dal/database"
	devicedal "github.com/akrck02/valhalla-core-dal/services/device"
	"github.com/akrck02/valhalla-core-sdk/http"
	"github.com/akrck02/valhalla-core-sdk/log"
	devicemodels "github.com/akrck02/valhalla-core-sdk/models/device"
	systemmodels "github.com/akrck02/valhalla-core-sdk/models/system"
	usersmodels "github.com/akrck02/valhalla-core-sdk/models/users"
	"github.com/akrck02/valhalla-core-sdk/utils"
	"github.com/akrck02/valhalla-core-sdk/valerror"

	"github.com/golang-jwt/jwt/v5"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EmailChangeRequest struct {
	Email    string `json:"email"`
	NewEmail string `json:"new_email"`
}

func Register(user *usersmodels.User) *systemmodels.Error {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	if utils.IsEmpty(user.Email) {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.EMPTY_USER_EMAIL,
			Message: "Email cannot be empty",
		}
	}

	if utils.IsEmpty(user.Password) {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.EMPTY_USER_PASSWORD,
			Message: "Password cannot be empty",
		}
	}

	if utils.IsEmpty(user.Username) {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.EMPTY_USERNAME,
			Message: "Username cannot be empty",
		}
	}

	var checkedPass = utils.ValidatePassword(user.Password)

	if checkedPass.Response != 200 {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   int(checkedPass.Response),
			Message: checkedPass.Message,
		}
	}

	checkedPass = utils.ValidateEmail(user.Email)

	if checkedPass.Response != 200 {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   checkedPass.Response,
			Message: checkedPass.Message,
		}
	}

	coll := client.Database(database.CurrentDatabase).Collection(database.USER)
	found := mailExists(user.Email, coll)

	if found != nil {

		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_CONFLICT,
			Error:   valerror.USER_ALREADY_EXISTS,
			Message: "User already exists",
		}
	}

	code, err := utils.GenerateValidationCode(user.Email)

	if err != nil {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   valerror.CANNOT_CREATE_VALIDATION_CODE,
			Message: "User not created",
		}
	}

	userToInsert := user.Clone()
	userToInsert.Password = utils.EncryptSha256(user.Clone().Password)
	userToInsert.ValidationCode = code

	// register user on database
	res, errInsert := coll.InsertOne(database.GetDefaultContext(), userToInsert)

	if errInsert != nil {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_CONFLICT,
			Error:   valerror.USER_ALREADY_EXISTS,
			Message: "User already exists",
		}
	}

	// get user from database
	user.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func Login(user *usersmodels.User, ip string, address string) (string, *systemmodels.Error) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	coll := client.Database(database.CurrentDatabase).Collection(database.USER)
	log.Info("Password: " + user.Password)
	found := authorizationOk(user.Email, user.Clone().Password, coll)

	if found == nil {
		return "", &systemmodels.Error{
			Status:  http.HTTP_STATUS_FORBIDDEN,
			Error:   valerror.USER_NOT_AUTHORIZED,
			Message: "Invalid credentials",
		}
	}

	device := &devicemodels.Device{Address: ip, UserAgent: address}
	device, err := devicedal.AddUserDevice(found, device)

	if err != nil {
		return "", &systemmodels.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   valerror.UNEXPECTED_ERROR,
			Message: "Cannot generate your auth token",
		}
	}

	return device.Token, nil
}

func LoginAuth(auth *usersmodels.AuthLogin, ip string, userAgent string) *systemmodels.Error {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	found, err := GetUser(&usersmodels.User{Email: auth.Email}, false)

	if err != nil {
		return err
	}

	// Search a user device with the same ip and user agent that has the token
	var filter = devicemodels.Device{
		User:      found.Email,
		UserAgent: userAgent,
		Address:   ip,
		Token:     auth.AuthToken,
	}

	var devices = client.Database(database.CurrentDatabase).Collection(database.DEVICE)
	device, deviceFindingError := devicedal.FindDeviceByAuthToken(devices, &filter)

	if deviceFindingError != nil || device == nil {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_NOT_FOUND,
			Error:   valerror.USER_NOT_AUTHORIZED,
			Message: "No possible login devices",
		}
	}

	return nil
}

func EditUser(user *usersmodels.User) *systemmodels.Error {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	users := client.Database(database.CurrentDatabase).Collection(database.USER)

	// validate email
	if user.Email != "" {
		checkedPass := utils.ValidateEmail(user.Email)

		if checkedPass.Response != 200 {
			return &systemmodels.Error{
				Status:  http.HTTP_STATUS_BAD_REQUEST,
				Error:   checkedPass.Response,
				Message: checkedPass.Message,
			}
		}
	}

	// validate password
	if user.Password != "" {

		checkedPass := utils.ValidatePassword(user.Password)

		if checkedPass.Response != 200 {
			return &systemmodels.Error{
				Status:  http.HTTP_STATUS_BAD_REQUEST,
				Error:   checkedPass.Response,
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
	res, err := users.UpdateOne(database.GetDefaultContext(), bson.M{"email": user.Email}, toUpdate)

	if err != nil {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   valerror.USER_NOT_UPDATED,
			Message: "User not updated",
		}
	}

	if res.MatchedCount == 0 && res.ModifiedCount == 0 {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_NOT_FOUND,
			Error:   valerror.USER_NOT_FOUND,
			Message: "Users not found",
		}
	}

	return nil
}

func EditUserEmail(mail *EmailChangeRequest) *systemmodels.Error {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	if utils.IsEmpty(mail.Email) || utils.IsEmpty(mail.NewEmail) {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.EMPTY_USER_EMAIL,
			Message: "Email cannot be empty",
		}
	}

	// Equal emails
	if mail.Email == mail.NewEmail {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.USER_EDITING_EMAILS_EQUAL,
			Message: "The new email is the same as the old one",
		}
	}

	// validate email
	var checkedPass = utils.ValidateEmail(mail.Email)
	if checkedPass.Response != 200 {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   checkedPass.Response,
			Message: checkedPass.Message,
		}
	}

	// Check if user exists
	users := client.Database(database.CurrentDatabase).Collection(database.USER)
	found := mailExists(mail.NewEmail, users)

	if found != nil {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_CONFLICT,
			Error:   valerror.USER_ALREADY_EXISTS,
			Message: "That email is already in use",
		}
	}

	// update user on database
	var checkedEmail = utils.ValidateEmail(mail.NewEmail)
	if checkedEmail.Response != 200 {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   int(checkedEmail.Response),
			Message: checkedEmail.Message,
		}

	}

	updateStatus, err := users.UpdateOne(database.GetDefaultContext(), bson.M{"email": mail.Email}, bson.M{"$set": bson.M{"email": mail.NewEmail}})

	if err != nil {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   valerror.USER_NOT_UPDATED,
			Message: "User not updated" + err.Error(),
		}
	}

	if updateStatus.MatchedCount == 0 {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_NOT_FOUND,
			Error:   valerror.USER_NOT_FOUND,
			Message: "User not found",
		}
	}

	if updateStatus.ModifiedCount == 0 {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   valerror.USER_NOT_UPDATED,
			Message: "User not updated",
		}
	}

	// update user devices on database
	devices := client.Database(database.CurrentDatabase).Collection(database.DEVICE)

	updateStatus, err = devices.UpdateMany(database.GetDefaultContext(), bson.M{"user": mail.Email}, bson.M{"$set": bson.M{"user": mail.NewEmail}})

	if err != nil {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   valerror.USER_NOT_UPDATED,
			Message: "User devices not updated",
		}
	}

	if updateStatus.MatchedCount != 0 && updateStatus.ModifiedCount == 0 {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   valerror.USER_NOT_UPDATED,
			Message: "User devices not updated",
		}
	}

	return nil
}

func EditUserProfilePicture(user *usersmodels.User, picture []byte) *systemmodels.Error {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	if utils.IsEmpty(user.Email) {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.EMPTY_USER_EMAIL,
			Message: "Email cannot be empty",
		}
	}

	//if the path base profile pic does not exist, create it
	var profilePathDir = utils.GetProfilePicturePath("", configuration.PROFILE_PICTURES_PATH)
	if !utils.ExistsDir(profilePathDir) {
		err := utils.CreateDir(profilePathDir)

		if err != nil {
			return &systemmodels.Error{
				Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
				Error:   valerror.USER_NOT_UPDATED,
				Message: "User not updated, image not saved :" + err.Error(),
			}
		}
	}

	// save the picture
	var profilePicPath = utils.GetProfilePicturePath(user.Email, configuration.PROFILE_PICTURES_PATH)
	err := utils.SaveFile(profilePicPath, picture)

	if err != nil {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   valerror.USER_NOT_UPDATED,
			Message: "User not updated, image not saved :" + err.Error(),
		}
	}

	user.ProfilePic = profilePicPath
	editErr := EditUser(user)

	if editErr != nil {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   valerror.USER_NOT_UPDATED,
			Message: "User not updated",
		}
	}

	return nil
}

func DeleteUser(user *usersmodels.User) *systemmodels.Error {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	if utils.IsEmpty(user.Email) {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.EMPTY_USER_EMAIL,
			Message: "Email cannot be empty",
		}
	}

	// delete user projects
	projects := client.Database(database.CurrentDatabase).Collection(database.PROJECT)
	_, err := projects.DeleteMany(database.GetDefaultContext(), bson.M{"owner": user.Email})

	if err != nil {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   valerror.USER_NOT_DELETED,
			Message: "User not deleted",
		}
	}

	// delete user devices
	devices := client.Database(database.CurrentDatabase).Collection(database.DEVICE)
	_, err = devices.DeleteMany(database.GetDefaultContext(), bson.M{"user": user.Email})

	if err != nil {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   valerror.USER_NOT_DELETED,
			Message: "User not deleted",
		}
	}

	// delete user on database
	users := client.Database(database.CurrentDatabase).Collection(database.USER)

	var deleteResult *mongo.DeleteResult
	deleteResult, err = users.DeleteOne(database.GetDefaultContext(), bson.M{"email": user.Email})

	if err != nil {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   valerror.USER_NOT_DELETED,
			Message: "User not deleted",
		}
	}

	if deleteResult.DeletedCount == 0 {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_NOT_FOUND,
			Error:   valerror.USER_NOT_FOUND,
			Message: "User not found",
		}
	}

	return nil
}

func GetUser(user *usersmodels.User, secure bool) (*usersmodels.User, *systemmodels.Error) { // get user from database

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	users := client.Database(database.CurrentDatabase).Collection(database.USER)
	var found usersmodels.User
	err := users.FindOne(database.GetDefaultContext(), bson.M{"email": user.Email}).Decode(&found)

	if err != nil {
		return nil, &systemmodels.Error{
			Status:  http.HTTP_STATUS_NOT_FOUND,
			Error:   valerror.USER_NOT_FOUND,
			Message: "User not found",
		}
	}

	if secure {
		found.Password = "****************"
	}

	return &found, nil
}

func ValidateUser(code string) *systemmodels.Error {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	if utils.IsEmpty(code) {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.INVALID_VALIDATION_CODE,
			Message: "Code cannot be empty",
		}
	}

	var user = &usersmodels.User{
		ValidationCode: code,
	}
	coll := client.Database(database.CurrentDatabase).Collection(database.USER)
	err := coll.FindOne(database.GetDefaultContext(), user).Decode(user)

	log.FormattedInfo("user: ${0}", user.Email)
	log.FormattedInfo("code: ${0}", code)

	if err != nil {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.INVALID_VALIDATION_CODE,
			Message: "Invalid validation code",
		}
	}

	if user.Validated {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_OK,
			Error:   valerror.USER_ALREADY_VALIDATED,
			Message: "User already validated",
		}
	}

	if user.ValidationCode != code {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_BAD_REQUEST,
			Error:   valerror.INVALID_VALIDATION_CODE,
			Message: "Invalid validation code",
		}
	}

	user.ValidationCode = ""
	user.Validated = true

	// update user on database
	result, editerr := coll.UpdateOne(database.GetDefaultContext(), bson.M{"email": user.Email}, bson.M{"$set": bson.M{"validation_code": "", "validated": true}})

	if result.MatchedCount == 0 {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_NOT_FOUND,
			Error:   valerror.USER_NOT_FOUND,
			Message: "User not found",
		}
	}

	if editerr != nil {
		return &systemmodels.Error{
			Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
			Error:   valerror.USER_NOT_UPDATED,
			Message: "User not validated: " + editerr.Error(),
		}
	}

	return nil
}

func mailExists(email string, coll *mongo.Collection) *usersmodels.User {

	filter := bson.D{{Key: "email", Value: email}}

	var result usersmodels.User
	err := coll.FindOne(database.GetDefaultContext(), filter).Decode(&result)

	if err != nil {
		log.FormattedError("Error: ${0}", err.Error())
		return nil
	}

	return &result
}

func authorizationOk(email string, password string, coll *mongo.Collection) *usersmodels.User {

	filter := bson.D{
		{Key: "email", Value: email},
		{Key: "password", Value: utils.EncryptSha256(password)},
	}

	var result usersmodels.User
	err := coll.FindOne(database.GetDefaultContext(), filter).Decode(&result)

	if err != nil {
		log.FormattedError("Error: ${0}", err.Error())
		return nil
	}

	return &result
}

func getUserFromToken(token string) (usersmodels.User, *systemmodels.Error) {

	// Connect database
	var client = database.Connect()
	defer database.Disconnect(*client)

	var tokenDevice devicemodels.Device

	devices := client.Database(database.CurrentDatabase).Collection(database.DEVICE)
	err := devices.FindOne(database.GetDefaultContext(), bson.M{"token": token}).Decode(&tokenDevice)

	if err != nil {
		return usersmodels.User{}, &systemmodels.Error{
			Status:  http.HTTP_STATUS_FORBIDDEN,
			Error:   valerror.INVALID_TOKEN,
			Message: "User not matching token",
		}
	}

	var tokenUser usersmodels.User

	users := client.Database(database.CurrentDatabase).Collection(database.USER)
	err = users.FindOne(database.GetDefaultContext(), bson.M{"email": tokenDevice.User}).Decode(&tokenUser)

	if err != nil {
		return usersmodels.User{}, &systemmodels.Error{
			Status:  http.HTTP_STATUS_FORBIDDEN,
			Error:   valerror.INVALID_TOKEN,
			Message: "User not matching token",
		}
	}

	return tokenUser, nil
}

func IsTokenValid(token string) (*usersmodels.User, *systemmodels.Error) {

	// decode token
	claims, err := utils.DecryptToken(token, configuration.Params.Secret)

	if err != nil {
		return nil, &systemmodels.Error{
			Status:  http.HTTP_STATUS_FORBIDDEN,
			Error:   valerror.INVALID_TOKEN,
			Message: "invalid token format",
		}
	}

	// log token claims
	log.Info("device: " + claims.Claims.(jwt.MapClaims)["device"].(string))
	log.Info("username: " + claims.Claims.(jwt.MapClaims)["username"].(string))
	log.Info("email: " + claims.Claims.(jwt.MapClaims)["email"].(string))

	email := claims.Claims.(jwt.MapClaims)["email"].(string)

	foundUser, tokenUserErr := getUserFromToken(token)

	if tokenUserErr != nil {
		return nil, &systemmodels.Error{
			Status:  http.HTTP_STATUS_FORBIDDEN,
			Error:   valerror.INVALID_TOKEN,
			Message: "invalid token",
		}
	}

	if foundUser.Email != email {
		return nil, &systemmodels.Error{
			Status:  http.HTTP_STATUS_FORBIDDEN,
			Error:   valerror.INVALID_TOKEN,
			Message: "invalid token",
		}
	}

	return &foundUser, nil
}
