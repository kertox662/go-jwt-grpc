# bwhttp
--
    import "github.com/kertox662/go-jwt-grpc/middleware/net_http"


## Usage

```go
var JWTContext bearerware.JWTContexter
```
JWTContext stores the request scoped tokens

#### func  JWTHandler

```go
func JWTHandler(
	h http.HandlerFunc,
	keyFunc jwt.Keyfunc,
	signingMethod jwt.SigningMethod,
) http.HandlerFunc
```
JWTHandler is JWT middleware for net/http
