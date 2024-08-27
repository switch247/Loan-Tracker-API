package controllers

import (
	"Loan-Tracker-API/Domain"
	"Loan-Tracker-API/Dtos"
	"Loan-Tracker-API/Utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LoanController struct {
	LoanUsecase Domain.LoanUseCases
}

func NewLoanController(service_reference Domain.LoanUseCases) *LoanController {
	return &LoanController{
		LoanUsecase: service_reference,
	}
}

func (controller *LoanController) GetLoans(c *gin.Context) {
	queryparams := c.Request.URL.Query()
	filter := Domain.Filter{}

	// fill in filter values from the request query
	if len(queryparams) > 0 {
		filter.Status = queryparams.Get("status")
		filter.LoanerName = queryparams.Get("loanerName")
		filter.Limit, _ = strconv.Atoi(queryparams.Get("limit"))
		filter.Page, _ = strconv.Atoi(queryparams.Get("page"))
		filter.OrderBy, _ = strconv.Atoi(queryparams.Get("orderBy"))
		filter.SortBy = queryparams.Get("sortBy")
	}

	loans, err, status, paginationMetaData := controller.LoanUsecase.GetLoans(c, filter)
	if err != nil {
		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(status, gin.H{"loans": loans, "pagination": paginationMetaData})
}

func (controller *LoanController) GetLoan(c *gin.Context) {
	user, err := Getclaim(c)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	id, err := Utils.StringToObjectId(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	loan, err, status := controller.LoanUsecase.GetLoansById(c, id, *user)
	if err != nil {
		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(status, loan)
}

func (controller *LoanController) CreateLoan(c *gin.Context) {
	user, err := Getclaim(c)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	loan := &Domain.Loan{}
	err = c.BindJSON(loan)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	loan.UserID = user.ID
	if loan.Status == "" {
		loan.Status = "pending"

	}
	newLoan, err, status := controller.LoanUsecase.CreateLoan(c, loan)
	if err != nil {
		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(status, newLoan)
}

func (controller *LoanController) UpdateLoan(c *gin.Context) {
	id, err := Utils.StringToObjectId(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	loan := &Dtos.UpdateLoan{}
	err = c.BindJSON(loan)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	newLoan, err, status := controller.LoanUsecase.UpdateLoansById(c, id, loan)
	if err != nil {
		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(status, newLoan)

}

func (controller *LoanController) DeleteLoan(c *gin.Context) {
	id, err := Utils.StringToObjectId(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	err, status := controller.LoanUsecase.DeleteLoansById(c, id)
	if err != nil {
		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(status, gin.H{
		"message": "Loan deleted successfully",
	})
}
