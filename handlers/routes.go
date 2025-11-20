package handlers

import (
	"AuthGo/middleware"
	"AuthGo/services"
	"net/http"
)

// amazonq-ignore-next-line

type ServiceContainer struct {
	RegisterService services.RegisterInterface
	LoginService    services.LoginInterface
	FetchService    services.FetchUsersInterface
}

func RouteSetup(services *ServiceContainer) *http.ServeMux {
	//using server mux to map the requests
	signupHandler := NewSignUpHandler(services.RegisterService)
	loginHandler := NewLoginHandler(services.LoginService)
	FetchUsersHandler := NewFetchHandler(services.FetchService) //the service must be in capital eg.FetchService

	mux := http.NewServeMux()
	mux.HandleFunc("/api/Signup", middleware.MethodChecker("POST")(signupHandler.SignUp))
	mux.HandleFunc("/api/Login", middleware.MethodChecker("POST")(loginHandler.Login))
	mux.HandleFunc("/api/FetchUsers", middleware.MethodChecker("GET")(FetchUsersHandler.FetchAllUsers))

	return mux
}
