package controllers

import (
	"Loan-Tracker-API/Domain"
	"Loan-Tracker-API/Dtos"
	"Loan-Tracker-API/Utils"

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
	loans, err, status := controller.LoanUsecase.GetLoans(c)
	if err != nil {
		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(status, loans)
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
	loan := &Domain.Loan{}
	err := c.BindJSON(loan)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
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
