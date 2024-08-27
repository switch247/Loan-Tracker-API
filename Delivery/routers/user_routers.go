package routers

import (
	"Loan-Tracker-API/Delivery/controllers"
	"Loan-Tracker-API/Infrastructure/auth_middleware"
	emailservice "Loan-Tracker-API/Infrastructure/email_service"
	"Loan-Tracker-API/Repositories"
	usecases "Loan-Tracker-API/UseCases"

	"github.com/go-playground/validator"
)

func UserRouter() {
	validator := validator.New()
	email_service_reference := emailservice.NewMailService()

	user_repo := Repositories.NewUserRepository(BlogCollections.Users, BlogCollections.RefreshTokens, validator, email_service_reference)

	user_usecase := usecases.NewUserUseCase(user_repo)
	user_controller := controllers.NewUserController(user_usecase)

	userRouter := Router.Group("/users")
	is_authenticated := auth_middleware.AuthMiddleware()
	is_admin := auth_middleware.IsAdminMiddleware(user_repo)
	{
		// users routes
		userRouter.GET("/", is_authenticated, is_admin, user_controller.GetUsers)
		userRouter.DELETE("/:id", is_authenticated, is_admin, user_controller.DeleteUser)
		// extra
		userRouter.POST("/", is_authenticated, is_admin, user_controller.CreateUser)
		userRouter.GET("/:id", is_authenticated, is_admin, user_controller.GetUser)
	}

	authRouter := userRouter.Group("")
	{

		// register
		authRouter.POST("/register", user_controller.Register)
		//login
		authRouter.POST("/login", user_controller.Login)

		//logout
		authRouter.GET("/logout", auth_middleware.AuthMiddleware(), user_controller.Logout)
		// forget password
		authRouter.POST("/password-reset", user_controller.ForgetPassword)
		authRouter.GET("/password-reset/:reset_token", user_controller.ForgetPasswordForm)
		authRouter.POST("/password-update/:reset_token", user_controller.ResetPassword)

		// activate account
		authRouter.GET("/verify-email/:activation_token", user_controller.ActivateAccount)

	}

	profileRouter := userRouter.Group("/profile")
	profileRouter.Use(auth_middleware.AuthMiddleware())
	{

		// get all users
		profileRouter.GET("/", user_controller.GetProfile)
		profileRouter.PATCH("/", user_controller.UpdateProfile)
		profileRouter.DELETE("/", user_controller.DeleteProfile)

	}

	refreshRouter := userRouter.Group("/token")
	refreshRouter.Use(auth_middleware.AuthMiddleware())
	{
		refreshRouter.GET("refresh", user_controller.Refresh)
	}
}
