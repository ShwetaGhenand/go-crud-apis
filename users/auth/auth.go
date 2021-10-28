package auth

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

type JWTToken struct {
	token *jwt.Token
}

type JWTSign struct {
	Key       interface{}
	Algorithm jwa.SignatureAlgorithm
}

func NewJWTToken(name, secret string) (string, error) {
	var tokenStr string
	js := &JWTSign{Key: []byte(secret), Algorithm: jwa.HS256}
	token := jwt.New()
	if err := token.Set(`name`, name); err != nil {
		return tokenStr, err
	}
	payload, err := jwt.Sign(token, js.Algorithm, js.Key)
	if err != nil {
		return tokenStr, err
	}
	tokenStr = string(payload)
	return tokenStr, nil
}

func ParseJWTToken(tokenStr, secret string) (*JWTToken, error) {
	token, err := jwt.Parse(
		[]byte(tokenStr),
		jwt.WithValidate(true),
		jwt.WithVerify(jwa.HS256, []byte(secret)))
	if err != nil {
		return nil, err
	}
	return &JWTToken{token: &token}, nil
}

func JWTMiddleware(secret string) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			authorizationHeader := req.Header.Get("authorization")
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) != 2 {
				log.Println("missing authorization token")
				http.Error(w, "access denied", http.StatusForbidden)
				return
			}
			if _, err := ParseJWTToken(bearerToken[1], secret); err != nil {
				log.Printf("error parsing authorization token %v", err)
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			h.ServeHTTP(w, req)
		})
	}
}
