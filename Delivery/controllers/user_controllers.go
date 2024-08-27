package controllers

import (
	templates "Loan-Tracker-API/Delivery/templates"
	"Loan-Tracker-API/Domain"
	"Loan-Tracker-API/Utils"
	"fmt"
	"net/http"

	"Loan-Tracker-API/Dtos"
	jwtservice "Loan-Tracker-API/Infrastructure/jwt_service"
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	UserUseCase Domain.UserUseCases
}

func NewUserController(service_reference Domain.UserUseCases) *UserController {
	return &UserController{
		UserUseCase: service_reference,
	}
}

func (pc *UserController) GetUsers(c *gin.Context) {

	users, err, statusCode := pc.UserUseCase.GetUsers(c)
	if err != nil {
		c.IndentedJSON(statusCode, gin.H{"error": err.Error()})
	} else {
		c.IndentedJSON(statusCode, gin.H{"users": users})
	}
}

func (pc *UserController) GetUser(c *gin.Context) {
	cur_user, err := Getclaim(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.IndentedJSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	user, err, statusCode := pc.UserUseCase.GetUsersById(c, objectID, *cur_user)
	if err != nil {
		c.IndentedJSON(statusCode, gin.H{"error": err.Error()})
	} else {
		c.IndentedJSON(statusCode, gin.H{"user": user})
	}
}

func (pc *UserController) CreateUser(c *gin.Context) {
	var user Domain.User
	if err := c.BindJSON(&user); err != nil {
		c.IndentedJSON(400, gin.H{"error": err.Error()})
		return
	}
	v := validator.New()
	if err := v.Struct(user); err != nil {
		fmt.Printf(err.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid or missing data", "error": err.Error()})
		return
	}

	OmitedUser, err, statusCode := pc.UserUseCase.CreateUser(c, &user)
	if err != nil {
		c.IndentedJSON(statusCode, gin.H{"error": err.Error()})
	} else {
		c.IndentedJSON(statusCode, gin.H{"user": OmitedUser})
	}
}

func (pc *UserController) DeleteUser(c *gin.Context) {
	user, err := Getclaim(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}

	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.IndentedJSON(400, gin.H{"error": "Invalid ID"})
		return
	}
	err, statusCode := pc.UserUseCase.DeleteUsersById(c, objectID, *user)
	if err != nil {
		c.IndentedJSON(statusCode, gin.H{"error": err.Error()})
	} else {
		c.IndentedJSON(statusCode, gin.H{"message": "User deleted successfully"})
	}
}

// authentication

// login
func (ac *UserController) Login(c *gin.Context) {
	var newUser Dtos.LoginUserDto
	v := validator.New()
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid data", "error": err.Error()})
		return
	}
	if err := v.Struct(newUser); err != nil {
		if newUser.Email == "" && newUser.UserName == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid or missing data", "error": "email or username is required"})
			return
		}
		// fmt.Println(err.Error())
		// c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid or missing data", "error": err.Error()})
		// return
	}
	fmt.Println(newUser)
	token, err, statusCode := ac.UserUseCase.Login(c, &newUser)
	if err != nil {
		c.IndentedJSON(statusCode, gin.H{"error": err.Error()})
	} else {
		//success
		c.IndentedJSON(http.StatusOK, gin.H{"message": "User logged in successfully",
			"acess_token": token.AccessToken})
	}

}

// register
func (ac *UserController) Register(c *gin.Context) {
	// return error
	var newUser Dtos.RegisterUserDto
	v := validator.New()
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid data", "error": err.Error()})
		return
	}

	if err := v.Struct(newUser); err != nil {
		fmt.Printf(err.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid or missing data", "error": err.Error()})
		return
	}

	createdUser, err, statusCode := ac.UserUseCase.Register(c, &newUser)

	if err != nil {
		c.IndentedJSON(statusCode, gin.H{"error": err.Error()})
	} else {
		//success
		c.IndentedJSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": createdUser})
	}

}

// logout
func (ac *UserController) Logout(c *gin.Context) {
	// return error
	// get the access token from the header
	claims, err := Getclaim(c)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	err, statusCode := ac.UserUseCase.Logout(c, claims.ID)
	if err != nil {
		c.IndentedJSON(statusCode, gin.H{"error": err.Error()})
	} else {
		//success
		c.IndentedJSON(http.StatusOK, gin.H{"message": "User logged out successfully"})
	}

}

// sends email with token and reset link
func (ac *UserController) ForgetPassword(c *gin.Context) {
	email := c.PostForm("email")
	err, statusCode := ac.UserUseCase.ForgetPassword(c, email)
	if err != nil {
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(statusCode, gin.H{"message": "reset token sent successfully"})
	}
}

