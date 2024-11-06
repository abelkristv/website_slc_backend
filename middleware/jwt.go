package middleware

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("hehe")

const ContextUserIDKey = "id"

func TokenValid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrNotSupported
			}
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			log.Print(claims["id"])

			var userIDStr string
			switch v := claims["id"].(type) {
			case string:
				userIDStr = v
			case float64:
				userIDStr = strconv.FormatFloat(v, 'f', -1, 64)
			default:
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			userID, err := strconv.Atoi(userIDStr)
			if err != nil {
				http.Error(w, "Invalid user ID", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ContextUserIDKey, uint(userID))
			r = r.WithContext(ctx)
		} else {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
