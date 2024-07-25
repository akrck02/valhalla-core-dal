package userdal

import (
	"net/http"

	"github.com/akrck02/valhalla-core-dal/configuration"
	"github.com/akrck02/valhalla-core-dal/database"
	devicedal "github.com/akrck02/valhalla-core-dal/services/device"

	apierror "github.com/akrck02/valhalla-core-sdk/error"
	"github.com/akrck02/valhalla-core-sdk/log"
	apimodels "github.com/akrck02/valhalla-core-sdk/models/api"
	devicemodels "github.com/akrck02/valhalla-core-sdk/models/device"
	usersmodels "github.com/akrck02/valhalla-core-sdk/models/users"
	"github.com/akrck02/valhalla-core-sdk/utils"

	"github.com/golang-jwt/jwt/v5"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EmailChangeRequest struct {
	Email    string `json:"email"`
	NewEmail string `json:"new_email"`
}

func Register(conn *mongo.Client, user *usersmodels.User) *apimodels.Error {

	if utils.IsEmpty(user.Email) {
		return &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   apierror.EmptyUserEmail,
			Message: "Email cannot be empty",
		}
	}

	if utils.IsEmpty(user.Password) {
		return &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   apierror.EmptyUserPassword,
			Message: "Password cannot be empty",
		}
	}

	if utils.IsEmpty(user.Username) {
		return &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   apierror.EmptyUsername,
			Message: "Username cannot be empty",
		}
	}

	var checkedError = utils.ValidatePassword(user.Password)
	if checkedError != nil {
		return checkedError
	}

	checkedError = utils.ValidateEmail(user.Email)
	if checkedError != nil {
		return checkedError
	}

	coll := conn.Database(database.CurrentDatabase).Collection(database.USER)
	found := mailExists(user.Email, coll)

	if found != nil {
		return &apimodels.Error{
			Status:  http.StatusConflict,
			Error:   apierror.UserAlreadyExists,
			Message: "User already exists",
		}
	}

	code, err := utils.GenerateValidationCode(user.Email)

	if err != nil {
		return &apimodels.Error{
			Status:  http.StatusInternalServerError,
			Error:   apierror.CannotCreateValidationCode,
			Message: "User not created",
		}
	}

	userToInsert := user.Clone()
	userToInsert.Password = utils.EncryptSha256(user.Clone().Password)
	userToInsert.ValidationCode = code

	creationDate := utils.GetCurrentMillis()
	userToInsert.CreationDate = &creationDate
	userToInsert.LastUpdate = userToInsert.CreationDate

	// register user on database
	res, errInsert := coll.InsertOne(database.GetDefaultContext(), userToInsert)

	if errInsert != nil {
		return &apimodels.Error{
			Status:  http.StatusConflict,
			Error:   apierror.UserAlreadyExists,
			Message: "User already exists",
		}
	}

	// get user from database
	user.ID = res.InsertedID.(primitive.ObjectID).Hex()
	user.CreationDate = userToInsert.CreationDate
	user.LastUpdate = userToInsert.LastUpdate
	user.ValidationCode = userToInsert.ValidationCode

	return nil
}

func Login(conn *mongo.Client, user *usersmodels.User, ip string, address string) (string, *apimodels.Error) {

	// Connect database
	coll := conn.Database(database.CurrentDatabase).Collection(database.USER)
	log.Info("Password: " + user.Password)
	found := authorizationOk(user.Email, user.Clone().Password, coll)

	if found == nil {
		return "", &apimodels.Error{
			Status:  http.StatusForbidden,
			Error:   apierror.UserNotAuthorized,
			Message: "Invalid credentials",
		}
	}

	device := &devicemodels.Device{Address: ip, UserAgent: address}
	device, err := devicedal.AddUserDevice(conn, found, device)

	if err != nil {
		return "", &apimodels.Error{
			Status:  http.StatusInternalServerError,
			Error:   apierror.UnexpectedError,
			Message: "Cannot generate your auth token",
		}
	}

	return device.Token, nil
}

