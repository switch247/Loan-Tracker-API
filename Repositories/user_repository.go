package Repositories

import (
	"Loan-Tracker-API/Domain"
	"Loan-Tracker-API/Dtos"
	emailservice "Loan-Tracker-API/Infrastructure/email_service"
	jwtservice "Loan-Tracker-API/Infrastructure/jwt_service"
	"Loan-Tracker-API/Infrastructure/password_services"
	ps "Loan-Tracker-API/Infrastructure/password_services"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "golang.org/x/oauth2"
	// "golang.org/x/oauth2/google"
)

type UserRepository struct {
	validator          *validator.Validate
	UserCollection     Domain.Collection
	mu                 sync.RWMutex
	refresh_collection Domain.Collection
	emailservice       *emailservice.MailService
}

func NewUserRepository(user_collection Domain.Collection, token_collection Domain.Collection, _valdator *validator.Validate, email_service_reference *emailservice.MailService) *UserRepository {

	return &UserRepository{
		validator:          _valdator,
		UserCollection:     user_collection,
		mu:                 sync.RWMutex{},
		refresh_collection: token_collection,
		emailservice:       email_service_reference,
		// oauth2Config:    *oauth_config,
	}
}

// create user
func (ur *UserRepository) CreateUser(ctx context.Context, user *Domain.User) (Dtos.OmitedUser, error, int) {
	ur.mu.RLock()
	defer ur.mu.RUnlock()

	// Check if user email is taken
	existingUserFilter := bson.D{{"email", user.Email}}
	existingUserCount, err := ur.UserCollection.CountDocuments(ctx, existingUserFilter)
	if err != nil {
		return Dtos.OmitedUser{}, err, 500
	}
	if existingUserCount > 0 {
		return Dtos.OmitedUser{}, errors.New("Email is already taken"), http.StatusBadRequest
	}
	// User registration logic
	hashedPassword, err := ps.GenerateFromPasswordCustom(user.Password)
	if err != nil {
		return Dtos.OmitedUser{}, err, 500
	}
	user.Password = string(hashedPassword)
	if user.Role == "" {
		user.Role = "user"
	}
	if user.UserName == "" {
		user.UserName = user.Email + "_" + user.Role
	}
	if user.EmailVerified == false {
		user.EmailVerified = true
	}
	fmt.Println(user)
	insertResult, err := ur.UserCollection.InsertOne(ctx, user)
	if err != nil {
		return Dtos.OmitedUser{}, err, 500
	}
	// Fetch the inserted task
	var fetched Dtos.OmitedUser
	err = ur.UserCollection.FindOne(context.TODO(), bson.D{{"_id", insertResult.InsertedID.(primitive.ObjectID)}}).Decode(&fetched)
	if err != nil {
		fmt.Println(err)
		return Dtos.OmitedUser{}, errors.New("User Not Created"), 500
	}
	if fetched.Email != user.Email {
		return Dtos.OmitedUser{}, errors.New("User Not Created"), 500
	}
	fetched.Password = ""
	return fetched, nil, 200
}

// get all users
func (ur *UserRepository) GetUsers(ctx context.Context) ([]*Dtos.OmitedUser, error, int) {
	ur.mu.RLock()
	defer ur.mu.RUnlock()
	var results []*Dtos.OmitedUser

	// Pass these options to the Find method
	findOptions := options.Find()
	// findOptions.SetLimit(2)
	filter := bson.D{{}}

	// Here's an array in which you can store the decoded documents

	// Passing bson.D{{}} ur the filter matches all documents in the collection
	cur, err := ur.UserCollection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Fatal("error in finding users", err)
		log.Fatal(err)
		return []*Dtos.OmitedUser{}, err, 0
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows ur to decode documents one at a time
	for cur.Next(ctx) {

		// create a value into which the single document can be decoded
		var elem Dtos.OmitedUser
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println("error in decoding user", err)
			fmt.Println(err.Error())
			// #handelthislater
			// should this say there was a decoding error and return?
			return []*Dtos.OmitedUser{}, err, 500
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		fmt.Println(err)
		return []*Dtos.OmitedUser{}, err, 500
	}

	// Close the cursor once finished
	cur.Close(ctx)

	fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)
	return results, nil, 200
}

