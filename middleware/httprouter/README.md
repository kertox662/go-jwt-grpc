# bwhttprouter
--
    import "github.com/kertox662/go-jwt-grpc/middleware/httprouter"


## Usage

```go
var JWTContext bearerware.JWTContexter
```
JWTContext stores the request scoped tokens

#### func  JWTHandler

```go
func JWTHandler(
	h httprouter.Handle,
	keyFunc jwt.Keyfunc,
	signingMethod jwt.SigningMethod,
) httprouter.Handle
```