func LoginAuth(conn *mongo.Client, auth *usersmodels.AuthLogin, ip string, userAgent string) *apimodels.Error {

	// Connect database

	found, err := GetUserByEmail(conn, auth.Email, false)

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

	var devices = conn.Database(database.CurrentDatabase).Collection(database.DEVICE)
	device, deviceFindingError := devicedal.FindDeviceByAuthToken(devices, &filter)

	if deviceFindingError != nil || device == nil {
		return &apimodels.Error{
			Status:  http.StatusNotFound,
			Error:   apierror.UserNotAuthorized,
			Message: "No possible login devices",
		}
	}

	return nil
}

func EditUser(conn *mongo.Client, user *usersmodels.User) *apimodels.Error {

	users := conn.Database(database.CurrentDatabase).Collection(database.USER)

	if utils.IsEmpty(user.ID) {
		return &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   apierror.EmptyUsername,
			Message: "User id cannot be empty",
		}
	}

	// Get if id is valid
	_, parseIdError := utils.StringToObjectId(user.ID)

	if parseIdError != nil {
		return &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   apierror.InvalidObjectId,
			Message: "Invalid user id",
		}
	}

	// Check if the user exists
	found := userExists(user.ID, users)

	if found == nil {
		return &apimodels.Error{
			Status:  http.StatusNotFound,
			Error:   apierror.UserNotFound,
			Message: "User not found",
		}
	}

	// validate email
	if user.Email != "" {
		checkedError := utils.ValidateEmail(user.Email)

		if checkedError != nil {
			return checkedError
		}
	}

	// validate password
	if user.Password != "" {

		checkedError := utils.ValidatePassword(user.Password)

		if checkedError != nil {
			return checkedError
		}
	}

	setObject := bson.M{}

	if user.Username != "" {
		setObject["username"] = user.Username
	}

	if user.Password != "" {
		encryptedPass := user.Password
		setObject["password"] = utils.EncryptSha256(encryptedPass)
	}

	if user.ProfilePic != "" {
		setObject["profilePic"] = user.ProfilePic
	}

	setObject["updatedate"] = utils.GetCurrentMillis()
	toUpdate := bson.M{"$set": setObject}

	// update user on database
	objID, parseObjectIdError := primitive.ObjectIDFromHex(user.ID)

	if parseObjectIdError != nil {
		return &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   apierror.InvalidObjectId,
			Message: "Invalid user id",
		}
	}

	res, err := users.UpdateByID(database.GetDefaultContext(), objID, toUpdate)

	if err != nil {
		return &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   apierror.DatabaseError,
			Message: "User not updated",
		}
	}

	if res.MatchedCount == 0 && res.ModifiedCount == 0 {
		return &apimodels.Error{
			Status:  http.StatusNotFound,
			Error:   apierror.UserNotFound,
			Message: "Users not found",
		}
	}

	return nil
}

