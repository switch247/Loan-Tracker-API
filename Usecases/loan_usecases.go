package UseCases

import (
	"Loan-Tracker-API/Domain"
	"Loan-Tracker-API/Dtos"
	"Loan-Tracker-API/Utils"
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoanUseCases struct {
	LoanRepository Domain.LoanRepository
}

func NewLoanUseCases(service_reference Domain.LoanRepository) *LoanUseCases {
	return &LoanUseCases{
		LoanRepository: service_reference,
	}
}

func (usecase *LoanUseCases) GetLoans(c *gin.Context, filter Domain.Filter) ([]*Dtos.GetLoan, error, int, Domain.PaginationMetaData) {
	loans, err, status, paginationMetaData := usecase.LoanRepository.GetLoans(context.Background(), filter)
	if err != nil {
		return nil, err, status, paginationMetaData
	}
	return loans, nil, 200, paginationMetaData
}

func (usecase *LoanUseCases) GetLoansById(c *gin.Context, id primitive.ObjectID, user Domain.AccessClaims) (*Dtos.GetLoan, error, int) {

	loan, err, status := usecase.LoanRepository.GetLoansById(context.Background(), id)
	if err != nil {
		return nil, err, status
	}
	if user.Role != "admin" && loan.UserID != user.ID {
		return nil, errors.New("Unauthorized Access"), 401
	}
	return loan, nil, 200
}

func (usecase *LoanUseCases) CreateLoan(c *gin.Context, loan *Domain.Loan) (*Dtos.GetLoan, error, int) {
	loan.ID = primitive.NewObjectID()
	loan.CreatedAt = primitive.NewDateTimeFromTime(Utils.GetCurrentTime())
	loan.UpdatedAt = primitive.NewDateTimeFromTime(Utils.GetCurrentTime())
	newLoan, err, status := usecase.LoanRepository.CreateLoan(context.Background(), loan)
	if err != nil {
		return nil, err, status
	}
	return newLoan, nil, 201
}

func (usecase *LoanUseCases) UpdateLoansById(c *gin.Context, id primitive.ObjectID, loan *Dtos.UpdateLoan) (*Dtos.GetLoan, error, int) {
	loan.UpdatedAt = primitive.NewDateTimeFromTime(Utils.GetCurrentTime())
	newLoan, err, status := usecase.LoanRepository.UpdateLoansById(context.Background(), id, loan)
	if err != nil {
		return nil, err, status
	}
	return newLoan, nil, 200
}

func (usecase *LoanUseCases) DeleteLoansById(c *gin.Context, id primitive.ObjectID) (error, int) {
	err, status := usecase.LoanRepository.DeleteLoansById(context.Background(), id)
	if err != nil {
		return err, status
	}
	return nil, 200
}
