package bwhttprouter

import (
	"net/http"

	"github.com/ckaznocha/go-JWTBearerware"
	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
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
