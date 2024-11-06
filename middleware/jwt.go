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
		// Retrieve the JWT token from the cookie
		cookie, err := r.Cookie("token") // Use the cookie's name, such as "auth_token"
		if err != nil {
			http.Error(w, "Authorization cookie missing", http.StatusUnauthorized)
			return
		}

		// Parse and validate the JWT token
		tokenStr := cookie.Value
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// Ensure the token is signed with the correct algorithm
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrNotSupported
			}
			return jwtKey, nil // Replace with your actual secret key
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Extract the claims from the token
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

			// Convert the user ID to uint
			userID, err := strconv.Atoi(userIDStr)
			if err != nil {
				http.Error(w, "Invalid user ID", http.StatusUnauthorized)
				return
			}

			// Store the user ID in the request context
			ctx := context.WithValue(r.Context(), ContextUserIDKey, uint(userID))
			r = r.WithContext(ctx)
		} else {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
