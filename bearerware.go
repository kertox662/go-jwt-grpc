package bearerware

import (
	"errors"
	"fmt"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	authHeader = "authorization"
	bearer     = "bearer "
	bearerLen  = len(bearer)

	invalidRequest = "invalid_request"
	invalidToken   = "invalid_token"
	//insufficientScope = "insufficient_scope"

	errSigningMethod = `Expected signing method %s but token is signed using %s`
)

var (
	errRestricted   = errors.New("Bearer realm=Restricted")
	errBearerFormat = errors.New(
		"Authorization header format must be Bearer {token}",
	)
	errTokenInvalid = errors.New("Token is invalid")
)

type jwtError struct {
	err  error
	code string
}

func (j *jwtError) Error() string {
	s := errRestricted.Error()
	if len(j.code) > 0 {
		s += fmt.Sprintf(`,error="%s",error_description="%s"`, j.code, j.err)
	}
	return s
}

func isBearerToken(s string) bool {
	return len(s) > bearerLen && strings.ToLower(s)[:bearerLen] == bearer
}

func tokenFromBearer(s string) (string, bool) {
	if isBearerToken(s) {
		return s[bearerLen:], true
	}
	return "", false
}

//CheckJWT checks if the JWT was issued by this app
func validJWTFromString(
	token string,
	keyFunc jwt.Keyfunc,
	signingMethod jwt.SigningMethod,
) (*jwt.Token, error) {
	var (
		err         error
		parsedToken *jwt.Token
	)
	if len(token) == 0 {
		return nil, &jwtError{
			err:  errBearerFormat,
			code: invalidRequest,
		}
	}
	parsedToken, err = jwt.Parse(token, keyFunc)
	if err != nil {
		return nil, &jwtError{
			err:  fmt.Errorf("Error parsing token: %v", err),
			code: invalidToken,
		}
	}
	if alg := parsedToken.Header["alg"]; alg != signingMethod.Alg() {
		return nil, &jwtError{
			err:  fmt.Errorf(errSigningMethod, signingMethod.Alg(), alg),
			code: invalidToken,
		}
	}
	if !parsedToken.Valid {
		return nil, &jwtError{
			err:  errTokenInvalid,
			code: invalidToken,
		}
	}
	return parsedToken, nil
}
