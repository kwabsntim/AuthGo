package handlers

import (
	"AuthGo/services"
	"net/http"
)

func RouteSetup(userService services.UserService) *http.ServeMux {
	//using server mux to map the requests
	signupHandler := NewSignUpHandler(userService)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/Signup", signupHandler.Handle)
	mux.HandleFunc("/api/Login", LoginUser)

	return mux
}