// get user by id
func (ur *UserRepository) GetUsersById(ctx context.Context, id primitive.ObjectID, current_user Domain.AccessClaims) (Dtos.OmitedUser, error, int) {
	ur.mu.RLock()
	defer ur.mu.RUnlock()

	var filter bson.D
	filter = bson.D{{"_id", id}}
	var result Dtos.OmitedUser
	err := ur.UserCollection.FindOne(ctx, filter).Decode(&result)
	// # handel this later
	if err != nil {
		return Dtos.OmitedUser{}, errors.New("User not found"), http.StatusNotFound
	}
	if current_user.Role == "user" && result.ID != current_user.ID {
		return Dtos.OmitedUser{}, errors.New("permission denied"), http.StatusForbidden

	}
	return result, nil, 200
}

// delete user by id
func (ur *UserRepository) DeleteUsersById(ctx context.Context, id primitive.ObjectID, current_user Domain.AccessClaims) (error, int) {
	ur.mu.RLock()
	defer ur.mu.RUnlock()

	filter := bson.D{{"_id", id}}
	if current_user.Role == "user" {
		return errors.New("permission denied"), http.StatusForbidden
	}

	deleteResult, err := ur.UserCollection.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return err, 500
	}
	if deleteResult.DeletedCount == 0 {
		return errors.New("User does not exist"), http.StatusNotFound
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
	if current_user.ID == id {
		// delete the refresh token if you are deleting you self
		err, statusCode := ur.DeleteToken(ctx, id)
		if err != nil {
			return err, statusCode
		}
	}
	return nil, 200
}

func (ur *UserRepository) ChangePassByEmail(ctx context.Context, email string, password string) (Dtos.OmitedUser, error, int) {
	ur.mu.RLock()
	defer ur.mu.RUnlock()

	statusCode := 200
	filter := bson.D{{"email", email}}
	update := bson.D{
		{"$set", bson.D{
			{"password", password},
		}},
	}
	updateResult, err := ur.UserCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		statusCode = 500
		return Dtos.OmitedUser{}, err, statusCode
	}
	if updateResult.ModifiedCount == 0 {
		statusCode = 400
		fmt.Println("user does not exist:", email)
		return Dtos.OmitedUser{}, errors.New("user does not exist"), statusCode
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	return Dtos.OmitedUser{}, nil, statusCode
}

// find by email
func (ur *UserRepository) FindByEmail(ctx context.Context, email string) (Dtos.OmitedUser, error, int) {
	ur.mu.RLock()
	defer ur.mu.RUnlock()
	filter := bson.D{{"email", email}}
	var result Dtos.OmitedUser
	err := ur.UserCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return Dtos.OmitedUser{}, err, 500
	}
	return result, nil, 200
}

// login.
func (ur *UserRepository) Login(ctx context.Context, user *Dtos.LoginUserDto) (Domain.Tokens, error, int) {
	ur.mu.RLock()
	defer ur.mu.RUnlock()

	filter := bson.D{
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "username", Value: user.UserName}},
			bson.D{{Key: "email", Value: user.Email}},
		}},
	}
	var existingUser Domain.User
	err := ur.UserCollection.FindOne(ctx, filter).Decode(&existingUser)
	fmt.Print("existingUser", existingUser)

	if err != nil || !password_services.CompareHashAndPasswordCustom(existingUser.Password, user.Password) {
		fmt.Printf("Login Called:%v, %v", existingUser.Password, user.Password)

		// cpmpare the hashed password
		hashedPassword, _ := password_services.GenerateFromPasswordCustom(user.Password)
		fmt.Print(existingUser.Password == hashedPassword)
		return Domain.Tokens{}, errors.New("Invalid credentials"), http.StatusBadRequest
	}
	fmt.Println("emailverified :", existingUser.EmailVerified, "email", existingUser.Email)
	if existingUser.EmailVerified == false {
		err, statusCode := ur.SendActivationEmail(user.Email)
		if err != nil {
			fmt.Println("error at sending email", err)
			return Domain.Tokens{}, err, statusCode
		}
		return Domain.Tokens{}, errors.New("email is not activated , an activation email has been sent"), http.StatusUnauthorized
	}
	return ur.GenerateTokenFromUser(ctx, existingUser)

}

