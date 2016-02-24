package bearerware

import (
	"net/http"
	"sync"

	"github.com/dgrijalva/jwt-go"
)

/*
JWTContexter provides and interface for safe access to a shared map to get a jwt
for the current request scope when using net/http.
*/
type JWTContexter interface {
	WriteJWT(*http.Request, *jwt.Token)
	ReadJWT(*http.Request) (*jwt.Token, bool)
	DeleteJWT(*http.Request)
}

type jwtContext struct {
	sync.RWMutex
	tokens map[*http.Request]*jwt.Token
}

//NewJWTContext creates a new JWTContexter
func NewJWTContext() JWTContexter {
	return &jwtContext{tokens: make(map[*http.Request]*jwt.Token)}
}

func (j *jwtContext) WriteJWT(req *http.Request, token *jwt.Token) {
	j.Lock()
	defer j.Unlock()
	j.tokens[req] = token
}

func (j *jwtContext) ReadJWT(req *http.Request) (*jwt.Token, bool) {
	j.RLock()
	defer j.RUnlock()
	token, ok := j.tokens[req]
	return token, ok
}

func (j *jwtContext) DeleteJWT(req *http.Request) {
	j.Lock()
	defer j.Unlock()
	delete(j.tokens, req)
}
