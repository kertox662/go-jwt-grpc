package bearerware

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/golang-jwt/jwt/v4"
)

func TestJWTContext(t *testing.T) {
	var (
		ctx       = NewJWTContext()
		testReq   = &http.Request{}
		testToken = &jwt.Token{}
	)
	ctx.WriteJWT(testReq, testToken)
	token, ok := ctx.ReadJWT(testReq)
	if !ok {
		t.Fatal("expected jwt but got none")
	}
	if !reflect.DeepEqual(token, testToken) {
		t.Fatal("wrong token")
	}
	ctx.DeleteJWT(testReq)
	_, ok = ctx.ReadJWT(testReq)
	if ok {
		t.Fatal("expected no jwt but got one")
	}
}
