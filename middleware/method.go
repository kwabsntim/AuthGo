package middleware

import (
	"net/http"
)

// method checking middleware
func MethodChecker(method string) func(http.HandlerFunc) http.HandlerFunc {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != method {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
			handler(w, r)
		}
	}
}
