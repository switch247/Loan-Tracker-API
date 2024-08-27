package UseCases

import (
	"Loan-Tracker-API/Domain"
	"Loan-Tracker-API/Dtos"
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserUseCases struct {
	UserRepository Domain.UserRepository
	contextTimeout time.Duration
}

func NewUserUseCase(service_reference Domain.UserRepository) *UserUseCases {
	return &UserUseCases{
		UserRepository: service_reference,
		contextTimeout: time.Second * 10,
	}
}

func (uc *UserUseCases) GetUsers(c *gin.Context) ([]*Dtos.OmitedUser, error, int) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.UserRepository.GetUsers(ctx)

}

func (uc *UserUseCases) GetUsersById(c *gin.Context, id primitive.ObjectID, current_user Domain.AccessClaims) (Dtos.OmitedUser, error, int) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.UserRepository.GetUsersById(ctx, id, current_user)

}

func (uc *UserUseCases) CreateUser(c *gin.Context, user *Domain.User) (Dtos.OmitedUser, error, int) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.UserRepository.CreateUser(ctx, user)

}

func (uc *UserUseCases) DeleteUsersById(c *gin.Context, id primitive.ObjectID, current_user Domain.AccessClaims) (error, int) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.UserRepository.DeleteUsersById(ctx, id, current_user)

}

// auth usecase

// login
func (a *UserUseCases) Login(c *gin.Context, user *Dtos.LoginUserDto) (Domain.Tokens, error, int) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	return a.UserRepository.Login(ctx, user)
}

// register
func (a *UserUseCases) Register(c *gin.Context, user *Dtos.RegisterUserDto) (*Dtos.OmitedUser, error, int) {
	// return error
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	return a.UserRepository.Register(ctx, user)
}

// logout

func (a *UserUseCases) Logout(c *gin.Context, user_id primitive.ObjectID) (error, int) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	return a.UserRepository.Logout(ctx, user_id)
}

func (a *UserUseCases) ForgetPassword(c *gin.Context, email string) (error, int) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	return a.UserRepository.ForgetPassword(ctx, email)
}

func (a *UserUseCases) ResetPassword(c *gin.Context, email string, password string, resetToken string) (error, int) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	return a.UserRepository.ResetPassword(ctx, email, password, resetToken)
}

func (a *UserUseCases) ActivateAccount(c *gin.Context, token string) (error, int) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	return a.UserRepository.ActivateAccount(ctx, token)
}

// profile usecase

func (uc *UserUseCases) GetProfile(c *gin.Context, id primitive.ObjectID, current_user Domain.AccessClaims) (Dtos.OmitedUser, error, int) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.UserRepository.GetProfile(ctx, id, current_user)

}

func (uc *UserUseCases) UpdateProfile(c *gin.Context, id primitive.ObjectID, user Domain.User, current_user Domain.AccessClaims) (Dtos.OmitedUser, error, int) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.UserRepository.UpdateProfile(ctx, id, user, current_user)

}

func (uc *UserUseCases) DeleteProfile(c *gin.Context, id primitive.ObjectID, current_user Domain.AccessClaims) (error, int) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.UserRepository.DeleteProfile(ctx, id, current_user)

}

// refresh usecase

// // Refresh function
// func (r *UserUseCases) UpdateToken(c *gin.Context,refreshToken string, userid primitive.ObjectID) (error, int) {
// 	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
// 	defer cancel()
// 	return r.UserRepository.UpdateToken(ctx, refreshToken, userid)
// }

// Delete function
func (r *UserUseCases) DeleteToken(c *gin.Context, userid primitive.ObjectID) (error, int) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()
	return r.UserRepository.DeleteToken(ctx, userid)
}

// Find function
func (r *UserUseCases) FindToken(c *gin.Context, userid primitive.ObjectID) (string, error, int) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()
	return r.UserRepository.FindToken(ctx, userid)
}

// store token function
func (r *UserUseCases) StoreToken(c *gin.Context, userid primitive.ObjectID, refreshToken string) (error, int) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()
	return r.UserRepository.StoreToken(ctx, userid, refreshToken)
}
