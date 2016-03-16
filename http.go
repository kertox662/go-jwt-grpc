package bearerware

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

/*
JWTFromHeader extracts a valid JWT from an http.Request or returns and error
*/
func JWTFromHeader(
	r *http.Request,
	keyFunc jwt.Keyfunc,
	signingMethod jwt.SigningMethod,
) (*jwt.Token, error) {
	authSlice, ok := r.Header[http.CanonicalHeaderKey(authHeader)]
	if !ok {
		return nil, errRestricted
	}
	for _, authString := range authSlice {
		tokenString, ok := tokenFromBearer(authString)
		if !ok {
			continue
		}
		return validJWTFromString(tokenString, keyFunc, signingMethod)
	}
	return nil, &jwtError{
		err:  errBearerFormat,
		code: invalidRequest,
	}
}

/*
WriteAuthError is a convenience function for setting the WWW-Authenticate header
and sending an http.Error()
*/
func WriteAuthError(w http.ResponseWriter, err error) {
	w.Header().Set(http.CanonicalHeaderKey("WWW-Authenticate"), err.Error())
	http.Error(
		w,
		fmt.Sprintf("%s: %s", http.StatusText(http.StatusUnauthorized), err),
		http.StatusUnauthorized,
	)
}
