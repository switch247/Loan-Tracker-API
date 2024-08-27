package routers

import (
	"Loan-Tracker-API/Delivery/controllers"
	"Loan-Tracker-API/Infrastructure/auth_middleware"
	"Loan-Tracker-API/Repositories"

	emailservice "Loan-Tracker-API/Infrastructure/email_service"
	usecases "Loan-Tracker-API/UseCases"

	"github.com/go-playground/validator"
)

func LoanRouter() {
	validator := validator.New()
	email_service_reference := emailservice.NewMailService()

	loan_repo := Repositories.NewLoanRepository(LoanCollections.Loans, validator, email_service_reference)

	loan_usecase := usecases.NewLoanUseCases(loan_repo)
	loan_controller := controllers.NewLoanController(loan_usecase)

	is_authenticated := auth_middleware.AuthMiddleware()
	user_repo := Repositories.NewUserRepository(LoanCollections.Users, LoanCollections.RefreshTokens, validator, email_service_reference)
	is_admin := auth_middleware.IsAdminMiddleware(user_repo)

	adminRouter := Router.Group("/admin/loans")
	adminRouter.Use(is_authenticated, is_admin)
	{
		// admin routes
		// accept/reject loan
		adminRouter.PATCH("/:id/status", loan_controller.UpdateLoan)
		adminRouter.GET("/", loan_controller.GetLoans)
		adminRouter.DELETE("/:id", loan_controller.DeleteLoan)
		// extra
		adminRouter.GET("/:id", loan_controller.GetLoan)
	}

	userRouter := Router.Group("/loans")
	userRouter.Use(is_authenticated)
	{

		// apply for loan
		userRouter.POST("", loan_controller.CreateLoan)
		//get loan
		userRouter.GET("/:id", loan_controller.GetLoan)

	}

}
