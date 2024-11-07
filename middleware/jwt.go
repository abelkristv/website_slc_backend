package middleware

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("hehe")

type contextKey string

const ContextUserIDKey contextKey = "userID"

func TokenValid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		log.Print(cookie)
		if err != nil {
			http.Error(w, "Authorization cookie missing", http.StatusUnauthorized)
			return
		}

		tokenStr := cookie.Value
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
			log.Print("JWT")
			log.Print(claims["id"])

			var userIDStr string
			switch v := claims["id"].(type) {
			case string:
				userIDStr = v
				log.Print("hehe")
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
