package handlers

import (
	"net/http"
)

func RouteSetup() *http.ServeMux {
	//using server mux to map the requests

	mux := http.NewServeMux()
	mux.HandleFunc("/api/Signup", SignUpHandler)

	return mux

}
