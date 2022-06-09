package bwhttp

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	bearerware "github.com/kertox662/go-jwt-grpc"
)

//JWTContext stores the request scoped tokens
var JWTContext bearerware.JWTContexter

func init() {
	JWTContext = bearerware.NewJWTContext()
}

//JWTHandler is JWT middleware for net/http
func JWTHandler(
	h http.HandlerFunc,
	keyFunc jwt.Keyfunc,
	signingMethod jwt.SigningMethod,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			h(w, r)
			return
		}
		j, err := bearerware.JWTFromHeader(r, keyFunc, signingMethod)
		if err != nil {
			bearerware.WriteAuthError(w, err)
			return
		}
		JWTContext.WriteJWT(r, j)
		defer JWTContext.DeleteJWT(r)
		h(w, r)
		return
	}
}