func EditUserEmail(conn *mongo.Client, mail *EmailChangeRequest) *apimodels.Error {

	// Connect database
	if utils.IsEmpty(mail.Email) || utils.IsEmpty(mail.NewEmail) {
		return &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   apierror.EmptyUserEmail,
			Message: "Email cannot be empty",
		}
	}

	// Equal emails
	if mail.Email == mail.NewEmail {
		return &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   apierror.UpdateParametersEqual,
			Message: "The new email is the same as the old one",
		}
	}

	// validate email
	checkedEmailError := utils.ValidateEmail(mail.Email)
	if checkedEmailError != nil {
		return checkedEmailError
	}

	// Check if user exists
	users := conn.Database(database.CurrentDatabase).Collection(database.USER)
	found := mailExists(mail.NewEmail, users)

	if found != nil {
		return &apimodels.Error{
			Status:  http.StatusConflict,
			Error:   apierror.UserAlreadyExists,
			Message: "That email is already in use",
		}
	}

	// update user on database
	checkedEmailError = utils.ValidateEmail(mail.NewEmail)
	if checkedEmailError != nil {
		return checkedEmailError

	}

	updateStatus, err := users.UpdateOne(database.GetDefaultContext(),
		bson.M{"email": mail.Email},
		bson.M{"$set": bson.M{
			"email":      mail.NewEmail,
			"updatedate": utils.GetCurrentMillis(),
		}},
	)

	if err != nil {
		return &apimodels.Error{
			Status:  http.StatusInternalServerError,
			Error:   apierror.DatabaseError,
			Message: "User not updated" + err.Error(),
		}
	}

	if updateStatus.MatchedCount == 0 {
		return &apimodels.Error{
			Status:  http.StatusNotFound,
			Error:   apierror.UserNotFound,
			Message: "User not found",
		}
	}

	if updateStatus.ModifiedCount == 0 {
		return &apimodels.Error{
			Status:  http.StatusInternalServerError,
			Error:   apierror.DatabaseError,
			Message: "User not updated",
		}
	}

	// update user devices on database
	devices := conn.Database(database.CurrentDatabase).Collection(database.DEVICE)
	updateStatus, err = devices.UpdateMany(database.GetDefaultContext(), bson.M{"user": mail.Email}, bson.M{"$set": bson.M{"user": mail.NewEmail}})

	if err != nil {
		return &apimodels.Error{
			Status:  http.StatusInternalServerError,
			Error:   apierror.DatabaseError,
			Message: "User devices not updated",
		}
	}

	if updateStatus.MatchedCount != 0 && updateStatus.ModifiedCount == 0 {
		return &apimodels.Error{
			Status:  http.StatusInternalServerError,
			Error:   apierror.DatabaseError,
			Message: "User devices not updated",
		}
	}

	return nil
}

func EditUserProfilePicture(conn *mongo.Client, user *usersmodels.User, picture []byte) *apimodels.Error {

	if utils.IsEmpty(user.ID) {
		return &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   apierror.EmptyUsername,
			Message: "User id cannot be empty",
		}
	}

	// Get if id is valid
	_, parseIdError := utils.StringToObjectId(user.ID)

	if parseIdError != nil {
		return &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   apierror.InvalidObjectId,
			Message: "Invalid user id",
		}
	}

	//if the path base profile pic does not exist, create it
	var profilePathDir = utils.GetProfilePicturePath(user.ID, configuration.PROFILE_PICTURES_PATH)
	if !utils.ExistsDir(profilePathDir) {
		err := utils.CreateDir(profilePathDir)

		if err != nil {
			return &apimodels.Error{
				Status:  http.StatusInternalServerError,
				Error:   apierror.DatabaseError,
				Message: "User not updated, image not saved :" + err.Error(),
			}
		}
	}

	// save the picture
	var profilePicPath = utils.GetProfilePicturePath(user.Email, configuration.PROFILE_PICTURES_PATH)
	saveFileError := utils.SaveFile(profilePicPath, picture)

	if parseIdError != nil {
		return &apimodels.Error{
			Status:  http.StatusInternalServerError,
			Error:   apierror.DatabaseError,
			Message: "User not updated, image not saved :" + saveFileError.Error(),
		}
	}

	user.ProfilePic = profilePicPath
	editErr := EditUser(conn, user)

	if editErr != nil {
		return &apimodels.Error{
			Status:  http.StatusInternalServerError,
			Error:   apierror.DatabaseError,
			Message: "User not updated",
		}
	}

	return nil
}

