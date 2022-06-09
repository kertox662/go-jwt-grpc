package bwhttprouter_test

import (
	"fmt"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/julienschmidt/httprouter"
	bwhttprouter "github.com/kertox662/go-jwt-grpc/middleware/httprouter"
)

func ExampleJWTHandler() {
	var (
		handler = func(
			w http.ResponseWriter,
			req *http.Request,
			_ httprouter.Params,
		) {
			token, ok := bwhttprouter.JWTContext.ReadJWT(req)
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

		router = httprouter.New()
	)
	router.GET(
		"/",
		bwhttprouter.JWTHandler(handler, jwtKeyFunc, jwt.SigningMethodHS256),
	)

	if err := http.ListenAndServe("localhost:8080", router); err != nil {
		log.Print(err)
	}
}
