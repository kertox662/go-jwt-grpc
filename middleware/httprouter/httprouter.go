package bwhttprouter

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/julienschmidt/httprouter"
	bearerware "github.com/kertox662/go-jwt-grpc"
)

//JWTContext stores the request scoped tokens
var JWTContext bearerware.JWTContexter

func init() {
	JWTContext = bearerware.NewJWTContext()
}

//JWTHandler is a JWT middleware for httprouter
func JWTHandler(
	h httprouter.Handle,
	keyFunc jwt.Keyfunc,
	signingMethod jwt.SigningMethod,
) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		j, err := bearerware.JWTFromHeader(r, keyFunc, signingMethod)
		if err != nil {
			bearerware.WriteAuthError(w, err)
			return
		}
		JWTContext.WriteJWT(r, j)
		defer JWTContext.DeleteJWT(r)

		h(w, r, ps)
		return
	}
}
