package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

type JWTSign struct {
	key       interface{}
	algorithm jwa.SignatureAlgorithm
}

func NewJWTToken(name, secret string) (string, error) {
	var tokenStr string
	js := &JWTSign{key: []byte(secret), algorithm: jwa.HS256}
	token := jwt.New()
	if err := token.Set(`name`, name); err != nil {
		return tokenStr, err
	}
	payload, err := jwt.Sign(token, js.algorithm, js.key)
	if err != nil {
		return tokenStr, err
	}
	tokenStr = string(payload)
	return tokenStr, nil
}

func ParseJWTToken(tokenStr, secret string) error {
	_, err := jwt.Parse(
		[]byte(tokenStr),
		jwt.WithValidate(true),
		jwt.WithVerify(jwa.HS256, []byte(secret)))
	if err != nil {
		return err
	}
	return nil
}

func authMiddleware(secret string) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			authorizationHeader := req.Header.Get("authorization")
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) != 2 {
				log.Println("missing authorization token")
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			if err := ParseJWTToken(bearerToken[1], secret); err != nil {
				log.Printf("error parsing authorization token %v", err)
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			h.ServeHTTP(w, req)
		})
	}
}