// register
func (ur *UserRepository) Register(ctx context.Context, user *Dtos.RegisterUserDto) (*Dtos.OmitedUser, error, int) {
	ur.mu.RLock()
	defer ur.mu.RUnlock()

	// Validate the user input
	err := ur.validator.Struct(user)
	if err != nil {
		return &Dtos.OmitedUser{}, err, http.StatusBadRequest
	}
	// Check if the email is already taken
	// existingUserFilter := bson.D{}
	// if user.UserName != "" {
	// 	existingUserFilter = bson.D{
	// 		{"$or", bson.A{
	// 			bson.D{{Key: "email", Value: user.Email}},
	// 			bson.D{{Key: "username", Value: user.UserName}},
	// 		}},
	// 	}
	// } else {
	// 	existingUserFilter = bson.D{
	// 		{Key: "email", Value: user.Email},
	// 	}
	// }
	// existingUserCount, err := ur.UserCollection.CountDocuments(ctx, existingUserFilter)
	// if err != nil {
	// 	fmt.Println("error at count", err)
	// 	return &Dtos.OmitedUser{}, err, 500
	// }
	// if existingUserCount > 0 {
	// 	return &Dtos.OmitedUser{}, errors.New("Email or UserName is already taken"), http.StatusBadRequest
	// }
	// check if password is following the rules
	err = password_services.CheckPasswordStrength(user.Password)
	if err != nil {
		return &Dtos.OmitedUser{}, err, http.StatusBadRequest
	}
	// User registration logic
	hashedPassword, err := password_services.GenerateFromPasswordCustom(user.Password)
	if err != nil {
		fmt.Println("error at hashing", err)
		return &Dtos.OmitedUser{}, err, 500
	}
	user.EmailVerified = false
	user.Password = string(hashedPassword)
	if user.UserName == "" {
		user.UserName = user.Email + "_user"
	}
	user.Role = "user"
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	InsertedID, err := ur.UserCollection.InsertOne(ctx, user)
	if err != nil {
		if strings.Contains(err.Error(), "E11000 duplicate key error") {
			fmt.Println("error at insert", err)
			return &Dtos.OmitedUser{}, errors.New("Email Already Taken"), 400
		}
		return &Dtos.OmitedUser{}, err, 500
	}

	// Fetch the inserted task
	var fetched Dtos.OmitedUser

	// Access the InsertedID field from the InsertOneResult struct
	insertedID := InsertedID.InsertedID.(primitive.ObjectID)

	err = ur.UserCollection.FindOne(context.TODO(), bson.D{{"_id", insertedID}}).Decode(&fetched)
	if err != nil {
		fmt.Println(err)
		return &Dtos.OmitedUser{}, errors.New("User Not Created"), 500
	}
	if fetched.Email != user.Email {
		return &Dtos.OmitedUser{}, errors.New("User Not Created"), 500
	}
	fetched.Password = ""
	err, statusCode := ur.SendActivationEmail(fetched.Email)
	if err != nil {
		// clean up
		filter := bson.D{{"_id", fetched.ID}}
		deleteResult, err2 := ur.UserCollection.DeleteOne(ctx, filter)
		if err2 != nil || deleteResult.DeletedCount == 0 {
			fmt.Println("error at cleanup", err2)
			return &Dtos.OmitedUser{}, errors.New("error clean up"), 500
		} else {
			fmt.Println(err)
			return &Dtos.OmitedUser{}, err, statusCode
		}
	}
	return &fetched, err, 200
}

// logout
func (ur *UserRepository) Logout(ctx context.Context, user_id primitive.ObjectID) (error, int) {
	ur.mu.RLock()
	defer ur.mu.RUnlock()

	// delete the refresh token
	err, statusCode := ur.DeleteToken(ctx, user_id)
	if err != nil {
		return err, statusCode
	}
	return nil, 200
}

func (ur *UserRepository) ForgetPassword(ctx context.Context, email string) (error, int) {
	ur.mu.RLock()
	defer ur.mu.RUnlock()
	_, err, status := ur.FindByEmail(ctx, email)
	if err != nil {
		return err, status
	}
	resetToken, err := jwtservice.GenerateToken(email)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	err = ur.emailservice.SendPasswordResetEmail(email, resetToken)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	return nil, http.StatusOK

}