func DeleteUser(conn *mongo.Client, user *usersmodels.User) *apimodels.Error {

	// Connect database

	if utils.IsEmpty(user.Email) {
		return &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   apierror.EmptyUserEmail,
			Message: "Email cannot be empty",
		}
	}

	// delete user projects
	projects := conn.Database(database.CurrentDatabase).Collection(database.PROJECT)
	_, err := projects.DeleteMany(database.GetDefaultContext(), bson.M{"owner": user.Email})

	if err != nil {
		return &apimodels.Error{
			Status:  http.StatusInternalServerError,
			Error:   apierror.UserNotDeleted,
			Message: "User not deleted",
		}
	}

	// delete user devices
	devices := conn.Database(database.CurrentDatabase).Collection(database.DEVICE)
	_, err = devices.DeleteMany(database.GetDefaultContext(), bson.M{"user": user.Email})

	if err != nil {
		return &apimodels.Error{
			Status:  http.StatusInternalServerError,
			Error:   apierror.UserNotDeleted,
			Message: "User not deleted",
		}
	}

	// delete user on database
	users := conn.Database(database.CurrentDatabase).Collection(database.USER)

	var deleteResult *mongo.DeleteResult
	deleteResult, err = users.DeleteOne(database.GetDefaultContext(), bson.M{"email": user.Email})

	if err != nil {
		return &apimodels.Error{
			Status:  http.StatusInternalServerError,
			Error:   apierror.UserNotDeleted,
			Message: "User not deleted",
		}
	}

	if deleteResult.DeletedCount == 0 {
		return &apimodels.Error{
			Status:  http.StatusNotFound,
			Error:   apierror.UserNotFound,
			Message: "User not found",
		}
	}

	return nil
}

func GetUser(conn *mongo.Client, user *usersmodels.User, secure bool) (*usersmodels.User, *apimodels.Error) {

	id, err := utils.StringToObjectId(user.ID)

	if err != nil {
		return nil, &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   apierror.InvalidObjectId,
			Message: "Invalid user id",
		}
	}

	filter := &bson.M{"_id": id}

	// get user from database
	users := conn.Database(database.CurrentDatabase).Collection(database.USER)
	var found usersmodels.User
	err = users.FindOne(database.GetDefaultContext(), filter).Decode(&found)

	// if an error occurs,
	if err != nil {
		return nil, &apimodels.Error{
			Status:  http.StatusInternalServerError,
			Error:   apierror.UnexpectedError,
			Message: err.Error(),
		}
	}

	// if user not found, return error
	if found.ID == "" {
		return nil, &apimodels.Error{
			Status:  http.StatusNotFound,
			Error:   apierror.UserNotFound,
			Message: "User not found",
		}
	}

	// if is secure search, hide password
	if secure {
		found.Password = "****************"
	}

	// return user
	return &found, nil
}

func GetUserByEmail(conn *mongo.Client, email string, secure bool) (*usersmodels.User, *apimodels.Error) {

	filter := bson.M{"email": email}

	// get user from database
	users := conn.Database(database.CurrentDatabase).Collection(database.USER)
	var found usersmodels.User
	err := users.FindOne(database.GetDefaultContext(), filter).Decode(&found)

	// if an error occurs,
	if err != nil {
		return nil, &apimodels.Error{
			Status:  http.StatusInternalServerError,
			Error:   apierror.UnexpectedError,
			Message: err.Error(),
		}
	}

	// if user not found, return error
	if found.ID == "" {
		return nil, &apimodels.Error{
			Status:  http.StatusNotFound,
			Error:   apierror.UserNotFound,
			Message: "User not found",
		}
	}

	// if is secure search, hide password
	if secure {
		found.Password = "****************"
	}

	// return user
	return &found, nil
}

