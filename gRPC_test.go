package bearerware

import (
	"errors"
	"reflect"
	"testing"

	"google.golang.org/grpc/metadata"

	"github.com/dgrijalva/jwt-go"

	"golang.org/x/net/context"
)

func TestJWTAccessCredentials(t *testing.T) {
	var tests = []struct {
		token string
		md    map[string]string
	}{
		{
			token: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiIiLCJpYXQiOm51bGwsImV4cCI6bnVsbCwiYXVkIjoiIiwic3ViIjoiIn0.IlffGJz3IyFX1ADQ6-jOTQ_0D-K0kuKq5SpB_oirCrk",
			md:    map[string]string{authHeader: "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiIiLCJpYXQiOm51bGwsImV4cCI6bnVsbCwiYXVkIjoiIiwic3ViIjoiIn0.IlffGJz3IyFX1ADQ6-jOTQ_0D-K0kuKq5SpB_oirCrk"},
		},
	}
	for _, test := range tests {
		cred, err := NewJWTAccessFromJWT(test.token)
		if err != nil {
			t.Error(err)
		}
		md, err := cred.GetRequestMetadata(context.Background())
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(md, test.md) {
			t.Errorf("Expected: %+v, Got: %+v", test.md, md)
		}
		require := cred.RequireTransportSecurity()
		if !require {
			t.Error("Expected RequireTransportSecurity to be true but got false")
		}
	}
}

func TestJWTFromContext(t *testing.T) {
	var (
		jwtKeyFunc = func(token *jwt.Token) (interface{}, error) {
			return []byte(" "), nil
		}

		tests = []struct {
			ctx           context.Context
			signingMethod jwt.SigningMethod
			err           error
		}{
			{
				ctx: metadata.NewIncomingContext(
					context.Background(),
					metadata.New(map[string]string{authHeader: "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiIiLCJpYXQiOm51bGwsImV4cCI6bnVsbCwiYXVkIjoiIiwic3ViIjoiIn0.IlffGJz3IyFX1ADQ6-jOTQ_0D-K0kuKq5SpB_oirCrk"}),
				),
				signingMethod: jwt.SigningMethodHS256,
				err:           nil,
			},
			{
				ctx:           context.Background(),
				signingMethod: jwt.SigningMethodHS256,
				err:           errors.New("Bearer realm=Restricted"),
			},
			{
				ctx: metadata.NewIncomingContext(
					context.Background(),
					metadata.New(map[string]string{}),
				),
				signingMethod: jwt.SigningMethodHS256,
				err:           errors.New("Bearer realm=Restricted"),
			},
			{
				ctx: metadata.NewIncomingContext(
					context.Background(),
					metadata.New(map[string]string{authHeader: "Beare eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiIiLCJpYXQiOm51bGwsImV4cCI6bnVsbCwiYXVkIjoiIiwic3ViIjoiIn0.IlffGJz3IyFX1ADQ6-jOTQ_0D-K0kuKq5SpB_oirCrk"}),
				),
				signingMethod: jwt.SigningMethodHS256,
				err:           errors.New("Authorization header format must be Bearer {token}"),
			},
		}
	)
	for _, test := range tests {
		token, err := JWTFromContext(test.ctx, jwtKeyFunc, test.signingMethod)
		if e := testJWTFrom(token, err, test.err); e != nil {
			t.Error(e)
		}
	}
}

func testJWTFrom(token *jwt.Token, err, testErr error) error {
	if err != nil && err.Error() != testErr.Error() {
		return err
	}
	if testErr == nil && token == nil {
		return errors.New("Token not valid")
	}
	return nil
}