func (ur *UserRepository) ResetPassword(ctx context.Context, email string, password string, resetToken string) (error, int) {
	ur.mu.RLock()
	defer ur.mu.RUnlock()
	_, err := jwtservice.VerifyToken(resetToken)
	if err != nil {
		return err, http.StatusBadRequest
	}
	if password == "" {
		return errors.New("password is required"), http.StatusBadRequest
	}
	err = password_services.CheckPasswordStrength(password)
	if err != nil {
		return err, http.StatusBadRequest
	}
	hashed, err := password_services.GenerateFromPasswordCustom(password)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	_, err, _ = ur.ChangePassByEmail(ctx, email, hashed)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	fmt.Println("password:", password, "reset_token", resetToken)
	return nil, http.StatusOK

}

func (ur *UserRepository) GenerateTokenFromUser(ctx context.Context, existingUser Domain.User) (Domain.Tokens, error, int) {
	ur.mu.RLock()
	defer ur.mu.RUnlock()

	// filter := bson.D{{Key: "email", Value: existingUser.Email}}
	// Generate JWT access
	jwtAccessToken, err := jwtservice.CreateAccessToken(existingUser)
	if err != nil {
		return Domain.Tokens{}, err, 500
	}
	refreshToken, err := jwtservice.CreateRefreshToken(existingUser)
	if err != nil {
		return Domain.Tokens{}, err, 500
	}

	// filter := primitive.D{{"_id", existingUser.ID}}
	existingToken, err, statusCode := ur.FindToken(ctx, existingUser.ID)
	if err != nil && err.Error() != "mongo: no documents in result" {
		fmt.Println("error at count", err)
		return Domain.Tokens{}, err, statusCode
	}

	if existingToken != "" {
		// update the refresh token
		err, statusCode := ur.UpdateToken(ctx, refreshToken, existingUser.ID)
		if err != nil {
			return Domain.Tokens{}, err, statusCode
		}

	} else {
		err, statusCode := ur.StoreToken(ctx, existingUser.ID, refreshToken)
		if err != nil {
			return Domain.Tokens{}, err, statusCode
		}
	}

	return Domain.Tokens{
		AccessToken:  jwtAccessToken,
		RefreshToken: refreshToken,
	}, nil, 200
}

func (ur *UserRepository) ActivateAccount(ctx context.Context, token string) (error, int) {
	ur.mu.RLock()
	defer ur.mu.RUnlock()
	email, err := jwtservice.VerifyToken(token)
	if err != nil {
		return err, http.StatusBadRequest
	}
	fmt.Println("email:", email, "token:", token)

	filter := bson.D{{"email", email}}

	update := bson.D{
		{"$set", bson.D{
			{"email_verified", true},
			{"created_at", time.Now()},
		}},
	}
	UpdatedResult, err := ur.UserCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	if UpdatedResult.ModifiedCount == 0 {
		return errors.New("user does not exist"), 400
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", UpdatedResult.MatchedCount, UpdatedResult.ModifiedCount)

	return nil, http.StatusOK

}

func (ur *UserRepository) SendActivationEmail(email string) (error, int) {
	ur.mu.RLock()
	defer ur.mu.RUnlock()

	activationToken, err := jwtservice.GenerateToken(email)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	err = ur.emailservice.SendActivationEmail(email, activationToken)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}

// profile repository

// get user by id
func (ps *UserRepository) GetProfile(ctx context.Context, id primitive.ObjectID, current_user Domain.AccessClaims) (Dtos.OmitedUser, error, int) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	if current_user.ID != id {
		return Dtos.OmitedUser{}, errors.New("permission denied"), http.StatusForbidden
	}

	var filter bson.D
	filter = bson.D{{"_id", id}}
	var result Dtos.OmitedUser
	err := ps.UserCollection.FindOne(ctx, filter).Decode(&result)
	// # handel this later
	if err != nil {
		return Dtos.OmitedUser{}, errors.New("User not found"), http.StatusNotFound
	}
	return result, nil, 200
}