func ValidateUser(conn *mongo.Client, code string) *apimodels.Error {

	// Connect database

	if utils.IsEmpty(code) {
		return &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   apierror.InvalidValidationCode,
			Message: "Code cannot be empty",
		}
	}

	var user = &usersmodels.User{
		ValidationCode: code,
	}
	coll := conn.Database(database.CurrentDatabase).Collection(database.USER)
	err := coll.FindOne(database.GetDefaultContext(), user).Decode(user)

	log.FormattedInfo("user: ${0}", user.Email)
	log.FormattedInfo("code: ${0}", code)

	if err != nil {
		return &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   apierror.InvalidValidationCode,
			Message: "Invalid validation code",
		}
	}

	if user.Validated {
		return &apimodels.Error{
			Status:  http.StatusOK,
			Error:   apierror.UserAlreadyValidated,
			Message: "User already validated",
		}
	}

	if user.ValidationCode != code {
		return &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   apierror.InvalidValidationCode,
			Message: "Invalid validation code",
		}
	}

	user.ValidationCode = ""
	user.Validated = true

	// update user on database
	result, editerr := coll.UpdateOne(database.GetDefaultContext(), bson.M{"email": user.Email}, bson.M{"$set": bson.M{"validation_code": "", "validated": true}})

	if result.MatchedCount == 0 {
		return &apimodels.Error{
			Status:  http.StatusNotFound,
			Error:   apierror.UserNotFound,
			Message: "User not found",
		}
	}

	if editerr != nil {
		return &apimodels.Error{
			Status:  http.StatusInternalServerError,
			Error:   apierror.DatabaseError,
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
		return nil
	}

	return &result
}

func userExists(userId string, coll *mongo.Collection) *usersmodels.User {

	objectId, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return nil
	}

	filter := bson.D{{Key: "_id", Value: objectId}}

	var result usersmodels.User
	err = coll.FindOne(database.GetDefaultContext(), filter).Decode(&result)

	if err != nil {
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
		return nil
	}

	return &result
}

func getUserFromToken(conn *mongo.Client, token string) (usersmodels.User, *apimodels.Error) {

	// Connect database

	var tokenDevice devicemodels.Device

	devices := conn.Database(database.CurrentDatabase).Collection(database.DEVICE)
	err := devices.FindOne(database.GetDefaultContext(), bson.M{"token": token}).Decode(&tokenDevice)

	if err != nil {
		return usersmodels.User{}, &apimodels.Error{
			Status:  http.StatusForbidden,
			Error:   apierror.InvalidToken,
			Message: "User not matching token",
		}
	}

	var tokenUser usersmodels.User

	users := conn.Database(database.CurrentDatabase).Collection(database.USER)
	err = users.FindOne(database.GetDefaultContext(), bson.M{"email": tokenDevice.User}).Decode(&tokenUser)

	if err != nil {
		return usersmodels.User{}, &apimodels.Error{
			Status:  http.StatusForbidden,
			Error:   apierror.InvalidToken,
			Message: "User not matching token",
		}
	}

	return tokenUser, nil
}

func IsTokenValid(conn *mongo.Client, token string) (*usersmodels.User, *apimodels.Error) {

	// decode token
	claims, err := utils.DecryptToken(token, configuration.Params.Secret)

	if err != nil {
		return nil, &apimodels.Error{
			Status:  http.StatusForbidden,
			Error:   apierror.InvalidToken,
			Message: "invalid token format",
		}
	}

	// log token claims
	log.Info("device: " + claims.Claims.(jwt.MapClaims)["device"].(string))
	log.Info("username: " + claims.Claims.(jwt.MapClaims)["username"].(string))
	log.Info("email: " + claims.Claims.(jwt.MapClaims)["email"].(string))

	email := claims.Claims.(jwt.MapClaims)["email"].(string)

	foundUser, tokenUserErr := getUserFromToken(conn, token)

	if tokenUserErr != nil {
		return nil, &apimodels.Error{
			Status:  http.StatusForbidden,
			Error:   apierror.InvalidToken,
			Message: "invalid token",
		}
	}

	if foundUser.Email != email {
		return nil, &apimodels.Error{
			Status:  http.StatusForbidden,
			Error:   apierror.InvalidToken,
			Message: "invalid token",
		}
	}

	return &foundUser, nil
}