// ForgetPasswordForm handles the rendering of the reset password form
func (ac *UserController) ForgetPasswordForm(c *gin.Context) {
	resetToken := c.Params.ByName("reset_token")
	_, err := jwtservice.VerifyToken(resetToken)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid token"})
		return
	}
	t, err := template.New("reset").Parse(templates.ResetTemplate)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error parsing template"})
		return
	}

	err = t.Execute(c.Writer, gin.H{"ResetToken": resetToken})

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error executing template"})
		return
	}
}

// reset password
func (ac *UserController) ResetPassword(c *gin.Context) {
	// extracts token and new_password from the request if correct update the password
	resetToken := c.Params.ByName("reset_token")
	email, err := jwtservice.VerifyToken(resetToken)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid token"})
		return
	}
	password := c.PostForm("password")
	if password == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "password is required"})
		return
	}

	err, statusCode := ac.UserUseCase.ResetPassword(c, email, password, resetToken)
	if err != nil {
		c.IndentedJSON(statusCode, gin.H{"message": err.Error()})
		return
	} else {
		c.IndentedJSON(statusCode, gin.H{"message": "password reset successfully"})
	}

	fmt.Println("password:", password, "reset_token", resetToken)
}

func (ac *UserController) ActivateAccount(c *gin.Context) {
	activationToken := c.Params.ByName("activation_token")
	err, statusCode := ac.UserUseCase.ActivateAccount(c, activationToken)
	if err != nil {
		c.IndentedJSON(statusCode, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(statusCode, gin.H{"message": "account activated successfully"})
}

// profile

func (pc *UserController) GetProfile(c *gin.Context) {
	cur_user, err := Getclaim(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}

	user, err, statusCode := pc.UserUseCase.GetProfile(c, cur_user.ID, *cur_user)
	if err != nil {
		c.IndentedJSON(statusCode, gin.H{"error": err.Error()})
	} else {
		c.IndentedJSON(statusCode, gin.H{"Profile": user})
	}
}

func (pc *UserController) UpdateProfile(c *gin.Context) {

	logeduser, err := Getclaim(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}
	// updateUser := Domain.UpdateUser{}

	// if err := c.BindJSON(&updateUser); err != nil {
	// 	fmt.Println("i am at thr top")
	// 	c.IndentedJSON(400, gin.H{"error": err.Error()})
	// 	return
	// }
	// get profile picture image from request
	file, _ := c.FormFile("profilepicture")
	// if err != nil {
	// 	c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	var user Domain.User
	user.ID = logeduser.ID
	user.Name = c.PostForm("name")
	user.UserName = c.PostForm("username")
	user.Password = c.PostForm("password")
	user.Bio = c.PostForm("bio")
	if file != nil {
		profilePicture, err := Utils.SetProfilePicture(file)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user.ProfilePicture = profilePicture
	}

	OmitedUser, err, statusCode := pc.UserUseCase.UpdateProfile(c, logeduser.ID, user, *logeduser)
	if err != nil {
		fmt.Println("i was here all along ")
		c.IndentedJSON(statusCode, gin.H{"error": err.Error()})
	} else {
		c.IndentedJSON(statusCode, gin.H{"Profile": OmitedUser})
	}
}

func (pc *UserController) DeleteProfile(c *gin.Context) {
	user, err := Getclaim(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}

	err, statusCode := pc.UserUseCase.DeleteProfile(c, user.ID, *user)
	if err != nil {
		c.IndentedJSON(statusCode, gin.H{"error": err.Error()})
	} else {
		c.IndentedJSON(statusCode, gin.H{"message": "Profile deleted successfully"})
	}
}

// Refresh  controller

// Refresh function
func (r *UserController) Refresh(c *gin.Context) {
	accessClaims, err := Getclaim(c)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	// get the refresh token
	refreshToken, err, statuscode := r.UserUseCase.FindToken(c, accessClaims.ID)

	if err != nil {
		c.JSON(statuscode, gin.H{"error": err.Error()})
		return
	}

	// check if the refresh token is valid
	if refreshToken == "" {
		c.JSON(401, gin.H{"error": "refresh token not found"})
		return
	}

	// verify the refresh token
	err = jwtservice.VerifyRefreshToken(refreshToken, accessClaims.ID)

	fmt.Println("refresh token", refreshToken)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	user := Domain.User{
		ID:   accessClaims.ID,
		Role: accessClaims.Role,
	}

	newAccessToken, err := jwtservice.CreateAccessToken(user)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// newRefreshToken, err := jwtservice.CreateRefreshToken(user)

	// if err != nil {
	// 	c.JSON(500, gin.H{"error": err.Error()})
	// 	return
	// }

	// store the refresh token
	// err, statuscode = r.UserUseCase.StoreToken(c,  accessClaims.ID, newRefreshToken)

	// if err != nil {
	// 	c.JSON(statuscode, gin.H{"error": err.Error()})
	// 	return
	// }

	c.JSON(200, gin.H{"access_token": newAccessToken, "refresh_token": refreshToken})
}
