package bearerware

import (
	"errors"
	"fmt"
	"testing"

	"github.com/dgrijalva/jwt-go"
)

func Test_validJWTFromString(t *testing.T) {
	var (
		jwtKeyFunc = func(token *jwt.Token) (interface{}, error) {
			return []byte(" "), nil
		}
		tests = []struct {
			jwt           string
			signingMethod jwt.SigningMethod
			isValidJWT    bool
			err           error
		}{
			{
				jwt:           "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiIiLCJpYXQiOm51bGwsImV4cCI6bnVsbCwiYXVkIjoiIiwic3ViIjoiIn0.IlffGJz3IyFX1ADQ6-jOTQ_0D-K0kuKq5SpB_oirCrk",
				signingMethod: jwt.SigningMethodHS256,
				isValidJWT:    true,
				err:           nil,
			},
			{
				jwt:           "",
				signingMethod: jwt.SigningMethodHS256,
				isValidJWT:    false,
				err:           errors.New(`Bearer realm=Restricted,error="invalid_request",error_description="Authorization header format must be Bearer {token}"`),
			},
			{
				jwt:           "foo",
				signingMethod: jwt.SigningMethodHS256,
				isValidJWT:    false,
				err:           errors.New(`Bearer realm=Restricted,error="invalid_token",error_description="Error parsing token: token contains an invalid number of segments"`),
			},
			{
				jwt:           "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiIiLCJpYXQiOm51bGwsImV4cCI6bnVsbCwiYXVkIjoiIiwic3ViIjoiIn0.IlffGJz3IyFX1ADQ6-jOTQ_0D-K0kuKq5SpB_oirCrk",
				signingMethod: jwt.SigningMethodES256,
				isValidJWT:    false,
				err:           errors.New(`Bearer realm=Restricted,error="invalid_token",error_description="Expected signing method ES256 but token is signed using HS256"`),
			},
			{
				jwt:           "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiIiLCJpYXQiOm51bGwsImV4cCI6bnVsbCwiYXVkIjoiIiwic3ViIjoiIn0.EK492gkeFLTSWLrQlu2hTNbFw3scKJqU5CA3sIsLf68",
				signingMethod: jwt.SigningMethodES256,
				isValidJWT:    false,
				err:           errors.New(`Bearer realm=Restricted,error="invalid_token",error_description="Error parsing token: signature is invalid"`),
			},
		}
	)
	for _, test := range tests {
		_, err := validJWTFromString(test.jwt, jwtKeyFunc, test.signingMethod)
		switch {
		case test.isValidJWT && err != nil:
			t.Error(err)
		case !test.isValidJWT && err.Error() != test.err.Error():
			t.Errorf("Expected Err: %s, Got: %s", test.err, err)

		}
	}
}

func Test_isBearerToken(t *testing.T) {
	var (
		token = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiIiLCJpYXQiOm51bGwsImV4cCI6bnVsbCwiYXVkIjoiIiwic3ViIjoiIn0.IlffGJz3IyFX1ADQ6-jOTQ_0D-K0kuKq5SpB_oirCrk"
		tests = []struct {
			s     string
			valid bool
		}{
			{s: fmt.Sprintf("Bearer %s", token), valid: true},
			{s: fmt.Sprintf("bearer %s", token), valid: true},
			{s: fmt.Sprintf("beaRer %s", token), valid: true},
			{s: fmt.Sprintf("Bearer%s", token), valid: false},
			{s: fmt.Sprintf("%s", token), valid: false},
			{s: fmt.Sprintf("MAC %s", token), valid: false},
		}
	)
	for _, test := range tests {
		got, valid := tokenFromBearer(test.s)
		if valid != test.valid {
			t.Errorf("Expected %t, got %t", test.valid, valid)
		}
		if test.valid && got != token {
			t.Errorf("Expected %s, got %s", token, got)
		}
	}
}
