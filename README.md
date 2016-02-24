# bearerware

[![Build Status](http://img.shields.io/travis/ckaznocha/go-JWTBearerware.svg?style=flat)](https://travis-ci.org/ckaznocha/go-JWTBearerware)
[![License](http://img.shields.io/:license-mit-blue.svg)](http://ckaznocha.mit-license.org)
[![GoDoc](https://godoc.org/github.com/ckaznocha/go-JWTBearerware?status.svg)](https://godoc.org/github.com/ckaznocha/go-JWTBearerware)
[![Go Report Card](https://goreportcard.com/badge/ckaznocha/go-JWTBearerware)](https://goreportcard.com/report/ckaznocha/go-JWTBearerware)

A Go library to make using [JSON Web Tokens](https://jwt.io/) in gRPC and HTTP requests more convenient. Middleware functions and examples for some popular routers are in the `midleware` directory.

This project was inspire by [auth0/go-jwt-middleware](https://github.com/auth0/go-jwt-middleware).

For more info see the example files and the [GoDoc](https://godoc.org/github.com/ckaznocha/go-JWTBearerware) page.

--
    import "github.com/ckaznocha/go-JWTBearerware"


## Usage

#### func  JWTFromContext

```go
func JWTFromContext(
	ctx context.Context,
	keyFunc jwt.Keyfunc,
	signingMethod jwt.SigningMethod,
) (*jwt.Token, error)
```
JWTFromContext extracts a valid JWT from a context.Contexts or returns and error

#### func  JWTFromHeader

```go
func JWTFromHeader(
	r *http.Request,
	keyFunc jwt.Keyfunc,
	signingMethod jwt.SigningMethod,
) (*jwt.Token, error)
```
JWTFromHeader extracts a valid JWT from an http.Request or returns and error

#### func  NewJWTAccessFromJWT

```go
func NewJWTAccessFromJWT(jsonKey string) (credentials.Credentials, error)
```
NewJWTAccessFromJWT creates a JWT credentials.Credentials which can be used in
gRPC requests.

#### func  WriteAuthError

```go
func WriteAuthError(w http.ResponseWriter, err error)
```
WriteAuthError is a convienence functon for setting the WWW-Authenticate header
and sending an http.Error()

#### type JWTContexter

```go
type JWTContexter interface {
	WriteJWT(*http.Request, *jwt.Token)
	ReadJWT(*http.Request) (*jwt.Token, bool)
	DeleteJWT(*http.Request)
}
```

JWTContexter provides and interface for safe access to a shared map to get a jwt
for the current request scope when using net/http.

#### func  NewJWTContext

```go
func NewJWTContext() JWTContexter
```
NewJWTContext creates a new JWTContexter
