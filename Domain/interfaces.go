package Domain

import (
	"Loan-Tracker-API/Dtos"
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) (Dtos.OmitedUser, error, int)
	GetUsers(ctx context.Context) ([]*Dtos.OmitedUser, error, int)
	GetUsersById(ctx context.Context, id primitive.ObjectID, user AccessClaims) (Dtos.OmitedUser, error, int)
	DeleteUsersById(ctx context.Context, id primitive.ObjectID, current_user AccessClaims) (error, int)

	ChangePassByEmail(ctx context.Context, email string, password string) (Dtos.OmitedUser, error, int)
	FindByEmail(ctx context.Context, email string) (Dtos.OmitedUser, error, int)

	Login(ctx context.Context, user *Dtos.LoginUserDto) (Tokens, error, int)
	Register(ctx context.Context, user *Dtos.RegisterUserDto) (*Dtos.OmitedUser, error, int)
	Logout(ctx context.Context, user_id primitive.ObjectID) (error, int)

	GenerateTokenFromUser(ctx context.Context, existingUser User) (Tokens, error, int)
	ResetPassword(ctx context.Context, email string, password string, resetToken string) (error, int)
	ForgetPassword(ctx context.Context, email string) (error, int)
	ActivateAccount(ctx context.Context, token string) (error, int)
	SendActivationEmail(email string) (error, int)

	GetProfile(ctx context.Context, id primitive.ObjectID, user AccessClaims) (Dtos.OmitedUser, error, int)
	UpdateProfile(ctx context.Context, id primitive.ObjectID, user User, current_user AccessClaims) (Dtos.OmitedUser, error, int)
	DeleteProfile(ctx context.Context, id primitive.ObjectID, current_user AccessClaims) (error, int)

	UpdateToken(ctx context.Context, refreshToken string, userid primitive.ObjectID) (error, int)
	DeleteToken(ctx context.Context, userid primitive.ObjectID) (error, int)
	FindToken(ctx context.Context, userid primitive.ObjectID) (string, error, int)
	StoreToken(ctx context.Context, userid primitive.ObjectID, refreshToken string) (error, int)
}

type UserUseCases interface {
	Login(c *gin.Context, user *Dtos.LoginUserDto) (Tokens, error, int)
	Register(c *gin.Context, user *Dtos.RegisterUserDto) (*Dtos.OmitedUser, error, int)
	Logout(c *gin.Context, user_id primitive.ObjectID) (error, int)

	ResetPassword(c *gin.Context, email string, password string, resetToken string) (error, int)
	ForgetPassword(c *gin.Context, email string) (error, int)
	ActivateAccount(c *gin.Context, token string) (error, int)

	CreateUser(c *gin.Context, user *User) (Dtos.OmitedUser, error, int)
	GetUsers(c *gin.Context) ([]*Dtos.OmitedUser, error, int)
	GetUsersById(c *gin.Context, id primitive.ObjectID, current_user AccessClaims) (Dtos.OmitedUser, error, int)
	DeleteUsersById(c *gin.Context, id primitive.ObjectID, current_user AccessClaims) (error, int)

	GetProfile(c *gin.Context, id primitive.ObjectID, current_user AccessClaims) (Dtos.OmitedUser, error, int)
	UpdateProfile(c *gin.Context, id primitive.ObjectID, user User, current_user AccessClaims) (Dtos.OmitedUser, error, int)
	DeleteProfile(c *gin.Context, id primitive.ObjectID, current_user AccessClaims) (error, int)
	// UpdateToken(c *gin.Context, refreshToken string, userid primitive.ObjectID) (error, int)
	DeleteToken(c *gin.Context, userid primitive.ObjectID) (error, int)
	FindToken(c *gin.Context, userid primitive.ObjectID) (string, error, int)
	StoreToken(c *gin.Context, userid primitive.ObjectID, refreshToken string) (error, int)
}

type LoanRepository interface {
	CreateLoan(ctx context.Context, loan *Loan) (*Dtos.GetLoan, error, int)
	GetLoans(ctx context.Context) ([]*Dtos.GetLoan, error, int)
	GetLoansById(ctx context.Context, id primitive.ObjectID) (*Dtos.GetLoan, error, int)
	UpdateLoansById(ctx context.Context, id primitive.ObjectID, loan *Dtos.UpdateLoan) (*Dtos.GetLoan, error, int)
	DeleteLoansById(ctx context.Context, id primitive.ObjectID) (error, int)
}

type LoanUseCases interface {
	CreateLoan(c *gin.Context, loan *Loan) (*Dtos.GetLoan, error, int)
	GetLoans(c *gin.Context) ([]*Dtos.GetLoan, error, int)
	GetLoansById(c *gin.Context, id primitive.ObjectID, user AccessClaims) (*Dtos.GetLoan, error, int)
	UpdateLoansById(c *gin.Context, id primitive.ObjectID, loan *Dtos.UpdateLoan) (*Dtos.GetLoan, error, int)
	DeleteLoansById(c *gin.Context, id primitive.ObjectID) (error, int)
}

type MailService interface {
	SendEmail(toEmail string, subject string, text string, category string) error
	SendActivationEmail(email string, activationToken string) error
	SendPasswordResetEmail(email string, resetToken string) error
}
