package bearerware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
)

func TestJWTFromHeader(t *testing.T) {
	var (
		jwtKeyFunc = func(token *jwt.Token) (interface{}, error) {
			return []byte(" "), nil
		}

		tests = []struct {
			req           *http.Request
			signingMethod jwt.SigningMethod
			err           error
		}{
			{
				req: &http.Request{
					Header: http.Header{
						"Authorization": {
							"Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiIiLCJpYXQiOm51bGwsImV4cCI6bnVsbCwiYXVkIjoiIiwic3ViIjoiIn0.IlffGJz3IyFX1ADQ6-jOTQ_0D-K0kuKq5SpB_oirCrk",
						},
					},
				},
				signingMethod: jwt.SigningMethodHS256,
				err:           nil,
			},
			{
				req: &http.Request{
					Header: http.Header{
						"Authorization": {
							"Beare eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiIiLCJpYXQiOm51bGwsImV4cCI6bnVsbCwiYXVkIjoiIiwic3ViIjoiIn0.IlffGJz3IyFX1ADQ6-jOTQ_0D-K0kuKq5SpB_oirCrk",
						},
					},
				},
				signingMethod: jwt.SigningMethodHS256,
				err:           errors.New(`Bearer realm=Restricted,error="invalid_request",error_description="Authorization header format must be Bearer {token}"`),
			},
			{
				req:           &http.Request{},
				signingMethod: jwt.SigningMethodHS256,
				err:           errors.New(`Bearer realm=Restricted`),
			},
		}
	)
	for _, test := range tests {
		token, err := JWTFromHeader(test.req, jwtKeyFunc, test.signingMethod)
		if e := testJWTFrom(token, err, test.err); e != nil {
			t.Error(e)
		}
	}
}

func TestWriteHeader(t *testing.T) {
	w := httptest.NewRecorder()
	WriteAuthError(w, errors.New("foo"))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
	authStatus := w.Header().Get(http.CanonicalHeaderKey("WWW-Authenticate"))
	if authStatus != "foo" {
		t.Fatalf("Expected 'foo', Got %s", authStatus)
	}
}
