package auth

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret = []byte("super-secret-key") // вынеси в env при проде

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPassword := os.Getenv("TODO_PASSWORD")
		if expectedPassword == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		tokenStr := cookie.Value

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return secret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		hashFromToken := claims["hash"].(string)
		currentHash := fmt.Sprintf("%x", sha256.Sum256([]byte(expectedPassword)))

		if hashFromToken != currentHash {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	})
}

// Генерация токена (используется в signin)
func GenerateToken(password string) (string, error) {
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": hash,
		"exp":  jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
	})

	return token.SignedString(secret)
}
