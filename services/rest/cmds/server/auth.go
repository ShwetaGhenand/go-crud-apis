package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

// JWTSign holds the key and algorithm information for signing JWT tokens.
type JWTSign struct {
	key       []byte
	algorithm jwt.SigningMethod
}

// NewJWTToken generates a new JWT token with the given name and secret.
func NewJWTToken(name, secret string) (string, error) {
	// Create JWT signer
	signer := &JWTSign{key: []byte(secret), algorithm: jwt.SigningMethodHS256}

	// Create claims for JWT
	claims := jwt.MapClaims{
		"name": name,
	}

	// Create token
	token := jwt.NewWithClaims(signer.algorithm, claims)

	// Sign the token
	tokenStr, err := token.SignedString(signer.key)
	if err != nil {
		return "", err
	}

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
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorizationHeader := r.Header.Get("Authorization")
			if authorizationHeader == "" {
				log.Println("Missing Authorization token")
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			bearerToken := strings.Split(authorizationHeader, "Bearer ")
			if len(bearerToken) != 2 {
				log.Println("Invalid Authorization token format")
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			token := strings.TrimSpace(bearerToken[1])
			if err := ParseJWTToken(token, secret); err != nil {
				log.Printf("Error parsing Authorization token: %v", err)
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
