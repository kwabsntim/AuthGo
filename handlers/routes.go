package handlers

import (
	"AuthGo/services"
	"net/http"
)

// amazonq-ignore-next-line

type ServiceContainer struct {
	RegisterService services.RegisterInterface
	LoginService    services.LoginInterface
}

func RouteSetup(services *ServiceContainer) *http.ServeMux {
	//using server mux to map the requests
	signupHandler := NewSignUpHandler(services.RegisterService)
	loginHandler := NewLoginHandler(services.LoginService)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/Signup", signupHandler.SignUp)
	mux.HandleFunc("/api/Login", loginHandler.Login)

	return mux
}
