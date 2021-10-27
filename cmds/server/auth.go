package server

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

func generateToken(c loginDto) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": c.Name,
		"exp":  time.Now().Add(time.Minute * 15).Unix(),
	})
	var tokenString string
	secret := os.Getenv("SECRET")
	tokenString, err := t.SignedString([]byte(secret))
	if err != nil {
		return tokenString, &customErr{Err: errors.New("Error generating signing token"), Code: 500}
	}
	return tokenString, nil
}

func authMiddleware() mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			authorizationHeader := req.Header.Get("authorization")
			secret := os.Getenv("SECRET")
			if authorizationHeader != "" {
				bearerToken := strings.Split(authorizationHeader, " ")
				if len(bearerToken) == 2 {
					token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
						if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
							return nil, errors.New("error signing token method")
						}
						return []byte(secret), nil
					})
					if err != nil {
						http.Error(w, err.Error(), http.StatusUnauthorized)
						return
					}
					if !token.Valid {
						http.Error(w, "invalid authorization token", http.StatusUnauthorized)
						return
					}
				}
			} else {
				http.Error(w, "authorization header is required", http.StatusUnauthorized)
				return
			}
			h.ServeHTTP(w, req)
		})
	}
}
