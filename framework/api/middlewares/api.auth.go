package middlewares

import (
	"net/http"

	"github.com/go-chi/jwtauth"
)

func DynamicVerifier(authFunc func() *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := authFunc()
			verifier := jwtauth.Verifier(auth)
			authenticator := jwtauth.Authenticator

			verifier(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				authenticator(next).ServeHTTP(w, r)
			})).ServeHTTP(w, r)
		})
	}
}
