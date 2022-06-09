package bwhttp_test

import (
	"fmt"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	bwhttp "github.com/kertox662/go-jwt-grpc/middleware/net_http"
)

func ExampleJWTHandler() {
	var (
		handler = func(w http.ResponseWriter, req *http.Request) {
			token, ok := bwhttp.JWTContext.ReadJWT(req)
			if !ok {
				http.Error(
					w,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError,
				)
			}
			fmt.Fprintf(w, "Token signed using %s", token.Method)
		}
		jwtKeyFunc = func(token *jwt.Token) (interface{}, error) {
			return []byte("MySecret"), nil
		}

		mux = http.NewServeMux()
	)
	mux.HandleFunc(
		"/",
		bwhttp.JWTHandler(handler, jwtKeyFunc, jwt.SigningMethodHS256),
	)

	if err := http.ListenAndServe("localhost:8080", mux); err != nil {
		log.Print(err)
	}
}
