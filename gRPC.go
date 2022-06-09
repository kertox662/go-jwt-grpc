package bearerware

import (
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/net/context"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type jwtAccess struct {
	jsonKey string
}

/*
NewJWTAccessFromJWT creates a JWT credentials.PerRPCCredentials for use in gRPC
requests.
*/
func NewJWTAccessFromJWT(
	jsonKey string,
) (credentials.PerRPCCredentials, error) {
	return jwtAccess{jsonKey}, nil
}

func (j jwtAccess) GetRequestMetadata(
	ctx context.Context,
	uri ...string,
) (map[string]string, error) {
	return map[string]string{
		authHeader: fmt.Sprintf("%s%s", strings.Title(bearer), j.jsonKey),
	}, nil
}

func (j jwtAccess) RequireTransportSecurity() bool {
	return true
}

/*
JWTFromContext **deprecated** use `JWTFromIncomingContext`
*/
func JWTFromContext(
	ctx context.Context,
	keyFunc jwt.Keyfunc,
	signingMethod jwt.SigningMethod,
) (*jwt.Token, error) {
	return JWTFromIncomingContext(ctx, keyFunc, signingMethod)
}

/*
JWTFromIncomingContext extracts a valid JWT from a context.Contexts or returns
and error
*/
func JWTFromIncomingContext(
	ctx context.Context,
	keyFunc jwt.Keyfunc,
	signingMethod jwt.SigningMethod,
) (*jwt.Token, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errRestricted
	}
	var tokenStrings []string
	for k := range md {
		if authHeader == k {
			tokenStrings = md[k]
			break
		}
	}
	if len(tokenStrings) == 0 {
		return nil, errRestricted
	}

	tokenString, ok := tokenFromBearer(tokenStrings[0])
	if !ok {
		return nil, errBearerFormat
	}

	return validJWTFromString(tokenString, keyFunc, signingMethod)
}
