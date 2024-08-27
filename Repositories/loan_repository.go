package Repositories

import (
	"Loan-Tracker-API/Domain"
	"Loan-Tracker-API/Dtos"
	emailservice "Loan-Tracker-API/Infrastructure/email_service"
	"Loan-Tracker-API/Utils"
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoanRepository struct {
	validator          *validator.Validate
	LoanCollection     Domain.Collection
	mu                 sync.RWMutex
	refresh_collection Domain.Collection
	emailservice       *emailservice.MailService
}

func NewLoanRepository(loan_collection Domain.Collection, _valdator *validator.Validate, email_service_reference *emailservice.MailService) *LoanRepository {

	return &LoanRepository{
		validator:      _valdator,
		LoanCollection: loan_collection,
		mu:             sync.RWMutex{},
		emailservice:   email_service_reference,
	}
}

func (lr *LoanRepository) CreateLoan(ctx context.Context, loan *Domain.Loan) (*Dtos.GetLoan, error, int) {
	lr.mu.RLock()
	defer lr.mu.RUnlock()

	loan.CreatedAt = primitive.NewDateTimeFromTime(Utils.GetCurrentTime())
	loan.UpdatedAt = primitive.NewDateTimeFromTime(Utils.GetCurrentTime())

	insertResult, err := lr.LoanCollection.InsertOne(ctx, loan)
	if err != nil {
		return &Dtos.GetLoan{}, err, 500
	}
	load_id := insertResult.InsertedID.(primitive.ObjectID)
	fmt.Println("Inserted a single document: ", load_id)
	newLoan := &Dtos.GetLoan{}
	err = lr.LoanCollection.FindOne(context.TODO(), bson.D{{"_id", load_id}}).Decode(&newLoan)
	if err != nil {
		fmt.Println(err)
		return &Dtos.GetLoan{}, errors.New("User Not Created"), 500
	}
	return newLoan, nil, 201
}

func (lr *LoanRepository) GetLoans(ctx context.Context) ([]*Dtos.GetLoan, error, int) {
	cursor, err := lr.LoanCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err, 500
	}
	var loans []*Dtos.GetLoan
	for cursor.Next(ctx) {
		var loan *Dtos.GetLoan
		err := cursor.Decode(&loan)
		if err != nil {
			return nil, err, 500
		}
		loans = append(loans, loan)
	}
	return loans, nil, 200
}

func (lr *LoanRepository) GetLoansById(ctx context.Context, id primitive.ObjectID) (*Dtos.GetLoan, error, int) {

	loan := &Dtos.GetLoan{}
	err := lr.LoanCollection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&loan)
	if err != nil {
		return nil, err, 500
	}
	return loan, nil, 200
}

func (lr *LoanRepository) UpdateLoansById(ctx context.Context, id primitive.ObjectID, loan *Dtos.UpdateLoan) (*Dtos.GetLoan, error, int) {
	lr.mu.RLock()
	defer lr.mu.RUnlock()

	loan.UpdatedAt = primitive.NewDateTimeFromTime(Utils.GetCurrentTime())
	_, err := lr.LoanCollection.UpdateOne(ctx, bson.D{{"_id", id}}, bson.D{{"$set", loan}})
	if err != nil {
		return nil, err, 500
	}
	newLoan := &Dtos.GetLoan{}
	err = lr.LoanCollection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&newLoan)
	if err != nil {
		return nil, err, 500
	}
	return newLoan, nil, 200
}

func (lr *LoanRepository) DeleteLoansById(ctx context.Context, id primitive.ObjectID) (error, int) {
	lr.mu.RLock()
	defer lr.mu.RUnlock()

	_, err := lr.LoanCollection.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return err, 500
	}
	return nil, 200
}