// update user by id
func (ps *UserRepository) UpdateProfile(ctx context.Context, id primitive.ObjectID, user Domain.User, current_user Domain.AccessClaims) (Dtos.OmitedUser, error, int) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	if current_user.ID != id {
		return Dtos.OmitedUser{}, errors.New("permission denied"), http.StatusForbidden
	}
	var NewUser Dtos.OmitedUser
	statusCode := 200

	// Retrieve the existing user
	NewUser, err, statusCode := ps.GetProfile(ctx, id, current_user)
	if err != nil {
		return Dtos.OmitedUser{}, err, 500
	}

	// Update only the specified fields
	if user.UserName != "" {
		NewUser.UserName = user.UserName
	}
	if user.Password != "" {
		err = password_services.CheckPasswordStrength(user.Password)
		if err != nil {
			return Dtos.OmitedUser{}, err, 400
		}
		newpass, er := password_services.GenerateFromPasswordCustom(user.Password)
		if er != nil {
			return Dtos.OmitedUser{}, er, 500
		}
		NewUser.Password = newpass
	}
	if user.Role != "" {
		NewUser.Role = user.Role
	}
	if user.ProfilePicture != "" {
		NewUser.ProfilePicture = user.ProfilePicture
	}
	if user.Bio != "" {
		NewUser.Bio = user.Bio
	}
	if user.Name != "" {
		NewUser.Name = user.Name
	}
	NewUser.UpdatedAt = time.Now()

	filter := bson.D{{"_id", id}}
	update := bson.D{
		{"$set", bson.D{
			{"email", NewUser.Email},
			{"username", NewUser.UserName},
			{"name", NewUser.Name},
			{"role", NewUser.Role},
			{"profile_picture", NewUser.ProfilePicture},
			{"bio", NewUser.Bio},
			{"updatedat", NewUser.UpdatedAt},
		}},
	}

	updateResult, err := ps.UserCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		statusCode = 500
		return Dtos.OmitedUser{}, err, statusCode
	}
	if updateResult.ModifiedCount == 0 {
		statusCode = 400
		return Dtos.OmitedUser{}, errors.New("user does not exist or no changes"), statusCode
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	return NewUser, nil, statusCode
}

// DeleteProfile removes a user profile, updates posts, and deletes the user's refresh token
func (ps *UserRepository) DeleteProfile(ctx context.Context, id primitive.ObjectID, current_user Domain.AccessClaims) (error, int) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	// Check if the current user has permission to delete the profile
	if current_user.ID != id {
		return errors.New("permission denied"), http.StatusForbidden
	}

	// Attempt to delete the user profile
	filter := bson.D{{"_id", id}}
	deleteResult, err := ps.UserCollection.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Println("Error deleting user profile:", err)
		return errors.New("internal server error"), http.StatusInternalServerError
	}
	if deleteResult.DeletedCount == 0 {
		return errors.New("user does not exist"), http.StatusNotFound
	}

	// Delete the refresh token associated with the user
	err, statusCode := ps.DeleteToken(ctx, id)
	if err != nil {
		return err, statusCode
	}

	return nil, http.StatusOK
}

func (ur *UserRepository) StoreToken(ctx context.Context, userid primitive.ObjectID, refreshToken string) (error, int) {
	token := Domain.RefreshToken{
		UserID:        userid,
		Refresh_token: refreshToken,
	}
	_, err := ur.refresh_collection.InsertOne(ctx, token)
	if err != nil {
		fmt.Println(err)
		return err, 500
	}
	return nil, 200
}

func (ur *UserRepository) UpdateToken(ctx context.Context, refreshToken string, userid primitive.ObjectID) (error, int) {
	//upaate the refresh token
	filter := primitive.D{{"_id", userid}}
	update := primitive.D{{"$set", primitive.D{{"refresh_token", refreshToken}}}}
	_, err := ur.refresh_collection.UpdateOne(ctx, filter, update)

	if err != nil {
		fmt.Println(err)
		return err, 500
	}

	return nil, 200
}

func (ur *UserRepository) DeleteToken(ctx context.Context, userid primitive.ObjectID) (error, int) {
	filter := primitive.D{{"_id", userid}}
	_, err := ur.refresh_collection.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return err, 500
	}
	return nil, 200
}

func (ur *UserRepository) FindToken(ctx context.Context, userid primitive.ObjectID) (string, error, int) {
	filter := primitive.D{{"_id", userid}}
	token := Domain.RefreshToken{}
	err := ur.refresh_collection.FindOne(ctx, filter).Decode(&token)
	if err != nil && err.Error() != "mongo: no documents in result" {
		fmt.Println(err)
		return "", err, 500
	}
	return token.Refresh_token, nil, 200
}
