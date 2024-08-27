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

func (lr *LoanRepository) GetLoans(ctx context.Context, filter Domain.Filter) ([]*Dtos.GetLoan, error, int, Domain.PaginationMetaData) {
	var loans []*Dtos.GetLoan
	// Initialize the filter for MongoDB query
	pipeline := []bson.M{}
	// Build the match stage for filtering
	matchStage := bson.M{}
	countfilter := bson.M{}

	// Set up pagination parameters
	page := 1
	if filter.Page > 1 {
		page = filter.Page
	}
	limit := 20
	if filter.Limit > 0 {
		limit = filter.Limit
	}

	// Count the number of documents that match the filter criteria

	// Default sort by updated_at in descending order
	orderBy := -1
	if filter.Status != "" {
		matchStage["status"] = filter.Status
		countfilter["status"] = filter.Status
		if filter.Status == "pending" {
			orderBy = 1
		} else {
			orderBy = -1
		}

	}
	if filter.OrderBy == 1 {
		orderBy = 1
	}
	// Add match stage to the pipeline if there are any filters
	if len(matchStage) > 0 {
		pipeline = append(pipeline, bson.M{"$match": matchStage})
	}

	count, err := lr.LoanCollection.CountDocuments(ctx, countfilter)
	if err != nil {
		return nil, err, 500, Domain.PaginationMetaData{}
	}
	sortBy := "updated_at"
	sort := bson.M{sortBy: orderBy}
	if filter.SortBy != "" {
		pipeline = append(pipeline, bson.M{"$sort": bson.M{sortBy: orderBy}})
	} else {
		pipeline = append(pipeline, bson.M{"$sort": sort})
	}

	pipeline = append(pipeline, bson.M{"$skip": int64((page - 1) * limit)})
	pipeline = append(pipeline, bson.M{"$limit": int64(limit)})

	cursor, err := lr.LoanCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err, 500, Domain.PaginationMetaData{}
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &loans); err != nil {
		return nil, err, 500, Domain.PaginationMetaData{}
	}
	// Return the list of posts, nil error, and a 200 status code
	paginationMetaData := Domain.PaginationMetaData{
		TotalRecords: int(count),
		TotalPages:   int(count / int64(limit)),
		PageSize:     limit,
		CurrentPage:  page,
	}

	return loans, nil, 200, paginationMetaData
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
	err := lr.validator.Struct(loan)
	if err != nil {
		return nil, errors.New("Invalid status"), 400
	}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: loan.Status}, {Key: "updatedAt", Value: loan.UpdatedAt}}}}
	_, err = lr.LoanCollection.UpdateOne(ctx, bson.D{{"_id", id}}, update)
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
